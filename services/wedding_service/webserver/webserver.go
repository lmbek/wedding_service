package webserver

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"wedding_service/config"
	"wedding_service/webserver/website/frontend"
)

type Webserver interface {
	init() (Webserver, error)
	ListenAndServe() error
	Close() error
}

type webserver struct {
	config     config.Config
	frontend   frontend.Frontend
	httpServer *http.Server
}

func NewWebserver(config config.Config, frontend frontend.Frontend) (Webserver, error) {
	w := &webserver{
		config:   config,
		frontend: frontend,
	}
	return w.init()
}

func (w *webserver) init() (Webserver, error) {
	mux := http.NewServeMux()
	useWebsite(w.config, mux, w.frontend)
	useApi(w.config, mux)

	httpServer := newHttpServer(w.config.HttpPort())
	httpServer.Handler = applyMiddleware(mux)

	w.httpServer = httpServer

	return w, nil
}

func (w *webserver) ListenAndServe() error {
	slog.Info("Listening on HTTP", slog.String("url", fmt.Sprintf("http://localhost:%s", w.config.HttpPort())))
	err := w.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server error: %w", err)
	}
	return nil
}

func (w *webserver) Close() error {
	return w.httpServer.Close()
}

// applyMiddleware applies a chain of middleware to a handler
// The middleware is applied in reverse order, so the first middleware in the list
// will be the outermost middleware in the chain
func applyMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
