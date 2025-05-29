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

// Embed the certificate and key files
// NOTE: you need to generate self-signed certificates and put them in this package (certificate)

//go:generate go run ./self_sign_cert/self_sign_cert.go
//go:embed localhost_wedding_service.crt
var embeddedCert []byte

//go:embed localhost_wedding_service.key
var embeddedKey []byte

// UseLOCALHOST loads the embedded certificate and private key as a TLS certificate.
func UseLOCALHOST() (*tls.Certificate, error) {
	cert, err := tls.X509KeyPair(embeddedCert, embeddedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded TLS certificate: %w", err)
	}
	return &cert, nil
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
