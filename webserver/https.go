package webserver

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"wedding_service/certificate"
	"wedding_service/env"
)

func newHttpsServer(port string, acmeManager *autocert.Manager) (*http.Server, error) {
	server := &http.Server{}

	server.Addr = fmt.Sprintf(":%s", port)
	server.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS13, // HTTP/2 requires TLS 1.3 or higher
		NextProtos: []string{"h2"},   // Enforce HTTP/2
	}

	if env.IsModeDevelopment() {
		server.Addr = fmt.Sprintf("localhost%s", server.Addr)

		cert, err := certificate.LoadSelfSigned(env.Env.CertPath, env.Env.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("could not use localhost certificate: %w", err)
		}
		server.TLSConfig.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return cert, nil
		}
	}

	if env.IsModeProduction() {
		server.TLSConfig.GetCertificate = acmeManager.GetCertificate
	}

	return server, nil
}
