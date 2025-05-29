// Note: in order to test this, we need to have certificates generated first (see top of certificate.go),
// then run the test with the tag certificate

package certificate

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestUseACME(t *testing.T) {
	autocertManager := UseACME()
	if autocertManager == nil {
		t.Errorf("autocertManager should not be nil")
	}

	//if autocertManager.Prompt == nil {
	//	t.Errorf("autocertManager.Prompt should not be nil")
	//}
	//
	//if autocertManager.Cache == nil {
	//	t.Errorf("autocertManager.Cache should not be nil")
	//}
	//
	//if autocertManager.Client == nil {
	//	t.Errorf("autocertManager.Client should not be nil")
	//}
	//
	//if autocertManager.HostPolicy == nil {
	//	t.Errorf("autocertManager.HostPolicy should not be nil")
	//}

	t.Run("hostpolicy_domain_test", func(t *testing.T) {
		type HostPolicyTestCase struct {
			inputHostname               string
			shouldExpectValidationError bool
		}

		directHostnameTestCases := []HostPolicyTestCase{
			{inputHostname: os.Getenv("WEDDING_SERVICE_EXTERNAL_HOSTNAME"), shouldExpectValidationError: false},
			{inputHostname: "unauthorized.com", shouldExpectValidationError: true},
			{inputHostname: "example.com", shouldExpectValidationError: true},
		}

		for _, hostPolicyTestCase := range directHostnameTestCases {
			validationError := autocertManager.HostPolicy(context.Background(), hostPolicyTestCase.inputHostname)

			if hostPolicyTestCase.shouldExpectValidationError && validationError == nil {
				t.Errorf("Expected validation error for hostname '%s', but got none.", hostPolicyTestCase.inputHostname)
			}
			if !hostPolicyTestCase.shouldExpectValidationError && validationError != nil {
				t.Errorf("Did not expect error for hostname '%s', but got: %v", hostPolicyTestCase.inputHostname, validationError)
			}
		}
	})

	t.Run("hostpolicy_domain_alias_test", func(t *testing.T) {
		type HostPolicyTestCase struct {
			inputHostname               string
			shouldExpectValidationError bool
		}

		aliasHostnameTestCases := []HostPolicyTestCase{
			{inputHostname: os.Getenv("WEDDING_SERVICE_EXTERNAL_HOSTNAME_ALIAS1"), shouldExpectValidationError: false},
			{inputHostname: "unauthorized.com", shouldExpectValidationError: true},
			{inputHostname: "api.example.com", shouldExpectValidationError: true},
			{inputHostname: "dev.test.com", shouldExpectValidationError: true},
		}

		for _, hostPolicyTestCase := range aliasHostnameTestCases {
			validationError := autocertManager.HostPolicy(context.Background(), hostPolicyTestCase.inputHostname)

			if hostPolicyTestCase.shouldExpectValidationError && validationError == nil {
				t.Errorf("Expected validation error for alias hostname '%s', but got none.", hostPolicyTestCase.inputHostname)
			}
			if !hostPolicyTestCase.shouldExpectValidationError && validationError != nil {
				t.Errorf("Did not expect error for alias hostname '%s', but got: %v", hostPolicyTestCase.inputHostname, validationError)
			}
		}
	})

	t.Run("client_retry_backoff_test", func(t *testing.T) {
		type RetryBackoffTestCase struct {
			attemptNumber int
			expectedDelay time.Duration
		}

		retryBackoffTestCases := []RetryBackoffTestCase{
			{attemptNumber: 0, expectedDelay: 1 * time.Second},  // 2^0 = 1
			{attemptNumber: 1, expectedDelay: 2 * time.Second},  // 2^1 = 2
			{attemptNumber: 2, expectedDelay: 4 * time.Second},  // 2^2 = 4
			{attemptNumber: 3, expectedDelay: 8 * time.Second},  // 2^3 = 8
			{attemptNumber: 4, expectedDelay: 16 * time.Second}, // 2^4 = 16
		}

		for _, retryBackoffTestCase := range retryBackoffTestCases {
			got := autocertManager.Client.RetryBackoff(retryBackoffTestCase.attemptNumber, nil, nil)

			if got != retryBackoffTestCase.expectedDelay {
				t.Errorf("For attempt %d: expected %v, got %v", retryBackoffTestCase.attemptNumber, retryBackoffTestCase.expectedDelay, got)
			}
		}
	})

}

func TestUseLOCALHOST(t *testing.T) {

	cert, key, err := getLocalhostCertAndKey(os.Getenv("LOCALHOST_CERT"), os.Getenv("LOCALHOST_KEY"))
	if err != nil {
		t.Errorf("Could not get cert and/or key: %v \n", err)
	}

	tlsCert, err := UseLOCALHOST(cert, key)
	if err != nil {
		t.Errorf("Got err: %v", err)
	}

	if tlsCert == nil {
		t.Errorf("tlsCert should not be nil")
	}

	t.Run("test_invalid_paths", func(t *testing.T) {
		cert, key, err := getLocalhostCertAndKey("invalid", "invalid")

		_, err = UseLOCALHOST(cert, key)
		if err == nil {
			t.Errorf("expected error due to invalid certificate")
		}
	})
}

func Test_getLocalhostCertAndKeys(t *testing.T) {
	cert, key, err := getLocalhostCertAndKey(os.Getenv("LOCALHOST_CERT"), os.Getenv("LOCALHOST_KEY"))
	if err != nil {
		t.Errorf("could not get localhost cert and keys: %v \n", err)
	}

	if cert == nil {
		t.Errorf("embeddedCert should not be nil \n")
	}

	if key == nil {
		t.Errorf("embeddedKey should not be nil \n")
	}

	t.Run("test_invalid_paths", func(t *testing.T) {
		_, _, err = getLocalhostCertAndKey("invalid_path", os.Getenv("LOCALHOST_KEY"))
		if err == nil {
			t.Errorf("should not get error \n")
		}

		_, _, err = getLocalhostCertAndKey("LOCALHOST_CERT", "invalid_path")
		if err == nil {
			t.Errorf("should not get error \n")
		}
	})
}
