package webserver

import (
	"github.com/swaggo/http-swagger"
	"golang.org/x/net/websocket"
	"net/http"
	"wedding_service/config"
	"wedding_service/webserver/api"
	"wedding_service/webserver/website"
	"wedding_service/webserver/website/frontend"
)

// automated swagger generate on general generate
//go:generate swag init --dir .. --output ../webserver/website/frontend/out/public/api/swagger --parseDependency

// manual swagger generate:
// swag init --output webserver/website/frontend/out/public/api/swagger --parseDependency

func useWebsite(config config.Config, m *http.ServeMux, newFrontend frontend.Frontend) {
	render := website.NewRender(config, newFrontend)
	// on the files on the frontend is not getting renewed
	m.HandleFunc("GET /", newFrontend.Serve)
	m.HandleFunc("GET /{$}", render.FrontPageHandler)
	m.HandleFunc("GET /invitation/{$}", render.InvitationPageHandler)
	m.Handle("GET /websocket/hotreload/{$}", websocket.Handler(func(ws *websocket.Conn) {
		frontend.HandleRegisterClient(ws, config.FrontendPath(), config.HotReloadEnabled())
	}))
}

func useApi(config config.Config, m *http.ServeMux) {
	// Configure Swagger UI to use doc.json
	m.HandleFunc("GET /api/swagger/", httpSwagger.WrapHandler)

	// Serve swagger.json as doc.json for Swagger UI
	m.HandleFunc("GET /api/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "webserver/website/frontend/out/public/api/swagger/swagger.json")
	})

	// OpenTelemetry metrics are collected via logs, no endpoint needed

	m.HandleFunc("GET /api/persons/{$}", api.ListPersonsHandler)
	m.HandleFunc("GET /api/persons/{id}/{$}", api.GetPersonHandler)
	m.HandleFunc("POST /api/persons/{$}", api.PostPersonHandler)
	m.HandleFunc("PUT /api/persons/{id}/{$}", api.PutPersonHandler)
	m.HandleFunc("DELETE /api/persons/{id}/{$}", api.DeletePersonHandler)
}
