package reverseproxy

import (
	"errors"
	"fmt"
	"gateway_service/config"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gateway_service/middleware"
	"gateway_service/tracer"
)

type ReverseProxy interface {
	AddHost(host string, target *url.URL) error
	Handler() http.Handler
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HTTPServer() *http.Server
	HTTPSServer() *http.Server
	ListenAndServe() error
	Close() error
}

type reverseProxy struct {
	hosts       map[string]*httputil.ReverseProxy
	middleware  middleware.Middleware
	transport   tracer.Transport
	httpServer  *http.Server
	httpsServer *http.Server
	config      config.Config
}

func (rp *reverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rp.Handler().ServeHTTP(w, r)
}

func NewReverseProxy(config config.Config, httpServer *http.Server, httpsServer *http.Server, middleware middleware.Middleware, transportTracer tracer.Transport) ReverseProxy {
	return &reverseProxy{
		config:      config,
		httpServer:  httpServer,
		httpsServer: httpsServer,
		hosts:       make(map[string]*httputil.ReverseProxy),
		middleware:  middleware,
		transport:   transportTracer,
	}
}

func (rp *reverseProxy) AddHost(host string, url *url.URL) error {
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" || url == nil {
		return nil
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = rp.transport
	rp.hosts[host] = proxy
	return nil
}

func (rp *reverseProxy) Handler() http.Handler {
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h := strings.ToLower(req.Host)
		if i := strings.Index(h, ":"); i > -1 {
			h = h[:i]
		}
		if proxy, ok := rp.hosts[h]; ok {
			proxy.ServeHTTP(w, req)
			return
		}
		slog.Warn("no host matched, gave forbidden", slog.String("host", req.Host))
		w.WriteHeader(http.StatusForbidden)
	})
	if rp.middleware != nil {
		return rp.middleware.Handler(proxyHandler)
	}
	return proxyHandler
}

func (rp *reverseProxy) ListenAndServe() error {
	errChan := make(chan error, 2)

	go func() {
		errChan <- rp.listenHTTPS()
	}()

	go func() {
		errChan <- rp.listenHTTP()
	}()

	return <-errChan
}

func (rp *reverseProxy) listenHTTPS() error {
	if rp.httpsServer == nil {
		return fmt.Errorf("HTTPS server not initialized")
	}

	// Listen and serve HTTPS / TLS
	// if there is a problem with these, the problem could be that the cert and key is not set before this (for example, missing docker-dev mode)
	err := rp.httpsServer.ListenAndServeTLS("", "")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTPS server error (cert not loaded?): %w", err)
	}

	return nil
}

func (rp *reverseProxy) listenHTTP() error {
	if rp.httpServer == nil {
		return fmt.Errorf("HTTP server not initialized")
	}

	// Listen and serve HTTP
	err := rp.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}

func (rp *reverseProxy) HTTPServer() *http.Server {
	return rp.httpServer
}

func (rp *reverseProxy) HTTPSServer() *http.Server {
	return rp.httpsServer
}

func (rp *reverseProxy) Close() error {
	var errClosingHttps, errClosingHttp error
	if rp.httpsServer != nil {
		errClosingHttps = rp.httpsServer.Close()
	}
	if rp.httpServer != nil {
		errClosingHttp = rp.httpServer.Close()
	}
	return errors.Join(errClosingHttp, errClosingHttps)
}
