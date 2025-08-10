package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"wedding_service/buildtag"
	"wedding_service/config"
	"wedding_service/webserver/certificate"
	"wedding_service/webserver/website/frontend"
)

type Webserver interface {
	init() (Webserver, error)
	ListenAndServe() error
	Close() error
}

type webserver struct {
	config      config.Config
	frontend    frontend.Frontend
	httpServer  *http.Server
	httpsServer *http.Server
}

func NewWebserver(config config.Config, frontend frontend.Frontend) (Webserver, error) {
	w := &webserver{
		config:   config,
		frontend: frontend,
	}
	return w.init()
}

func (w *webserver) init() (Webserver, error) {
	// TODO: a suggestion could be to use: useHostProtection(mux) and then wrap the mux (maybe)

	initAllowedHosts(w.config)

	mux := http.NewServeMux()
	useWebsite(w.config, mux, w.frontend)
	useApi(w.config, mux)

	acmeManager, err := certificate.NewAutoCertManager(w.config)
	if err != nil {
		return nil, fmt.Errorf("could not use acme manager: %w", err)
	}

	httpServer := newHttpServer(w.config.HttpPort())
	httpsServer, err := newHttpsServer(w.config.CertPath(), w.config.KeyPath(), w.config.HttpsPort(), acmeManager)
	if err != nil {
		return nil, err
	}

	// TODO: middleware can be collected in one func
	httpServer.Handler = LoggingAndMetricsMiddleware(acmeManager.HTTPHandler(redirectToHTTPS(w.config.HttpsPort()))) //ProtectHostsMiddleware(LoggingAndMetricsMiddleware(acmeManager.HTTPHandler(redirectToHTTPS(env.Env.HttpsPort))))
	httpsServer.Handler = LoggingAndMetricsMiddleware(mux)                                                           //ProtectHostsMiddleware(LoggingAndMetricsMiddleware(mux))
	//httpServer.Handler = ProtectHostsMiddleware(LoggingAndMetricsMiddleware(acmeManager.HTTPHandler(redirectToHTTPS(env.Env.HttpsPort))))
	//httpsServer.Handler = ProtectHostsMiddleware(LoggingAndMetricsMiddleware(mux))

	w.httpServer = httpServer
	w.httpsServer = httpsServer

	return w, nil
}

func (w *webserver) ListenAndServe() error {
	errChan := make(chan error, 2)

	go func() {
		errChan <- w.listenHTTPS()
	}()

	go func() {
		errChan <- w.listenHTTP()
	}()

	return <-errChan
}

func (w *webserver) listenHTTPS() error {
	if w.config.IsDebugInfoEnabled() && buildtag.IsDevelopment() {
		fmt.Println("Listening on", fmt.Sprintf("https://localhost:%s", w.config.HttpsPort()))
	}

	// Listen and serve HTTPS / TLS
	// if there is a problem with these, the problem could be that the cert and key is not set before this (for example, missing docker-dev mode)
	err := w.httpsServer.ListenAndServeTLS("", "")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTPS server error (cert not loaded?): %w", err)
	}

	return nil
}

func (w *webserver) listenHTTP() error {
	if w.config.IsDebugInfoEnabled() && buildtag.IsDevelopment() {
		fmt.Println("Listening on", fmt.Sprintf("http://localhost:%s", w.config.HttpPort()))
	}

	// Listen and serve HTTP
	err := w.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}

func (w *webserver) Close() error {
	errClosingHttps := w.httpsServer.Close()
	errClosingHttp := w.httpServer.Close()

	return errors.Join(errClosingHttp, errClosingHttps)
}
