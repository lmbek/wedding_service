package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Env interface {
	Reset() (Env, error)
	DebugLevel() int
	ServiceName() string
	CertPath() string
	KeyPath() string
}

type env struct {
	debugLevel  int
	serviceName string
	certPath    string
	keyPath     string
}

func NewEnv() (Env, error) {
	var (
		err         error
		serviceName string
		certPath    string
		keyPath     string
	)

	serviceName, err = required("SERVICE_NAME")
	if err != nil {
		return nil, err
	}
	certPath, err = required("CERT_PATH")
	if err != nil {
		return nil, err
	}
	keyPath, err = required("KEY_PATH")
	if err != nil {
		return nil, err
	}
	debugLevel := optional(os.Getenv("DEBUG_LEVEL"))

	// Optional DEBUG (default 4)
	level := 4
	if debugLevel != "" {
		level, err = strconv.Atoi(debugLevel)
		if err != nil {
			return nil, fmt.Errorf("%s must be an integer: %w", "DEBUG_LEVEL", err)
		}
	}

	return &env{
		debugLevel:  level,
		serviceName: serviceName,
		certPath:    certPath,
		keyPath:     keyPath,
	}, nil
}

func (e *env) Reset() (Env, error) { return NewEnv() }
func (e *env) DebugLevel() int     { return e.debugLevel }
func (e *env) ServiceName() string { return e.serviceName }
func (e *env) CertPath() string    { return e.certPath }
func (e *env) KeyPath() string     { return e.keyPath }

func required(name string) (string, error) {
	val := strings.TrimSpace(os.Getenv(name))
	if val == "" {
		return "", fmt.Errorf("missing required environment variable: %s", name)
	}
	return val, nil
}

func optional(env string) string {
	return env
}
