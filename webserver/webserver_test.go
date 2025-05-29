package webserver

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"
	"wedding_service/certificate"
)

func TestStart(t *testing.T) {
	err := Start()
	if err != nil {
		t.Errorf("got error: %s\n", err)
	}
}

func TestListenAndServe(t *testing.T) {
	httpServer := http.Server{Addr: "localhost:8080"}
	httpsServer := http.Server{Addr: "localhost:8443"}
	autoCertManager := certificate.UseACME()
	httpsServer.TLSConfig = &tls.Config{
		MinVersion:     tls.VersionTLS12, // HTTP/2 requires TLS 1.2 or higher
		NextProtos:     []string{"h2"},   // Enforce HTTP/2
		GetCertificate: autoCertManager.GetCertificate,
	}

	go func() {
		time.Sleep(2 * time.Second)
		httpServer.Close()
		httpsServer.Close()
	}()

	err := ListenAndServe(&httpServer, &httpsServer)
	if err != nil {
		t.Errorf("got error: %s\n", err)
	}
}
