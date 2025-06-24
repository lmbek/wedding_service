package api

import (
	"encoding/json"
	"net/http"
)

// Person represents a guest/person in the wedding service.
type Person struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdatePerson allows partial updates.
type UpdatePerson struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

var persons = []Person{
	{ID: 1, Name: "Alice", Email: "alice@example.com"},
	{ID: 2, Name: "Bob", Email: "bob@example.com"},
}

// ListPersonsHandler godoc
// @Summary      List all persons
// @Description  Returns a list of all persons in the system.
// @Tags         persons
// @Produce      json
// @Success      200  {array}  Person
// @Failure      500  {object}  map[string]string
// @Router       /api/persons/ [get]
func ListPersonsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(persons)
}

// GetPersonHandler godoc
// @Summary      Get a person by ID
// @Description  Returns details of a single person by their ID.
// @Tags         persons
// @Produce      json
// @Param        id   path      int  true  "Person ID"
// @Success      200  {object}  Person
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/persons/{id}/ [get]
func GetPersonHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}
	for _, p := range persons {
		if p.ID == id {
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, `{"error":"Person not found"}`, http.StatusNotFound)
}

// PostPersonHandler godoc
// @Summary      Create a new person
// @Description  Creates and returns a new person entry.
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        person  body      Person  true  "New person data"
// @Success      201     {object}  Person
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api/persons/ [post]
func PostPersonHandler(w http.ResponseWriter, r *http.Request) {
	var p Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	p.ID = getNextID()
	persons = append(persons, p)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// PutPersonHandler godoc
// @Summary      Update a person partially
// @Description  Partially updates a person by ID. Fields left out will remain unchanged.
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        id      path      int           true  "Person ID"
// @Param        person  body      UpdatePerson  true  "Updated person data (partial)"
// @Success      200     {object}  Person
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api/persons/{id}/ [put]
func PutPersonHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}
	var patch UpdatePerson
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	for i, p := range persons {
		if p.ID == id {
			if patch.Name != nil {
				p.Name = *patch.Name
			}
			if patch.Email != nil {
				p.Email = *patch.Email
			}
			persons[i] = p
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, `{"error":"Person not found"}`, http.StatusNotFound)
}

// DeletePersonHandler godoc
// @Summary      Delete a person
// @Description  Deletes a person from the system by ID.
// @Tags         persons
// @Produce      json
// @Param        id   path      int  true  "Person ID"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/persons/{id}/ [delete]
func DeletePersonHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}
	for i, p := range persons {
		if p.ID == id {
			persons = append(persons[:i], persons[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, `{"error":"Person not found"}`, http.StatusNotFound)
}
