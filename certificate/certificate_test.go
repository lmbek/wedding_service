// Note: in order to test this, we need to have certificates generated first (see top of certificate.go),
// then run the test with the tag certificate

package certificate

import (
	"os"
	"strings"
	"testing"
	"time"
	"wedding_service/env"
)

func Test_getSelfSignedCertAndKey(t *testing.T) {
	certContent := []byte("dummy cert content")
	keyContent := []byte("dummy key content")

	certPath := createTempFile(t, "cert-*.crt", certContent)
	defer os.Remove(certPath)

	keyPath := createTempFile(t, "key-*.key", keyContent)
	defer os.Remove(keyPath)

	t.Run("success", func(t *testing.T) {
		cert, key, err := getSelfSignedCertAndKey(certPath, keyPath)
		if err != nil {
			t.Errorf("could not getLocalhostCertAndKey: %v", err)
			return
		}
		if string(cert) != string(certContent) {
			t.Errorf("cert content mismatch, got %q, want %q", cert, certContent)
			return
		}
		if string(key) != string(keyContent) {
			t.Errorf("key content mismatch, got %q, want %q", key, keyContent)
			return
		}
	})

	t.Run("fail cert read", func(t *testing.T) {
		_, _, err := getSelfSignedCertAndKey("nonexistent-cert-file", keyPath)
		if err == nil {
			t.Errorf("expected error, got nil")
			return
		}
		if !strings.Contains(err.Error(), "(cert file)") {
			t.Errorf("expected error about cert file, got %v", err)
			return
		}
	})

	t.Run("fail key read", func(t *testing.T) {
		_, _, err := getSelfSignedCertAndKey(certPath, "nonexistent-key-file")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "(key file)") {
			t.Errorf("expected error about key file, got %v", err)
		}
	})
}

func TestUseSelfSigned(t *testing.T) {
	t.Chdir("..")
	env.Init()
	defer env.Reset()

	_, err := UseSelfSigned(env.Env.CertPath, env.Env.KeyPath)
	if err != nil {
		t.Errorf("could not UseLocalhost: %s", err)
		return
	}

	t.Run("invalid cert file paths", func(t *testing.T) {
		defer env.Reset()
		env.Env.CertPath = "nonexistent-cert-file"
		env.Env.KeyPath = "nonexistent-key-file"

		UseSelfSigned(env.Env.CertPath, env.Env.KeyPath)

	})

	t.Run("invalid cert files", func(t *testing.T) {
		defer env.Reset()

		tempDir := t.TempDir()

		// Create temp cert and key files with garbage content
		tempCert, err := os.CreateTemp(tempDir, "cert-*.pem")
		if err != nil {
			t.Fatalf("failed to create temp cert: %v", err)
		}
		defer tempCert.Close()

		tempKey, err := os.CreateTemp(tempDir, "key-*.pem")
		if err != nil {
			t.Fatalf("failed to create temp key: %v", err)
		}
		defer tempKey.Close()

		// Write invalid contents to simulate broken cert/key
		//tempCert.WriteString("not a real cert")
		//tempKey.WriteString("not a real key")

		// Update the env to point to the files
		env.Env.CertPath = tempCert.Name()
		env.Env.KeyPath = tempKey.Name()

		_, err = UseSelfSigned(env.Env.CertPath, env.Env.KeyPath)
		if err == nil {
			t.Errorf("expected an error due to invalid cert/key content")
		}
	})
}

func Test_loadTLSKeyPair(t *testing.T) {
	t.Chdir("..")
	cert, key, err := getSelfSignedCertAndKey(env.Env.CertPath, env.Env.KeyPath)
	if err != nil {
		t.Errorf("could not getLocalhostCertAndKey: %s", err)
		return
	}

	_, err = loadTLSKeyPair(cert, key)
	if err != nil {
		t.Errorf("could not loadTLSKeyPair: %s", err)
		return
	}

	t.Run("fail", func(t *testing.T) {
		_, err = loadTLSKeyPair(nil, nil)
		if err == nil {
			t.Errorf("error is nil")
			return
		}
	})
}

