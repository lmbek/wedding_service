package webserver

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"wedding_service/certificate"
	"wedding_service/helper"
)

func setupHttpsServer(httpsServer *http.Server, certPath string, keyPath string) error {
	httpsServer.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS12, // HTTP/2 requires TLS 1.2 or higher
		NextProtos: []string{"h2"},   // Enforce HTTP/2
	}

	err := setupCertificate(httpsServer, certPath, keyPath)
	if err != nil {
		return err
	}

	return nil
}

func setupCertificate(httpsServer *http.Server, certPath string, keyPath string) error {
	envMode := os.Getenv("MODE")
	mode, err := helper.ParseMode(envMode)
	if err != nil {
		return err
	}

	if mode == helper.Development {
		cert, err := certificate.UseLOCALHOST(certPath, keyPath)
		if err != nil {
			return fmt.Errorf("could not use localhost certificate: %v\n", err)
		}
		httpsServer.TLSConfig.Certificates = []tls.Certificate{*cert}
	}

	if mode == helper.Production {
		autoCertManager := certificate.UseACME()
		httpsServer.TLSConfig.GetCertificate = autoCertManager.GetCertificate
	}

	return nil
}
