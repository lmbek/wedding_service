//go:build self_sign_cert

// NOTE: this test is ignored by the default tests
// the test coverage of this tool is 67.9%
// as it is a tool, we don't count this test coverage in our total

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
