//go:build self_sign_cert

// NOTE: this file is ignored by the default build
// as it is a tool, we don't count this file in our test coverage total

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	// TEST SERVER CAN BE USED, BUT BE AWARE THAT EMBEDDING USES PREVIOUS LOCALHOST CERT,
	// SO THE CREATE SELFSIGNED AND TESTSERVER CANT RUN AT THE SAME TIME!

	// FIRST FLIP THE IF STATEMENT
	// SECOND UNCOMMENT THE EMBEDS INSIDE certificate_windows.go
	// if you don't uncomment the embeds this err will come:
	// selfsigning failed: failed to load embedded TLS certificate and key: tls: failed to find any PEM data in certificate input

	if true {
		certName := "localhost"
		orgNames := []string{"Local Beksoft Cert"}
		dnsNames := []string{"localhost"}

		createSelfSignedCertificateFile(certName, orgNames, dnsNames)
		return
	}

	// Test HTTPS server setup
	//testServer()
}

// CreateSelfSignedCertificateFile generates a self-signed certificate and saves it to files with a timestamp
func createSelfSignedCertificateFile(certName string, orgNames []string, dnsNames []string) {
	if certName == "" {
		certName = "localhost"
	}

	// Generate the certificate and private key first
	cert, privKey, err := generateSelfSignedCertificate(orgNames, dnsNames)
	if err != nil {
		fmt.Printf("Error generating self-signed certificate: %v\n", err)
		return
	}

	// Create the certificate directory if it doesn't exist
	oldCertDir := "old_certs"
	os.MkdirAll(oldCertDir, 0700) // Restrictive permissions for directories

	// Move existing certificate to the old folder
	err = moveOldCertificates(".", oldCertDir)
	if err != nil {
		fmt.Printf("Error moving old certificates: %v\n", err)
		return
	}

	// Create a timestamp for the file names
	//timestamp := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS

	// Generate file names with the timestamp
	//certFileName := filepath.Join(certDir, fmt.Sprintf("%s_%s.crt", certName, timestamp))
	//keyFileName := filepath.Join(certDir, fmt.Sprintf("%s_%s.key", certName, timestamp))
	certFileName := fmt.Sprintf("localhost_%s.crt", "wedding_service")
	keyFileName := fmt.Sprintf("localhost_%s.key", "wedding_service")

	// Save the timestamped certificate and private key to files
	err = os.WriteFile(certFileName, cert, 0644) // Public read, owner write
	if err != nil {
		fmt.Printf("Error writing certificate to file: %v\n", err)
		return
	}

	err = os.WriteFile(keyFileName, privKey, 0600) // Owner read/write only
	if err != nil {
		fmt.Printf("Error writing private key to file: %v\n", err)
		return
	}

	// Now create the non-timestamped version with the same certificate bytes
	defaultCertFileName := fmt.Sprintf("localhost_%s.crt", "wedding_service")
	defaultKeyFileName := fmt.Sprintf("localhost_%s.key", "wedding_service")

	// Save the non-timestamped certificate and private key with the same bytes
	err = os.WriteFile(defaultCertFileName, cert, 0644) // Public read, owner write
	if err != nil {
		fmt.Printf("Error writing non-timestamped certificate to file: %v\n", err)
		return
	}

	err = os.WriteFile(defaultKeyFileName, privKey, 0600) // Owner read/write only
	if err != nil {
		fmt.Printf("Error writing non-timestamped private key to file: %v\n", err)
		return
	}

	fmt.Printf("Certificate and private key have been written to '%s' and '%s'\n", certFileName, keyFileName)

	// Add the certificate to the Windows CA store
	err = AddCertificateToWindowsCAStore(certName, orgNames, cert)
	if err != nil {
		fmt.Printf("Error adding certificate to Windows CA store: %v\n", err)
		return
	}

	fmt.Println("Certificate successfully added to the Windows CA store.")

	// Now load the certificate and key into a TLS config
	err = loadTLSConfig(certFileName, keyFileName)
	if err != nil {
		fmt.Printf("Error loading TLS config: %v\n", err)
	}
}

// generateSelfSignedCertificate generates a self-signed certificate and private key
func generateSelfSignedCertificate(orgNames []string, dnsNames []string) ([]byte, []byte, error) {
	// Generate ECDSA private key (P-256 curve)
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate ECDSA private key: %w", err)
	}

	// Create certificate template
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Valid for 1 year
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: orgNames,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
	}

	// Self-sign the certificate using ECDSA private key
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	key, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	// PEM encode the certificate and private key
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: key})

	return certPEM, keyPEM, nil
}

