package webserver

import (
	"crypto/tls"
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
	httpPort := env.Env.HttpPort
	httpsPort := env.Env.HttpsPort

	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%s", httpPort),
	}

	httpsServer := &http.Server{
		Addr: fmt.Sprintf(":%s", httpsPort),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // HTTP/2 requires TLS 1.2 or higher
			NextProtos: []string{"h2"},   // Enforce HTTP/2
		},
	}

	if env.IsModeDevelopment() {
		cert, err := certificate.UseLOCALHOST()
		if err != nil {
			return nil, fmt.Errorf("could not use localhost certificate: %w", err)
		}
		httpsServer.TLSConfig.Certificates = []tls.Certificate{*cert}
	}

	if env.IsModeProduction() {
		acmeManager, err := certificate.UseACME()
		if err != nil {
			return nil, fmt.Errorf("could not use acme manager: %w", err)
		}
		httpsServer.TLSConfig.GetCertificate = acmeManager.GetCertificate
	}

	if env.IsModeNotSet() {
		return nil, errors.New("no MODE set in .env")
	}

	return &webserver{
		httpServer:  httpServer,
		httpsServer: httpsServer,
	}, nil
}

func (w *webserver) ListenAndServe() error {
	errChan := make(chan error, 2)

	go func() {
		if env.IsDebugInfoEnabled() && env.IsModeDevelopment() {
			addr := "https://localhost" + w.httpsServer.Addr
			fmt.Println("Listening on", addr)
		}

		err := w.httpsServer.ListenAndServeTLS("", "")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("HTTPS server error: %w", err)
			return
		}
		errChan <- nil
	}()

	go func() {
		if env.IsDebugInfoEnabled() && env.IsModeDevelopment() {
			addr := "http://localhost" + w.httpServer.Addr
			fmt.Println("Listening on", addr)
		}

		err := w.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
			return
		}
		errChan <- nil
	}()

	return <-errChan
}

func (w *webserver) Close() error {
	errClosingHttps := w.httpsServer.Close()
	errClosingHttp := w.httpServer.Close()

	return errors.Join(errClosingHttp, errClosingHttps)
}
