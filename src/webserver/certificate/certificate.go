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
	"wedding_service/config"
)

type AutoCertManager interface {
	init(config config.Config) (acmeManager AutoCertManager, err error)
	HTTPHandler(fallback http.Handler) http.Handler
	LoadSelfSigned(certPath string, keyPath string) (tlsCert *tls.Certificate, err error)
	GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error)
}

type autoCertManager struct {
	config   config.Config
	autoCert *autocert.Manager
}

func NewAutoCertManager(config config.Config) (AutoCertManager, error) {
	a := &autoCertManager{}
	return a.init(config)
}

// InitAcme initializes the autocert.Manager for managing certificates.
func (a *autoCertManager) init(config config.Config) (acmeManager AutoCertManager, err error) {
	hostnames := config.Hostnames()
	if hostnames == nil || len(config.Hostnames()) == 0 {
		return nil, errors.New("hostnames must not be empty, should have format: domain:alias,alias2|domain2:alias3")
	}
	a.autoCert = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"),
		HostPolicy: a.handleHostPolicy(hostnames),
		Client: &acme.Client{
			DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory", // zerossl: "https://acme.zerossl.com/v2/DV90"
			RetryBackoff: a.handleRetryBackoff,
		},
	}

	return a, nil
}

func (a *autoCertManager) getSelfSignedCertAndKey(certPath string, keyPath string) (cert []byte, key []byte, err error) {
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

// LoadSelfSigned loads the certificate and private key as a TLS certificate.
func (a *autoCertManager) LoadSelfSigned(certPath string, keyPath string) (tlsCert *tls.Certificate, err error) {
	cert, key, err := a.getSelfSignedCertAndKey(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("could not GetSelfSignedCertAndKey: %w", err)
	}

	tlsCert, err = a.loadTLSKeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("could not loadTLSKeyPair: %w", err)
	}

	return
}

func (a *autoCertManager) loadTLSKeyPair(cert []byte, key []byte) (*tls.Certificate, error) {
	x509KeyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}
	return &x509KeyPair, nil
}

// Validate host against configured domains and aliases
func (a *autoCertManager) handleHostPolicy(domainAliases map[string][]string) autocert.HostPolicy {
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
func (a *autoCertManager) handleRetryBackoff(n int, r *http.Request, resp *http.Response) time.Duration {
	return time.Duration(math.Pow(2, float64(n))) * time.Second
}

func (a *autoCertManager) HTTPHandler(fallback http.Handler) http.Handler {
	return a.autoCert.HTTPHandler(fallback)
}

func (a *autoCertManager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return a.autoCert.GetCertificate(hello)
}
