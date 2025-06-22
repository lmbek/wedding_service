package webserver

import (
	"fmt"
	"net/http"
	"wedding_service/webserver/website"
)

func useWebsite(m *http.ServeMux) {

	fs := NewFileServer()

	m.HandleFunc("GET /", fs.Serve)

	m.HandleFunc("GET /{$}", website.FrontPageHandler)
	m.HandleFunc("GET /invitation/{$}", website.InvitationPageHandler)
}

func useApi(m *http.ServeMux) {
	fmt.Println("using api")
	//m.HandleFunc("GET /api/{$}", null)
}
