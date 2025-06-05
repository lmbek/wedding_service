package webserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"wedding_service/certificate"
)

func Start(httpServer *http.Server, httpsServer *http.Server, certPath string, keyPath string) error {
	httpsServer.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS12, // HTTP/2 requires TLS 1.2 or higher
		NextProtos: []string{"h2"},   // Enforce HTTP/2
	}

	if os.Getenv("MODE") == "development" {
		cert, err := certificate.UseLOCALHOST(certPath, keyPath)
		if err != nil {
			return fmt.Errorf("could not use localhost certificate: %v\n", err)
		}
		httpsServer.TLSConfig.Certificates = []tls.Certificate{*cert}
	}

	if os.Getenv("MODE") == "production" {
		autoCertManager := certificate.UseACME()
		httpsServer.TLSConfig.GetCertificate = autoCertManager.GetCertificate
	}

	return ListenAndServe(httpServer, httpsServer)
}

func ListenAndServe(httpServer *http.Server, httpsServer *http.Server) error {
	var wg sync.WaitGroup
	wg.Add(2)

	var httpErr, httpsErr error

	go func() {
		defer wg.Done()
		fmt.Println("Listening on http://" + httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("HTTP server error: %s\n", err)
			httpErr = err
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Listening on https://" + httpsServer.Addr)
		err := httpsServer.ListenAndServeTLS("", "")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("HTTPS server error: %s\n", err)
			httpsErr = err
		}
	}()

	wg.Wait()

	if httpErr != nil {
		return httpErr
	}
	if httpsErr != nil {
		return httpsErr
	}
	return nil
}
