package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

// AcceptApp is an interface-first in-memory RSVP accept service.
// We do not expose concrete structs outside this package.
type AcceptApp interface {
	Accept(name string) (AcceptedResponse, error)
	Decline(name string) (AcceptedResponse, error)
	List() AcceptedResponse
}

type acceptApp struct {
	mu       sync.Mutex
	set      map[string]struct{}
	order    []string
	capacity int
}

func NewAcceptApp() AcceptApp {
	return &acceptApp{set: make(map[string]struct{}), capacity: 45}
}

type AcceptedResponse struct {
	Accepted []string `json:"accepted"`
	Count    int      `json:"count"`
	Capacity int      `json:"capacity"`
}

func (a *acceptApp) snapshot() AcceptedResponse {
	out := make([]string, len(a.order))
	copy(out, a.order)
	return AcceptedResponse{Accepted: out, Count: len(out), Capacity: a.capacity}
}

func (a *acceptApp) Accept(name string) (AcceptedResponse, error) {
	if name == "" {
		name = "Gæst"
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.set[name]; !ok {
		a.set[name] = struct{}{}
		a.order = append(a.order, name)
	}
	return a.snapshot(), nil
}

func (a *acceptApp) Decline(name string) (AcceptedResponse, error) {
	if name == "" {
		name = "Gæst"
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.set[name]; ok {
		delete(a.set, name)
		// remove from order slice
		for i, n := range a.order {
			if n == name {
				a.order = append(a.order[:i], a.order[i+1:]...)
				break
			}
		}
	}
	return a.snapshot(), nil
}

func (a *acceptApp) List() AcceptedResponse {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.snapshot()
}

// Singleton instance (interface type) kept inside package scope.
var acceptor AcceptApp = NewAcceptApp()

// HTTP handlers (thin layer) using the interface above.

type acceptPayload struct {
	Name string `json:"name"`
}

func GetAcceptedHandler(w http.ResponseWriter, r *http.Request) {
	resp := acceptor.List()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func PostAcceptHandler(w http.ResponseWriter, r *http.Request) {
	var p acceptPayload
	_ = json.NewDecoder(r.Body).Decode(&p)
	resp, _ := acceptor.Accept(p.Name)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func PostDeclineHandler(w http.ResponseWriter, r *http.Request) {
	var p acceptPayload
	_ = json.NewDecoder(r.Body).Decode(&p)
	resp, _ := acceptor.Decline(p.Name)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
