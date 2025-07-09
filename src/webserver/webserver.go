package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"wedding_service/buildtag"
	"wedding_service/certificate"
	"wedding_service/env"
	"wedding_service/webserver/website/frontend"
)

type Webserver interface {
	ListenAndServe() error
	Close() error
}

type webserver struct {
	httpServer  *http.Server
	httpsServer *http.Server
	certPath    string
	keyPath     string
}

func NewWebserver(newFrontend frontend.Frontend) (w Webserver, err error) {
	// TODO: a suggestion could be to use: useHostProtection(mux) and then wrap the mux (maybe)

	initAllowedHosts()

	mux := http.NewServeMux()
	useWebsite(mux, newFrontend)
	useApi(mux)

	acmeManager, err := certificate.InitAcme()
	if err != nil {
		return nil, fmt.Errorf("could not use acme manager: %w", err)
	}

	httpServer := newHttpServer(env.Env.HttpPort)
	httpsServer, err := newHttpsServer(env.Env.HttpsPort, acmeManager)
	if err != nil {
		return nil, err
	}

	// TODO: middleware can be collected in one func
	httpServer.Handler = LoggingAndMetricsMiddleware(acmeManager.HTTPHandler(redirectToHTTPS(env.Env.HttpsPort))) //ProtectHostsMiddleware(LoggingAndMetricsMiddleware(acmeManager.HTTPHandler(redirectToHTTPS(env.Env.HttpsPort))))
	httpsServer.Handler = LoggingAndMetricsMiddleware(mux)                                                        //ProtectHostsMiddleware(LoggingAndMetricsMiddleware(mux))
	//httpServer.Handler = ProtectHostsMiddleware(LoggingAndMetricsMiddleware(acmeManager.HTTPHandler(redirectToHTTPS(env.Env.HttpsPort))))
	//httpsServer.Handler = ProtectHostsMiddleware(LoggingAndMetricsMiddleware(mux))

	return &webserver{
		httpServer:  httpServer,
		httpsServer: httpsServer,
	}, nil
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
	if env.IsDebugInfoEnabled() && buildtag.IsDevelopment() {
		fmt.Println("Listening on", fmt.Sprintf("https://localhost:%s", env.Env.HttpsPort))
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
	if env.IsDebugInfoEnabled() && buildtag.IsDevelopment() {
		fmt.Println("Listening on", fmt.Sprintf("http://localhost:%s", env.Env.HttpPort))
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
