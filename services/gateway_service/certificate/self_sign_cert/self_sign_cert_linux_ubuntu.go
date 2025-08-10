//go:build linux && self_sign_cert

// TODO: the self_sign_cert ignore is for test coverage ignore. Tests can be written to remove this

// READ THIS FIRST:
// create the certificate and key files for localhost first!
// You need to generate self-signed certificates
// by running go generate

//// go:generate go run ./self_sign_cert/self_sign_cert_linux_ubuntu.go
//go:generate go run self_sign_cert_linux_ubuntu.go

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("Self sign certificate")

	certPEM, certFileName, err := generateSelfSignedCert()
	if err != nil {
		fmt.Printf("Failed to generate cert: %v\n", err)
		return
	}

	err = importCertToLinuxTrustStore(certPEM, certFileName)
	if err != nil {
		fmt.Printf("Failed to import cert to trust store: %v\n", err)
		return
	}

	fmt.Println("Certificate successfully imported to system trust store")
}

// generateSelfSignedCert generates the self-signed certificate and writes cert, key, and pem files.
// It returns the PEM certificate bytes and the certificate filename (for import).
func generateSelfSignedCert() ([]byte, string, error) {
	cert, key, _, err := generateSelfSignedPEM()
	if err != nil {
		return nil, "", fmt.Errorf("error generating self-signed certificate: %w", err)
	}

	baseName := "localhost_wedding_service"
	certFileName := baseName + ".crt"
	keyFileName := baseName + ".key"
	pemFileName := baseName + ".pem"

	if err := os.WriteFile(certFileName, cert, 0644); err != nil {
		return nil, "", fmt.Errorf("error writing certificate file: %w", err)
	}

	if err := os.WriteFile(keyFileName, key, 0644); err != nil {
		return nil, "", fmt.Errorf("error writing private key file: %w", err)
	}

	// Combine cert + key for pem file (optional)
	pemContent := append(cert, key...)
	if err := os.WriteFile(pemFileName, pemContent, 0644); err != nil {
		return nil, "", fmt.Errorf("error writing pem file: %w", err)
	}

	fmt.Printf("Certificate and key saved as '%s', '%s', and combined PEM '%s'\n", certFileName, keyFileName, pemFileName)
	return cert, certFileName, nil
}

// generateSelfSignedPEM creates a self-signed cert and returns cert, key PEM bytes plus raw certificate bytes.
func generateSelfSignedPEM() (cert []byte, key []byte, rawCert []byte, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Local Development"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	rawCert, err = x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawCert})
	key = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return cert, key, rawCert, nil
}

// importCertToLinuxTrustStore writes the cert to the system cert dir and runs update-ca-certificates.
func importCertToLinuxTrustStore(certPEM []byte, certName string) error {
	if len(certPEM) == 0 {
		return fmt.Errorf("certificate PEM data is empty")
	}
	if certName == "" {
		return fmt.Errorf("certificate name cannot be empty")
	}

	if filepath.Ext(certName) != ".crt" {
		certName += ".crt"
	}

	destPath := filepath.Join("/usr/local/share/ca-certificates", certName)

	// Remove old certificate if it exists
	if _, err := os.Stat(destPath); err == nil {
		fmt.Printf("Old certificate found at %s, deleting it...\n", destPath)
		if err := os.Remove(destPath); err != nil {
			return fmt.Errorf("failed to remove old certificate: %w", err)
		}
	}

	// Write new certificate
	if err := os.WriteFile(destPath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}

	// Update system CA certificates
	cmd := exec.Command("update-ca-certificates")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update CA certificates: %v, output: %s", err, out)
	}

	fmt.Println("Certificate successfully imported into system trust store.")
	return nil
}
