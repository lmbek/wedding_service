package webserver

import (
	"log"
	"net/http"
	"wedding_service/config"
	"wedding_service/webserver/api"
	"wedding_service/webserver/database"
	"wedding_service/webserver/website"
	"wedding_service/webserver/website/frontend"

	"github.com/swaggo/http-swagger"
	"golang.org/x/net/websocket"
)

// automated swagger generate on general generate
//go:generate swag init --dir .. --output ../webserver/website/frontend/out/public/api/swagger --parseDependency

// manual swagger generate:
// swag init --output webserver/website/frontend/out/public/api/swagger --parseDependency

func useWebsite(cfg config.Config, m *http.ServeMux, newFrontend frontend.Frontend) {
	// Try MySQL-backed invites first
	var invites database.Invites
	user := cfg.MySQLUser()
	host := cfg.MySQLHost()
	port := cfg.MySQLPort()
	pass := cfg.MySQLPassword()
	dbname := cfg.MySQLDatabase()
	dsn := ""
	if user != "" && host != "" && port != "" && dbname != "" {
		// Attempt to ensure the database exists using the application user only
		dbInit := database.NewDBInit()
		errInit := dbInit.EnsureDatabase(host, port, user, pass, dbname)
		if errInit != nil {
			log.Printf("[startup] database init attempt failed for %s@%s:%s/%s: %v", user, host, port, dbname, errInit)
		} else {
			log.Printf("[startup] database ensured (or already exists) for %s@%s:%s/%s", user, host, port, dbname)
		}
		dsn = user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true&charset=utf8mb4,utf8"
	} else {
		log.Printf("[startup] missing DB env; using read-only invites (no in-memory RSVP state) (user=%q host=%q port=%q db=%q)", user, host, port, dbname)
	}
	var err error
	if dsn != "" {
		invites, err = database.NewInvitesMySQL(dsn)
		if err != nil {
			log.Printf("[startup] MySQL invites backend failed; using read-only fallback: %v", err)
			invites = database.NewInvites()
		}
	} else {
		invites = database.NewInvites()
	}
	render := website.NewRender(cfg, newFrontend, invites)
	// Wire RSVP and Accept services with invites (DI)
	api.SetRSVP(api.NewRSVP(invites))
	api.SetAccept(api.NewAcceptApp(invites))
	// on the files on the frontend is not getting renewed
	// Root redirects to /bryllup/
	m.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/bryllup/", http.StatusMovedPermanently)
	})

	// Static assets
	m.HandleFunc("GET /bryllup/", newFrontend.Serve)
	// Serve embedded/public assets under /bryllup/ by stripping the /bryllup/ prefix

	// Serve static frontend assets under /bryllup/
	m.HandleFunc("GET /bryllup/{$}", render.FrontPageHandler)
	m.HandleFunc("GET /bryllup/invitation/{$}", render.InvitationPageHandler)
	m.HandleFunc("GET /bryllup/invitation/{code}/{$}", render.InvitationPageHandler)
	m.HandleFunc("GET /bryllup/menu/{$}", render.MenuPageHandler)
	m.HandleFunc("GET /bryllup/rsvp/{$}", render.RSVPPageHandler)
	m.HandleFunc("GET /bryllup/info/{$}", render.InfoPageHandler)
	m.HandleFunc("GET /bryllup/wishes/{$}", render.WishesPageHandler)
	// Hidden endpoint: not linked in navigation
	m.HandleFunc("GET /bryllup/reception/{$}", render.ReceptionPageHandler)

	//m.Handle("GET /bryllup/{rest...}", http.StripPrefix("/bryllup/", http.HandlerFunc(newFrontend.Serve)))
	m.Handle("GET /bryllup/websocket/hotreload/{$}", websocket.Handler(func(ws *websocket.Conn) {
		frontend.HandleRegisterClient(ws, cfg.FrontendPath(), cfg.HotReloadEnabled())
	}))
}

func useApi(config config.Config, m *http.ServeMux) {
	// Configure Swagger UI to use doc.json
	m.HandleFunc("GET /bryllup/api/swagger/", httpSwagger.WrapHandler)

	// Serve swagger.json as doc.json for Swagger UI
	m.HandleFunc("GET /bryllup/api/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "webserver/website/frontend/out/public/api/swagger/swagger.json")
	})

	// OpenTelemetry metrics are collected via logs, no endpoint needed

	m.HandleFunc("GET /bryllup/api/persons/{$}", api.ListPersonsHandler)
	m.HandleFunc("GET /bryllup/api/persons/{id}/{$}", api.GetPersonHandler)
	m.HandleFunc("POST /bryllup/api/persons/{$}", api.PostPersonHandler)
	m.HandleFunc("PUT /bryllup/api/persons/{id}/{$}", api.PutPersonHandler)
	m.HandleFunc("DELETE /bryllup/api/persons/{id}/{$}", api.DeletePersonHandler)

	// Invitations acceptance endpoints
	m.HandleFunc("GET /bryllup/api/invites/accepted/{$}", api.GetAcceptedHandler)
	m.HandleFunc("POST /bryllup/api/invites/accept/{$}", api.PostAcceptHandler)
	m.HandleFunc("POST /bryllup/api/invites/decline/{$}", api.PostDeclineHandler)
	// Per-invitation endpoints
	m.HandleFunc("GET /bryllup/api/invites/{code}/accepted/{$}", api.GetAcceptedByCodeHandler)
	m.HandleFunc("POST /bryllup/api/invites/{code}/accept/{$}", api.PostAcceptByCodeHandler)
	m.HandleFunc("POST /bryllup/api/invites/{code}/decline/{$}", api.PostDeclineByCodeHandler)
}
