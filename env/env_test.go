package env

import (
	"os"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	t.Chdir("..")
	defer Reset()

	_ = os.Setenv("DEBUG", "3")
	_ = os.Setenv("MODE", "Development")
	_ = os.Setenv("WEDDING_SERVICE_HTTP_PORT", "8080")
	_ = os.Setenv("WEDDING_SERVICE_HTTPS_PORT", "8443")
	_ = os.Setenv("WEDDING_SERVICE_HOSTNAMES", "example.com:www.example.com,www2.example.com|example2.com:www.example2.com")
	_ = os.Setenv("SELF_SIGNED_CERT_PATH", "./certificate/localhost_wedding_service.crt")
	_ = os.Setenv("SELF_SIGNED_KEY_PATH", "./certificate/localhost_wedding_service.key")

	err := Init()
	if err != nil {
		t.Errorf("%s", err)
	}

	if Env.DebugLevel != Info {
		t.Errorf("Expected DebugLevel to be %d, got %d", Info, Env.DebugLevel)
	}
	if Env.Mode != "development" {
		t.Errorf("Expected Mode to be 'development', got '%s'", Env.Mode)
	}
	if Env.HttpPort != "8080" {
		t.Errorf("Expected HttpPort to be '8080', got '%s'", Env.HttpPort)
	}
	if Env.HttpsPort != "8443" {
		t.Errorf("Expected HttpsPort to be '8443', got '%s'", Env.HttpsPort)
	}
	expected := readHostnames("example.com:www.example.com,www2.example.com|example2.com:www.example2.com")
	if !reflect.DeepEqual(Env.Hostnames, expected) {
		t.Errorf("Expected hostnames to be '%s', got '%s'", expected, Env.Hostnames)
	}
	if Env.CertPath != "./certificate/localhost_wedding_service.crt" {
		t.Errorf("Expected CertificatePath to be './certificate/localhost_wedding_service.crt', got '%s'", Env.CertPath)
	}
	if Env.KeyPath != "./certificate/localhost_wedding_service.key" {
		t.Errorf("Expected KeyPath to be './certificate/localhost_wedding_service.key', got '%s'", Env.KeyPath)
	}

	t.Run("test with invalid .env path", func(t *testing.T) {
		originalEnvPath := EnvPath
		defer func() {
			EnvPath = originalEnvPath
		}()
		EnvPath = "invalid file path"
		err = Init()
		if err == nil {
			t.Errorf("err should not be nil: %s", err)
		}
	})
}

func TestReset(t *testing.T) {
	t.Chdir("..")
	Reset()
}

func TestIsDebugInfoEnabled(t *testing.T) {
	t.Chdir("..")
	Init()
	defer Reset()

	Env = &environment{DebugLevel: Info}
	if !IsDebugInfoEnabled() {
		t.Error("Expected IsDebugInfoEnabled to return true")
	}

	Env = &environment{DebugLevel: All}
	if !IsDebugInfoEnabled() {
		t.Error("Expected IsDebugInfoEnabled to return true for All")
	}

	Env = &environment{DebugLevel: Warning}
	if IsDebugInfoEnabled() {
		t.Error("Expected IsDebugInfoEnabled to return false")
	}
}

func TestIsDebugWarningsEnabled(t *testing.T) {
	t.Chdir("..")
	Init()
	defer Reset()

	Env = &environment{DebugLevel: Warning}
	if !IsDebugWarningsEnabled() {
		t.Error("Expected IsDebugWarningsEnabled to return true")
	}

	Env = &environment{DebugLevel: All}
	if !IsDebugWarningsEnabled() {
		t.Error("Expected IsDebugWarningsEnabled to return true for All")
	}

	Env = &environment{DebugLevel: Error}
	if IsDebugWarningsEnabled() {
		t.Error("Expected IsDebugWarningsEnabled to return false")
	}
}

func TestIsDebugErrorsEnabled(t *testing.T) {
	t.Chdir("..")
	Init()
	defer Reset()

	Env = &environment{DebugLevel: Error}
	if !IsDebugErrorsEnabled() {
		t.Error("Expected IsDebugErrorsEnabled to return true")
	}

	Env = &environment{DebugLevel: All}
	if !IsDebugErrorsEnabled() {
		t.Error("Expected IsDebugErrorsEnabled to return true for All")
	}

	Env = &environment{DebugLevel: None}
	if IsDebugErrorsEnabled() {
		t.Error("Expected IsDebugErrorsEnabled to return false")
	}
}

func TestIsDebugDisabled(t *testing.T) {
	t.Chdir("..")
	Init()
	defer Reset()

	Env = &environment{DebugLevel: None}
	if !IsDebugDisabled() {
		t.Error("Expected IsDebugDisabled to return true")
	}

	Env = &environment{DebugLevel: All}
	if IsDebugDisabled() {
		t.Error("Expected IsDebugDisabled to return false")
	}
}

func TestIsModeDevelopment(t *testing.T) {
	t.Chdir("..")
	Init()
	defer Reset()

	Env = &environment{Mode: "development"}
	if !IsModeDevelopment() {
		t.Error("Expected IsModeDevelopment to return true")
	}

	Env = &environment{Mode: "production"}
	if IsModeDevelopment() {
		t.Error("Expected IsModeDevelopment to return false")
	}
}

func TestIsModeProduction(t *testing.T) {
	t.Chdir("..")
	Init()
	defer Reset()

	Env = &environment{Mode: "production"}
	if !IsModeProduction() {
		t.Error("Expected IsModeProduction to return true")
	}

	Env = &environment{Mode: "development"}
	if IsModeProduction() {
		t.Error("Expected IsModeProduction to return false")
	}
}

func TestIsModeNotSet(t *testing.T) {
	t.Chdir("..")
	Init()
	defer Reset()

	Env = &environment{Mode: ""}
	if !IsModeNotSet() {
		t.Error("Expected IsModeNotSet to return true")
	}

	Env = &environment{Mode: "something"}
	if IsModeNotSet() {
		t.Error("Expected IsModeNotSet to return false")
	}
}

func TestConvertEnvToInt(t *testing.T) {
	if got := convertEnvToInt("2"); got != 2 {
		t.Errorf("Expected 2, got %d", got)
	}

	if got := convertEnvToInt("invalid"); got != -1 {
		t.Errorf("Expected -1 for invalid input, got %d", got)
	}
}

func TestReadHostnames(t *testing.T) {
	expected := map[string][]string{
		"example.com":  {"www.example.com", "www2.example.com"},
		"example2.com": {"www.example2.com"},
	}

	result := readHostnames("example.com:www.example.com,www2.example.com|example2.com:www.example2.com")
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected hostnames to be '%s', got '%s'", expected, result)
	}
}
