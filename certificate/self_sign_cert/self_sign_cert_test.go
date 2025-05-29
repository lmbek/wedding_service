//go:build self_sign_cert

package main

import (
	"testing"
)

func Test_createSelfSignedCertificateFile(t *testing.T) {
	createSelfSignedCertificateFile()
}

func Test_generateSelfSignedCertificatePEM(t *testing.T) {
	cert, privKeyPEM, err := generateSelfSignedCertificatePEM()
	if err != nil {
		t.Errorf("Got err: %v", err)
	}

	if cert == nil {
		t.Errorf("cert should not be nil")
	}

	if privKeyPEM == nil {
		t.Errorf("privKeyPEM should not be nil")
	}
}
