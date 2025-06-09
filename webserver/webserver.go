package webserver

import (
	"errors"
	"fmt"
	"net/http"
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
	var httpServer *http.Server = newHttpServer(env.Env.HttpPort)
	var httpsServer *http.Server = newHttpsServer(env.Env.HttpsPort)

	// use certificate for https/tls
	err = useCertificate(httpsServer)
	if err != nil {
		return nil, err
	}

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
	err := w.httpsServer.ListenAndServeTLS("", "")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTPS server error: %w", err)
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
