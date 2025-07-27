package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

type Env interface {
	Reset() (Env, error)

	DebugLevel() int
	IsDebugInfoEnabled() bool
	IsDebugWarningsEnabled() bool
	IsDebugErrorsEnabled() bool
	IsDebugDisabled() bool

	HttpPort() string
	HttpsPort() string
	Hostnames() map[string][]string
	CertPath() string
	KeyPath() string
}

type env struct {
	path       string
	debugLevel int
	httpPort   string
	httpsPort  string
	hostnames  map[string][]string
	certPath   string
	keyPath    string
}

func NewEnv(path string) (Env, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	e := &env{
		path:       path,
		debugLevel: convertEnvToInt(os.Getenv("DEBUG")),
		httpPort:   os.Getenv("WEDDING_SERVICE_HTTP_PORT"),
		httpsPort:  os.Getenv("WEDDING_SERVICE_HTTPS_PORT"),
		hostnames:  readHostnames(os.Getenv("WEDDING_SERVICE_HOSTNAMES")),
		certPath:   os.Getenv("SELF_SIGNED_CERT_PATH"),
		keyPath:    os.Getenv("SELF_SIGNED_KEY_PATH"),
	}
	return e, nil
}

func (e *env) Reset() (Env, error) {
	return NewEnv(e.path)
}

// Getters

func (e *env) DebugLevel() int {
	return e.debugLevel
}

func (e *env) HttpPort() string {
	return e.httpPort
}

func (e *env) HttpsPort() string {
	return e.httpsPort
}

func (e *env) Hostnames() map[string][]string {
	return e.hostnames
}

func (e *env) CertPath() string {
	return e.certPath
}

func (e *env) KeyPath() string {
	return e.keyPath
}

// Debug level constants
const (
	None = iota
	Error
	Warning
	Info
	All
)

// Debug-level check methods

func (e *env) IsDebugInfoEnabled() bool     { return e.debugLevel == Info || e.debugLevel == All }
func (e *env) IsDebugWarningsEnabled() bool { return e.debugLevel == Warning || e.debugLevel == All }
func (e *env) IsDebugErrorsEnabled() bool   { return e.debugLevel == Error || e.debugLevel == All }
func (e *env) IsDebugDisabled() bool        { return e.debugLevel == None }

func convertEnvToInt(envVar string) int {
	i, err := strconv.Atoi(envVar)
	if err != nil {
		fmt.Printf("DEBUG should be set and be an integer: %v\n", err)
		return -1
	}
	return i
}

// Parses hostnames from a string like "host1:alias1,alias2|host2:alias3,alias4"
func readHostnames(hostnames string) map[string][]string {
	domainAliases := make(map[string][]string)
	if hostnames == "" {
		return domainAliases
	}

	groups := strings.Split(hostnames, "|")
	for _, group := range groups {
		if group == "" {
			continue
		}
		parts := strings.SplitN(group, ":", 2)
		hostname := strings.TrimSpace(parts[0])
		var aliases []string
		if len(parts) == 2 {
			for _, a := range strings.Split(parts[1], ",") {
				a = strings.TrimSpace(a)
				if a != "" {
					aliases = append(aliases, a)
				}
			}
		}
		if hostname != "" {
			domainAliases[hostname] = aliases
		}
	}

	return domainAliases
}
