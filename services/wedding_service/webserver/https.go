package webserver

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func newHttpsServer(port string, acmeManager *autocert.Manager) (*http.Server, error) {
	server := &http.Server{}

	server.Addr = fmt.Sprintf(":%s", port)
	server.TLSConfig = &tls.Config{
		MinVersion:     tls.VersionTLS12,           // allow TLS 1.2 and 1.3
		NextProtos:     []string{"h2", "http/1.1"}, // allow HTTP/2 and HTTP/1.1
		GetCertificate: acmeManager.GetCertificate,
	}

	server.TLSConfig.GetCertificate = acmeManager.GetCertificate
	return server, nil
}
