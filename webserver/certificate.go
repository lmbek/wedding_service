package webserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"wedding_service/certificate"
	"wedding_service/env"
)

func useCertificate(httpsServer *http.Server) error {
	if env.IsModeDevelopment() {
		cert, err := certificate.UseSelfSigned(env.Env.CertPath, env.Env.KeyPath)
		if err != nil {
			return fmt.Errorf("could not use localhost certificate: %w", err)
		}
		httpsServer.TLSConfig.GetCertificate = wrapCert(cert)
	}

	if env.IsModeProduction() {
		acmeManager, err := certificate.UseAcme()
		if err != nil {
			return fmt.Errorf("could not use acme manager: %w", err)
		}
		httpsServer.TLSConfig.GetCertificate = acmeManager.GetCertificate
	}

	if env.IsModeNotSet() {
		return errors.New("no MODE set in .env")
	}

	return nil
}

func wrapCert(cert *tls.Certificate) func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return cert, nil
	}
}
