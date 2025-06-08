package env

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	oldEnv := Env
	defer func() { Env = oldEnv }()

	_ = os.Setenv("DEBUG", "3")
	_ = os.Setenv("MODE", "Development")
	_ = os.Setenv("WEDDING_SERVICE_HTTP_PORT", "8080")

	Init()

	if Env.DebugLevel != Info {
		t.Errorf("Expected DebugLevel to be %d, got %d", Info, Env.DebugLevel)
	}
	if Env.Mode != "development" {
		t.Errorf("Expected Mode to be 'development', got '%s'", Env.Mode)
	}
	if Env.HttpPort != "8080" {
		t.Errorf("Expected HttpPort to be '8080', got '%s'", Env.HttpPort)
	}
	if Env.HttpsPort != "8080" {
		t.Errorf("Expected HttpsPort to be '8080', got '%s'", Env.HttpsPort)
	}
}

func TestReset(t *testing.T) {
	Reset()
}

func TestIsDebugInfoEnabled(t *testing.T) {
	oldEnv := Env
	defer func() { Env = oldEnv }()

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
	oldEnv := Env
	defer func() { Env = oldEnv }()

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
	oldEnv := Env
	defer func() { Env = oldEnv }()

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
	oldEnv := Env
	defer func() { Env = oldEnv }()

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
	oldEnv := Env
	defer func() { Env = oldEnv }()

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
	oldEnv := Env
	defer func() { Env = oldEnv }()

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
	oldEnv := Env
	defer func() { Env = oldEnv }()

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