// loadTLSConfig loads the certificate and private key into a TLS configuration
func loadTLSConfig(certFileName, keyFileName string) error {
	// Load the certificate
	certPEM, err := os.ReadFile(certFileName)
	if err != nil {
		return fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Load the private key
	keyPEM, err := os.ReadFile(keyFileName)
	if err != nil {
		return fmt.Errorf("failed to read key file: %w", err)
	}

	// Parse the certificate and key
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return fmt.Errorf("failed to parse certificate and key: %w", err)
	}

	// Create a TLS config with the loaded certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Use the TLS config (e.g., with a server or client)
	fmt.Println("TLS configuration loaded successfully.")
	_ = tlsConfig // To prevent unused variable warning

	return nil
}

// AddCertificateToWindowsCAStore adds the provided certificate to the Current User's CA store (no admin required)
func AddCertificateToWindowsCAStore(certName string, orgNames []string, certPEM []byte) error {
	// remove already existing PEM's already exported with same name
	err := RemoveCertificateFromCurrentUser(orgNames)
	if err != nil {
		return err
	}

	// Define the certificate file path inside the certificate directory with a clearer name
	timestamp := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
	certFileName := fmt.Sprintf("%s_ca_%s.pem", certName, timestamp)
	certFilePath := certFileName

	// Save the certificate directly to the certificate folder
	err = os.WriteFile(certFilePath, certPEM, 0644) // Public read, owner write
	if err != nil {
		return err
	}

	// Use PowerShell to add the certificate to the Current User's CA store (no admin required)
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(
		`Import-Certificate -FilePath '%s' -CertStoreLocation Cert:\\CurrentUser\\Root`,
		certFilePath,
	))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add certificate to Windows CA store: %w, output: %s", err, string(output))
	}

	return nil
}

// RemoveCertificateFromCurrentUser removes a certificate from the CurrentUser CA store
func RemoveCertificateFromCurrentUser(orgNames []string) error {
	for _, org := range orgNames {
		// Use PowerShell to find and remove the certificate with the given subject from the CurrentUser's CA store
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf(`$subjectName = "%s"
				$certificates = Get-ChildItem -Path Cert:\CurrentUser\Root | Where-Object {$_.Subject -like "*$subjectName*"}
				foreach ($cert in $certificates) {
					Remove-Item -Path $cert.PSPath -Force
				}`,
				org,
			),
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to remove certificate for '%s' from CurrentUser CA store: %w, output: %s", org, err, string(output))
		}
	}
	return nil
}

// moveOldCertificates moves any existing certificates to the 'old' folder in the certificates directory.
// It only moves .crt, .key, and .pem files. If the 'old' folder has more than 100 files, it removes the oldest.
func moveOldCertificates(certDir, oldCertDir string) error {
	// Check if the certificate directory exists
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		// If the directory doesn't exist, return early
		return nil
	}

	// List files in the certificate directory
	files, err := os.ReadDir(certDir)
	if err != nil {
		return fmt.Errorf("failed to read certificate directory: %w", err)
	}

	// Move each file to the 'old' directory
	for _, file := range files {
		// Skip directories, we only want to move files
		if file.IsDir() {
			continue
		}

		// Get the file extension
		ext := strings.ToLower(filepath.Ext(file.Name()))

		// Only move .crt, .key, and .pem files
		if ext == ".crt" || ext == ".key" || ext == ".pem" {
			// Construct the source and destination file paths
			sourceFilePath := filepath.Join(certDir, file.Name())
			destFilePath := filepath.Join(oldCertDir, file.Name())

			// Move the file
			err := os.Rename(sourceFilePath, destFilePath)
			if err != nil {
				return fmt.Errorf("failed to move file '%s' to '%s': %w", sourceFilePath, destFilePath, err)
			}
		}
	}

	// Check if the 'old' folder has more than 100 files
	oldFiles, err := os.ReadDir(oldCertDir)
	if err != nil {
		return fmt.Errorf("failed to read old certificate directory: %w", err)
	}

	if len(oldFiles) > 100 {
		// Sort files by creation time (oldest first)
		sort.Slice(oldFiles, func(i, j int) bool {
			fileInfoI, errI := oldFiles[i].Info()
			fileInfoJ, errJ := oldFiles[j].Info()
			if errI != nil || errJ != nil {
				// If there's an error getting file info, treat them as equal
				return false
			}
			return fileInfoI.ModTime().Before(fileInfoJ.ModTime())
		})

		// Remove the oldest file
		oldestFile := oldFiles[0]
		oldestFilePath := filepath.Join(oldCertDir, oldestFile.Name())

		// Delete the oldest file
		err := os.Remove(oldestFilePath)
		if err != nil {
			return fmt.Errorf("failed to delete oldest file '%s': %w", oldestFilePath, err)
		}
		log.Printf("Deleted the oldest file: %s", oldestFilePath)
	}

	return nil
}
