package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"wedding_service/certificate"
	"wedding_service/env"
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

func NewWebserver() (w Webserver, err error) {
	if env.IsModeNotSet() {
		return nil, errors.New("no MODE set in .env")
	}

	mux := http.NewServeMux()
	useWebsite(mux)
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

	// TODO: add middleware with protect host names and logging and more to these.
	// Maybe useMiddleware(httpsServer) and then put handle there etc.
	httpServer.Handler = acmeManager.HTTPHandler(mux)
	httpsServer.Handler = mux

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
	if env.IsDebugInfoEnabled() && env.IsModeDevelopment() {
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
	if env.IsDebugInfoEnabled() && env.IsModeDevelopment() {
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
