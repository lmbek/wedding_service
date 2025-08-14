package certificate

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"math"
	"net/http"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func loadTLSKeyPair(cert []byte, key []byte) (*tls.Certificate, error) {
	x509KeyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}
	return &x509KeyPair, nil
}

// InitAcme initializes the autocert.Manager for managing certificates.
func InitAcme() (acmeManager *autocert.Manager, err error) {
	return &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"),
		HostPolicy: autocert.HostWhitelist("lmbek.dk", "www.lmbek.dk"),
		Client: &acme.Client{
			DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory", // zerossl: "https://acme.zerossl.com/v2/DV90"
			RetryBackoff: handleRetryBackoff,
		},
	}, nil
}

// Validate host against configured domains and aliases
func handleHostPolicy(domainAliases map[string][]string) autocert.HostPolicy {
	return func(ctx context.Context, host string) error {
		for domain, aliases := range domainAliases {
			if host == domain {
				return nil
			}
			for _, alias := range aliases {
				if host == alias {
					return nil
				}
			}
		}
		return fmt.Errorf("unauthorized host: %s", host)
	}
}

// handleRetryBackoff handles exponential backoff: Retry after 2^n seconds, where n is the attempt number
func handleRetryBackoff(n int, r *http.Request, resp *http.Response) time.Duration {
	return time.Duration(math.Pow(2, float64(n))) * time.Second
}
