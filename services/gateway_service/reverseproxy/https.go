package reverseproxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"gateway_service/certificate"
)

type HTTPSServer struct {
	Port string
	// SelfSignedCert og SelfSignedKey ignoreres i prod-only opsætning
	SelfSignedCert string
	SelfSignedKey  string
	ACME           certificate.AutoCertManager
	Handler        http.Handler
	//PreferHTTP2Only bool // true => NextProtos: ["h2"]
}

func NewHTTPSServer(server HTTPSServer) (*http.Server, error) {
	if server.Port == "" {
		server.Port = "443"
	}
	if server.Handler == nil {
		server.Handler = http.DefaultServeMux
	}
	if server.ACME == nil {
		return nil, fmt.Errorf("ACME manager is required for production certificates")
	}

	srv := &http.Server{
		Addr:    ":" + server.Port,
		Handler: server.Handler,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.TLSConfig.NextProtos = []string{"h2", "http/1.1"}

	// Prod-only: altid brug ACME-certifikater
	getCert := server.ACME.GetCertificate
	srv.TLSConfig.GetCertificate = func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return getCert(chi)
	}

	return srv, nil
}
