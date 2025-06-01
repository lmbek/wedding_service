package webserver

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"
)

// helper function to restore env vars
func setEnvAndRestore(key, value string) func() {
	old := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, old)
	}
}

func TestSetupCertificate_InvalidMode(t *testing.T) {
	restore := setEnvAndRestore("MODE", "invalid-mode")
	defer restore()

	server := &http.Server{TLSConfig: &tls.Config{}}
	err := setupCertificate(server, "", "")
	if err == nil {
		t.Errorf("Expected error for invalid mode, got nil")
	}
}

func TestSetupCertificate_DevelopmentMode(t *testing.T) {
	restore := setEnvAndRestore("MODE", "development")
	defer restore()

	server := &http.Server{TLSConfig: &tls.Config{}}
	err := setupCertificate(server, "../certificate/localhost_wedding_service.crt", "../certificate/localhost_wedding_service.key")
	if err != nil {
		t.Errorf("Expected no error in development mode, got: %v", err)
	}

	if len(server.TLSConfig.Certificates) == 0 {
		t.Errorf("Expected certificate to be set in development mode")
	}
}

func TestSetupCertificate_ProductionMode(t *testing.T) {
	restore := setEnvAndRestore("MODE", "production")
	defer restore()

	server := &http.Server{TLSConfig: &tls.Config{}}
	err := setupCertificate(server, "", "")
	if err != nil {
		t.Errorf("Expected no error in production mode, got: %v", err)
	}

	if server.TLSConfig.GetCertificate == nil {
		t.Errorf("Expected GetCertificate to be set in production mode")
	}
}
