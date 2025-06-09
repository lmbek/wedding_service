// READ THIS FIRST:
// create the certificate and key files for localhost first!
// You need to generate self-signed certificates
// by running go generate

//// go:generate go run ./self_sign_cert/self_sign_cert.go
//go:generate go run ./self_sign_cert_windows/self_sign_cert_windows.go

//////

package certificate

import (
	"context"
	"crypto/tls"
	_ "embed"
	"errors"
	"fmt"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"math"
	"net/http"
	"os"
	"time"
	"wedding_service/env"
)

func getSelfSignedCertAndKey(certPath string, keyPath string) (cert []byte, key []byte, err error) {
	cert, err = os.ReadFile(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("(cert file) please run go generate first: %w", err)
	}

	key, err = os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("(key file) please run go generate first: %w", err)
	}

	return cert, key, nil
}

// UseSelfSigned loads the certificate and private key as a TLS certificate.
func UseSelfSigned(certPath string, keyPath string) (tlsCert *tls.Certificate, err error) {
	cert, key, err := getSelfSignedCertAndKey(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("could not GetSelfSignedCertAndKey: %w", err)
	}

	tlsCert, err = loadTLSKeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("could not loadTLSKeyPair: %w", err)
	}

	return
}

func loadTLSKeyPair(cert []byte, key []byte) (*tls.Certificate, error) {
	x509KeyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}
	return &x509KeyPair, nil
}

// UseAcme initializes the autocert.Manager for managing certificates.
func UseAcme() (acmeManager *autocert.Manager, err error) {
	hostnames := env.Env.Hostnames
	if hostnames == nil && len(hostnames) >= 1 {
		return nil, errors.New("hostnames must not be empty, should have format: domain:alias,alias2|domain2:alias3")
	}

	return &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"),
		HostPolicy: handleHostPolicy(hostnames),
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
