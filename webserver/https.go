package webserver

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"wedding_service/env"
)

func newHttpsServer(port string) *http.Server {
	var addr string

	if env.IsModeDevelopment() {
		addr = fmt.Sprintf("localhost:%s", port)
	} else {
		addr = fmt.Sprintf(":%s", port)
	}

	return &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // HTTP/2 requires TLS 1.2 or higher
			NextProtos: []string{"h2"},   // Enforce HTTP/2
		},
	}
}
