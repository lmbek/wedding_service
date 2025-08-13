package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"wedding_service/certificate"
	"wedding_service/config"
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

func NewWebserver(config config.Config, frontend frontend.Frontend) (w Webserver, err error) {
	mux := http.NewServeMux()
	useWebsite(config, mux, frontend)
	useApi(config, mux)

	acmeManager, err := certificate.InitAcme()
	if err != nil {
		return nil, fmt.Errorf("could not use acme manager: %w", err)
	}

	httpServer := newHttpServer(config.HttpPort())
	httpsServer, err := newHttpsServer(config.HttpsPort(), acmeManager)
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
	// Listen and serve HTTPS / TLS
	err := w.httpsServer.ListenAndServeTLS("", "")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTPS server error: %w", err)
	}

	return nil
}

func (w *webserver) listenHTTP() error {
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
