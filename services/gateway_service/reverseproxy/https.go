package reverseproxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"gateway_service/certificate"
)

type HTTPSServer struct {
	Port            string
	SelfSignedCert  string
	SelfSignedKey   string
	ACME            certificate.AutoCertManager
	Handler         http.Handler
	PreferHTTP2Only bool // true => NextProtos: ["h2"]
}

func NewHTTPSServer(server HTTPSServer) (*http.Server, error) {
	if server.Port == "" {
		server.Port = "443"
	}
	if server.Handler == nil {
		server.Handler = http.DefaultServeMux
	}

	srv := &http.Server{
		Addr:    ":" + server.Port,
		Handler: server.Handler,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// HTTP/2 enforcement hvis ønsket
	if server.PreferHTTP2Only {
		srv.TLSConfig.NextProtos = []string{"h2"}
	} else {
		srv.TLSConfig.NextProtos = []string{"h2", "http/1.1"}
	}

	// Prefer self-signed in local/dev when provided; otherwise fall back to ACME manager
	if server.SelfSignedCert != "" && server.SelfSignedKey != "" {
		if server.ACME != nil {
			cert, err := server.ACME.LoadSelfSigned(server.SelfSignedCert, server.SelfSignedKey)
			if err != nil {
				return nil, fmt.Errorf("load self-signed: %w", err)
			}
			srv.TLSConfig.GetCertificate = func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return cert, nil
			}
		}
	} else if server.ACME != nil {
		// Use real certificates managed by ACME when no self-signed is specified
		getCert := server.ACME.GetCertificate
		srv.TLSConfig.GetCertificate = func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return getCert(chi)
		}
	}

	return srv, nil
}
