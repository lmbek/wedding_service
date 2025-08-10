package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"wedding_service/webserver/database"
)

// RSVP is an interface-first service that manages per-invitation accept/decline state.
// We never expose the concrete struct outside the package.
type RSVP interface {
	List(code string) AcceptedByCode
	Accept(code, name string) AcceptedByCode
	Decline(code, name string) AcceptedByCode
}

type acceptedStore struct {
	mu      sync.Mutex
	invites database.Invites
	// per code we keep a set of accepted names
	acc map[string]map[string]struct{}
}

type AcceptedByCode struct {
	Code     string   `json:"code"`
	Members  []string `json:"members"`
	Accepted []string `json:"accepted"`
	Count    int      `json:"count"`
	Capacity int      `json:"capacity"`
}

func NewRSVP(invites database.Invites) RSVP {
	return &acceptedStore{invites: invites, acc: make(map[string]map[string]struct{})}
}

func (s *acceptedStore) snapshot(code string) AcceptedByCode {
	inv, ok := s.invites.FindByCode(code)
	if !ok {
		return AcceptedByCode{Code: code}
	}
	// ensure order stable: iterate members and include only accepted subset
	st := s.acc[code]
	accepted := make([]string, 0, len(inv.Members))
	for _, m := range inv.Members {
		if st != nil {
			if _, ok2 := st[m]; ok2 {
				accepted = append(accepted, m)
			}
		}
	}
	return AcceptedByCode{
		Code:     code,
		Members:  append([]string(nil), inv.Members...),
		Accepted: accepted,
		Count:    len(accepted),
		Capacity: len(inv.Members),
	}
}

func (s *acceptedStore) List(code string) AcceptedByCode {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshot(code)
}

func (s *acceptedStore) Accept(code, name string) AcceptedByCode {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.acc[code]; !ok {
		s.acc[code] = make(map[string]struct{})
	}
	s.acc[code][name] = struct{}{}
	return s.snapshot(code)
}

func (s *acceptedStore) Decline(code, name string) AcceptedByCode {
	s.mu.Lock()
	defer s.mu.Unlock()
	if st, ok := s.acc[code]; ok {
		delete(st, name)
	}
	return s.snapshot(code)
}

// Package-level singleton wired in mux (via SetRSVP) to respect DI style.
var rsvp RSVP

func SetRSVP(s RSVP) { rsvp = s }

// HTTP handlers for per-code operations.

type memberPayload struct {
	Name string `json:"name"`
}

func GetAcceptedByCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	resp := rsvp.List(code)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func PostAcceptByCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	var p memberPayload
	_ = json.NewDecoder(r.Body).Decode(&p)
	resp := rsvp.Accept(code, p.Name)
	// also update global accepted list for the overview table
	acceptor.Accept(p.Name)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func PostDeclineByCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	var p memberPayload
	_ = json.NewDecoder(r.Body).Decode(&p)
	resp := rsvp.Decline(code, p.Name)
	// update global list accordingly
	acceptor.Decline(p.Name)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