func TestUseACME(t *testing.T) {
	env.Init()
	_, err := UseAcme()
	if err != nil {
		t.Errorf("could not UseAcme: %s", err)
		return
	}

	t.Run("EmptyEnv_ShouldFail", func(t *testing.T) {
		defer env.Reset()

		env.Env.Hostnames = nil

		// Backup and clear environment variable
		_, err := UseAcme()
		if err == nil {
			t.Errorf("error is nil")
		}
	})
}

//	t.Run("ValidEnv_ShouldInitializeManager", func(t *testing.T) {
//		// Backup and set valid environment
//		oldEnv := os.Getenv("WEDDING_SERVICE_HOSTNAMES")
//		defer os.Setenv("WEDDING_SERVICE_HOSTNAMES", oldEnv)
//		os.Setenv("WEDDING_SERVICE_HOSTNAMES", "example.com:www.example.com,api.example.com|test.com:dev.test.com")
//
//		acmeManager, err := UseAcme()
//		if err != nil {
//			t.Fatalf("UseAcme failed: %v", err)
//		}
//		if acmeManager == nil {
//			t.Fatal("expected non-nil autocert.Manager")
//		}
//
//		t.Run("HostPolicy_Domain", func(t *testing.T) {
//			tests := []struct {
//				name       string
//				host       string
//				shouldFail bool
//			}{
//				{"example.com allowed", "example.com", false},
//				{"test.com allowed", "test.com", false},
//				{"unauthorized.com denied", "unauthorized.com", true},
//				{"another.com denied", "another.com", true},
//			}
//			for _, tt := range tests {
//				t.Run(tt.name, func(t *testing.T) {
//					err := acmeManager.HostPolicy(context.Background(), tt.host)
//					if tt.shouldFail && err == nil {
//						t.Errorf("expected error for host %q, got nil", tt.host)
//					}
//					if !tt.shouldFail && err != nil {
//						t.Errorf("unexpected error for host %q: %v", tt.host, err)
//					}
//				})
//			}
//		})
//
//		t.Run("HostPolicy_DomainAliases", func(t *testing.T) {
//			tests := []struct {
//				name       string
//				host       string
//				shouldFail bool
//			}{
//				{"www.example.com allowed", "www.example.com", false},
//				{"api.example.com allowed", "api.example.com", false},
//				{"dev.test.com allowed", "dev.test.com", false},
//				{"invalid.example.com denied", "invalid.example.com", true},
//			}
//			for _, tt := range tests {
//				t.Run(tt.name, func(t *testing.T) {
//					err := acmeManager.HostPolicy(context.Background(), tt.host)
//					if tt.shouldFail && err == nil {
//						t.Errorf("expected error for alias host %q, got nil", tt.host)
//					}
//					if !tt.shouldFail && err != nil {
//						t.Errorf("unexpected error for alias host %q: %v", tt.host, err)
//					}
//				})
//			}
//		})
//
//		t.Run("Client_RetryBackoff", func(t *testing.T) {
//			for n := 0; n <= 4; n++ {
//				expected := time.Duration(1<<n) * time.Second
//				actual := acmeManager.Client.RetryBackoff(n, nil, nil)
//				if actual != expected {
//					t.Errorf("retry backoff failed at attempt %d: got %v, want %v", n, actual, expected)
//				}
//			}
//		})
//	})
//}

func createTempFile(t *testing.T, pattern string, content []byte) string {
	t.Helper()
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	return f.Name()
}

func Test_handleRetryBackoff(t *testing.T) {
	for n := 0; n <= 10; n++ { // Test a few values to verify exponential backoff
		expected := time.Duration(1<<n) * time.Second
		actual := handleRetryBackoff(n, nil, nil)
		if actual != expected {
			t.Errorf("handleRetryBackoff(%d) = %v; want %v", n, actual, expected)
		}
	}
}
