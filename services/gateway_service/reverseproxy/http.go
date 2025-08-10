package reverseproxy

import (
	"net/http"
	"strings"
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

	// TODO: fjern det her og bare antag at acme altid virker
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

// NewHTTPRedirectOrProxyHandler returns a handler that proxies WebSocket upgrade
// requests to the provided proxy handler, while redirecting normal HTTP traffic
// to HTTPS. ACME challenges are still handled when an ACME manager is provided.
func NewHTTPRedirectOrProxyHandler(
	acme certificate.AutoCertManager,
	redirectFunc func(http.ResponseWriter, *http.Request),
	proxy http.Handler,
) http.Handler {
	redirect := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if redirectFunc != nil {
			redirectFunc(w, r)
			return
		}
		target := "https://" + r.Host + r.URL.RequestURI()
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connHdr := strings.ToLower(r.Header.Get("Connection"))
		upgrade := strings.ToLower(r.Header.Get("Upgrade"))
		if strings.Contains(connHdr, "upgrade") && upgrade == "websocket" {
			proxy.ServeHTTP(w, r)
			return
		}
		redirect.ServeHTTP(w, r)
	})

	if acme != nil {
		return acme.HTTPHandler(handler)
	}
	return handler
}
