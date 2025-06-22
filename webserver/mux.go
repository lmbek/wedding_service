package webserver

import (
	"net/http"
	"wedding_service/webserver/website"
)

func useWebsite(m *http.ServeMux) {
	m.HandleFunc("GET /{$}", website.FrontPageHandler)
	m.HandleFunc("GET /invitation/{$}", website.InvitationPageHandler)
}

func useApi(m *http.ServeMux) {
	m.HandleFunc("GET /api/{$}", null)
}

func null(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
