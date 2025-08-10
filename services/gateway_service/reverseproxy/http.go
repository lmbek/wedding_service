package reverseproxy

import (
	"net/http"
	"time"

	"gateway_service/certificate"
)

type HTTPServer struct {
	Port         string
	ACME         certificate.AutoCertManager
	RedirectFunc func(http.ResponseWriter, *http.Request) // fallback redirect
}

func NewHTTPServer(server HTTPServer) *http.Server {
	if server.Port == "" {
		server.Port = "80"
	}
	redirect := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if server.RedirectFunc != nil {
			server.RedirectFunc(w, r)
			return
		}
		target := "https://" + r.Host + r.URL.RequestURI()
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})

	var handler http.Handler = redirect
	if server.ACME != nil {
		handler = server.ACME.HTTPHandler(redirect) // håndterer ACME HTTP-01 challenges
	}

	return &http.Server{
		Addr:         ":" + server.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
