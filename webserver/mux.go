package webserver

import "net/http"

func useWebsite(m *http.ServeMux) {
	m.HandleFunc("GET /{$}", null)
}

func useApi(m *http.ServeMux) {
	m.HandleFunc("GET /api/{$}", null)
}

func null(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
