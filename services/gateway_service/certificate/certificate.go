package certificate

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gateway_service/config"
	"math"
	"net/http"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type AutoCertManager interface {
	init(config config.Config) (acmeManager AutoCertManager, err error)
	HTTPHandler(fallback http.Handler) http.Handler
	GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error)
}

type autoCertManager struct {
	config   config.Config
	autoCert *autocert.Manager
}

func NewAutoCertManager(cfg config.Config) (AutoCertManager, error) {
	a := &autoCertManager{}
	return a.init(cfg)
}

func (a *autoCertManager) init(cfg config.Config) (AutoCertManager, error) {
	hostnames := cfg.Hostnames()
	if hostnames == nil || len(hostnames) == 0 {
		return nil, errors.New("hostnames must not be empty, should have format: domain:alias,alias2|domain2:alias3")
	}

	// Flad alle domæner + aliaser til én whitelist
	allowed := make([]string, 0, 8)
	for domain, aliases := range hostnames {
		if domain != "" {
			allowed = append(allowed, domain)
		}
		for _, a := range aliases {
			if a != "" {
				allowed = append(allowed, a)
			}
		}
	}
	if len(allowed) == 0 {
		return nil, fmt.Errorf("no valid hostnames collected for ACME whitelist")
	}

	a.config = cfg
	a.autoCert = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"), // sørg for persistent/skrivebar sti
		HostPolicy: autocert.HostWhitelist(allowed...),
		Client: &acme.Client{
			DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory",
			RetryBackoff: a.handleRetryBackoff,
		},
	}
	return a, nil
}

func (a *autoCertManager) handleRetryBackoff(n int, r *http.Request, resp *http.Response) time.Duration {
	return time.Duration(math.Pow(2, float64(n))) * time.Second
}

func (a *autoCertManager) HTTPHandler(fallback http.Handler) http.Handler {
	return a.autoCert.HTTPHandler(fallback)
}

func (a *autoCertManager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return a.autoCert.GetCertificate(hello)
}
