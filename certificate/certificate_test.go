// Note: in order to test this, we need to have certificates generated first (see top of certificate.go),
// then run the test with the tag certificate

package certificate

import (
	"context"
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
	t.Chdir("..")
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

func TestHandleHostPolicy(t *testing.T) {
	domainAliases := map[string][]string{
		"example.com":   {"www.example.com", "api.example.com"},
		"another.com":   {"www.another.com"},
		"onlydomain.dk": {},
	}

	policy := handleHostPolicy(domainAliases)

	testCases := []struct {
		name        string
		host        string
		expectError bool
	}{
		{"exact match domain", "example.com", false},
		{"match alias", "www.example.com", false},
		{"another domain match", "another.com", false},
		{"another alias match", "www.another.com", false},
		{"no match at all", "unknown.com", true},
		{"match domain with no aliases", "onlydomain.dk", false},
		{"empty host", "", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := policy(context.Background(), testCase.host)
			if testCase.expectError && err == nil {
				t.Errorf("expected error for host %q but got none", testCase.host)
			} else if !testCase.expectError && err != nil {
				t.Errorf("unexpected error for host %q: %v", testCase.host, err)
			}
		})
	}
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
