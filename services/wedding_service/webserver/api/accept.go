package api

import (
	"encoding/json"
	"net/http"
	"wedding_service/webserver/database"
)

// AcceptApp is an interface-first in-memory RSVP accept service.
// We do not expose concrete structs outside this package.
type AcceptApp interface {
	Accept(name string) (AcceptedResponse, error)
	Decline(name string) (AcceptedResponse, error)
	List() AcceptedResponse
}

type acceptApp interface {
	Accept(name string) (AcceptedResponse, error)
	Decline(name string) (AcceptedResponse, error)
	List() AcceptedResponse
}

type dbAcceptApp struct {
	invites  database.Invites
	capacity int
}

func NewAcceptApp(invites database.Invites) AcceptApp {
	return &dbAcceptApp{invites: invites, capacity: 45}
}

type AcceptedResponse struct {
	Accepted []string `json:"accepted"`
	Count    int      `json:"count"`
	Capacity int      `json:"capacity"`
}

func (a *dbAcceptApp) snapshot() AcceptedResponse {
	list, _ := a.invites.ListAllAccepted()
	out := make([]string, len(list))
	copy(out, list)
	return AcceptedResponse{Accepted: out, Count: len(out), Capacity: a.capacity}
}

func (a *dbAcceptApp) Accept(name string) (AcceptedResponse, error) {
	// Global accept endpoint is read-only; per-invite accept persists via DB in RSVP handlers.
	return a.snapshot(), nil
}

func (a *dbAcceptApp) Decline(name string) (AcceptedResponse, error) {
	// Global decline endpoint is read-only; per-invite decline persists via DB in RSVP handlers.
	return a.snapshot(), nil
}

func (a *dbAcceptApp) List() AcceptedResponse {
	return a.snapshot()
}

// Singleton instance (interface type) kept inside package scope.
var acceptor AcceptApp

// SetAccept wires the global AcceptApp implementation (DI from mux).
func SetAccept(a AcceptApp) { acceptor = a }

// HTTP handlers (thin layer) using the interface above.

type acceptPayload struct {
	Name string `json:"name"`
}

func GetAcceptedHandler(w http.ResponseWriter, r *http.Request) {
	var resp AcceptedResponse
	if acceptor != nil {
		resp = acceptor.List()
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func PostAcceptHandler(w http.ResponseWriter, r *http.Request) {
	var p acceptPayload
	_ = json.NewDecoder(r.Body).Decode(&p)
	var resp AcceptedResponse
	if acceptor != nil {
		resp, _ = acceptor.Accept(p.Name)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func PostDeclineHandler(w http.ResponseWriter, r *http.Request) {
	var p acceptPayload
	_ = json.NewDecoder(r.Body).Decode(&p)
	var resp AcceptedResponse
	if acceptor != nil {
		resp, _ = acceptor.Decline(p.Name)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
