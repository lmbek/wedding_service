// READ THIS FIRST:
// create the certificate and key files for localhost first!
// You need to generate self-signed certificates
// by running go generate

//go:generate go run ./self_sign_cert/self_sign_cert.go

//////

package certificate

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"math"
	"net/http"
	"os"
	"time"
)

func getLocalhostCertAndKey(crtPath string, keyPath string) (cert []byte, key []byte, err error) {
	cert, err = os.ReadFile(crtPath)
	if err != nil {
		return nil, nil, fmt.Errorf("(crt file) please run go generate on the certificate folder first, could not read file: %v", err)
	}

	key, err = os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("(key file) please run go generate on the certificate folder first, could not read file: %v", err)
	}
	return cert, key, nil
}

// UseLOCALHOST loads the certificate and private key as a TLS certificate.
func UseLOCALHOST() (tlsCert *tls.Certificate, err error) {
	cert, key, err := getLocalhostCertAndKey(os.Getenv("LOCALHOST_CERT"), os.Getenv("LOCALHOST_KEY"))
	if err != nil {
		return nil, err
	}

	x509KeyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %v", err)
	}
	return &x509KeyPair, nil
}

// UseACME initializes the autocert.Manager for managing certificates.
func UseACME() *autocert.Manager {
	domainAliases := map[string][]string{
		os.Getenv("WEDDING_SERVICE_EXTERNAL_HOSTNAME"): {os.Getenv("WEDDING_SERVICE_EXTERNAL_HOSTNAME_ALIAS1")},
	}

	return &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("certs"),
		HostPolicy: func(ctx context.Context, host string) error {
			// Validate host against configured domains and aliases
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
		},
		Client: &acme.Client{
			DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory", // zerossl: "https://acme.zerossl.com/v2/DV90"
			RetryBackoff: func(n int, r *http.Request, resp *http.Response) time.Duration {
				// Exponential backoff: Retry after 2^n seconds, where n is the attempt number
				return time.Duration(math.Pow(2, float64(n))) * time.Second
			},
		},
	}
}
