//go:build self_sign_cert

// NOTE: this file is ignored by the default build
// as it is a tool, we don't count this file in our test coverage total

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
	"time"
)

func main() {
	createSelfSignedCertificateFile()
}

// CreateSelfSignedCertificateFile generates a self-signed certificate and saves it to files with a timestamp
func createSelfSignedCertificateFile() {
	// Generate the certificate
	cert, privKeyPEM, err := generateSelfSignedCertificatePEM()
	if err != nil {
		fmt.Printf("Error generating self-signed certificate: %v\n", err)
		return
	}

	// Create a timestamp for the file names
	//timestamp := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS

	// Generate file names with the timestamp
	//certFileName := fmt.Sprintf("localhost_%s.crt", timestamp)
	//keyFileName := fmt.Sprintf("localhost_%s.key", timestamp)
	certFileName := fmt.Sprintf("localhost_%s.crt", "wedding_service")
	keyFileName := fmt.Sprintf("localhost_%s.key", "wedding_service")

	// Save the certificate to a file
	err = os.WriteFile(certFileName, cert, 0644)
	if err != nil {
		fmt.Printf("Error writing certificate to file: %v\n", err)
		return
	}

	// Save the private key to a file
	err = os.WriteFile(keyFileName, privKeyPEM, 0600) // Use more restrictive permissions
	if err != nil {
		fmt.Printf("Error writing private key to file: %v\n", err)
		return
	}

	fmt.Printf("Certificate and private key have been written to '%s' and '%s'\n", certFileName, keyFileName)
}

// generateSelfSignedCertificatePEM generates a self-signed certificate and private key as PEM-encoded bytes
func generateSelfSignedCertificatePEM() ([]byte, []byte, error) {
	// Generate private key (RSA in this case)
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Valid for 1 year

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

	// Self-sign the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// PEM encode the certificate and private key
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certPEM, keyPEM, nil
}
