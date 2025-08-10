package app

import (
	"context"
	"fmt"
	"net/http"
)

type App interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type app struct {
	httpServer *http.Server
}

func NewApp() App {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello batman"))
	})

	return &app{
		httpServer: &http.Server{
			Addr:    ":80",
			Handler: mux,
		},
	}
}

func (a *app) Run() error {
	fmt.Println("tinyserver listening on http://0.0.0.0:80")
	return a.httpServer.ListenAndServe()
}

func (a *app) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
