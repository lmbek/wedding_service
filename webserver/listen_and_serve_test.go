package webserver

import (
	"crypto/tls"
	"errors"
	"net/http"
	"testing"
	"time"
	"wedding_service/certificate"
)

func TestListenAndServe(t *testing.T) {
	httpServer := http.Server{Addr: "localhost:8080"}
	httpsServer := http.Server{Addr: "localhost:8443"}
	autoCertManager := certificate.UseACME()
	httpsServer.TLSConfig = &tls.Config{
		MinVersion:     tls.VersionTLS12, // HTTP/2 requires TLS 1.2 or higher
		NextProtos:     []string{"h2"},   // Enforce HTTP/2
		GetCertificate: autoCertManager.GetCertificate,
	}

	// TODO - use synctest
	go func() {
		time.Sleep(2 * time.Second)
		httpServer.Close()
		httpsServer.Close()
	}()

	err := ListenAndServe(&httpServer, &httpsServer)
	if err != nil {
		t.Errorf("could not ListenAndServe: %s\n", err)
	}
}

func TestListenAndServe_Error(t *testing.T) {
	// Start a dummy server on localhost:8080 to block that port
	blocker := &http.Server{Addr: "localhost:8080"}
	go blocker.ListenAndServe()
	defer blocker.Close()
	time.Sleep(100 * time.Millisecond) // Let blocker start

	httpServer := http.Server{Addr: "localhost:8080"}  // same port, will error
	httpsServer := http.Server{Addr: "localhost:8443"} // unused port

	err := ListenAndServe(&httpServer, &httpsServer)

	if err == nil {
		t.Errorf("expected error due to port in use, got nil")
	} else if errors.Is(err, http.ErrServerClosed) {
		t.Errorf("expected real error, not http.ErrServerClosed")
	}
}
