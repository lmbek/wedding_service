package webserver

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func newHttpsServer(port string) *http.Server {
	return &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // HTTP/2 requires TLS 1.2 or higher
			NextProtos: []string{"h2"},   // Enforce HTTP/2
		},
	}
}
