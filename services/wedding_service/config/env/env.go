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
	HttpPort() string
	CertPath() string
	KeyPath() string
	MySQLHost() string
	MySQLPort() string
	MySQLUser() string
	MySQLPassword() string
	MySQLDatabase() string
	MySQLRootPassword() string
}

type env struct {
	path       string
	debugLevel int
	httpPort   string
	certPath   string
	keyPath    string
	mysqlHost  string
	mysqlPort  string
	mysqlUser  string
	mysqlPass  string
	mysqlDB    string
	mysqlRoot  string
}

// NewEnv reads configuration exclusively from the process environment.
// Required variables: DEBUG, WEDDING_SERVICE_HTTP_PORT.
// Optional variables:
//   - WEDDING_SERVICE_HTTPS_PORT (defaults to "8443" if empty)
//   - WEDDING_SERVICE_HOSTNAMES (format: "host:alias1,alias2|host2:alias...")
//   - SELF_SIGNED_CERT_PATH (defaults to "/data/certs/localhost_wedding_service.crt")
//   - SELF_SIGNED_KEY_PATH (defaults to "/data/certs/localhost_wedding_service.key")
func NewEnv(path string) (Env, error) {

	debugStr := strings.TrimSpace(os.Getenv("DEBUG"))
	if debugStr == "" {

	}

	httpPort := strings.TrimSpace(os.Getenv("WEDDING_SERVICE_HTTP_PORT"))
	if httpPort == "" {

	}

	// Parse DEBUG level (must be integer)
	lvl, err := strconv.Atoi(debugStr)
	if err != nil {
		return nil, fmt.Errorf("DEBUG must be an integer, got %q: %w", debugStr, err)
	}

	e := &env{
		path:       path, // kept for Reset() compatibility; not used for file loading
		debugLevel: lvl,
		httpPort:   httpPort,
		certPath:   os.Getenv("SELF_SIGNED_CERT_PATH"),
		keyPath:    os.Getenv("SELF_SIGNED_KEY_PATH"),
		mysqlHost:  os.Getenv("MYSQL_HOST"),
		mysqlPort:  os.Getenv("MYSQL_PORT"),
		mysqlUser:  os.Getenv("MYSQL_USER"),
		mysqlPass:  os.Getenv("MYSQL_PASSWORD"),
		mysqlDB:    os.Getenv("MYSQL_DATABASE"),
		mysqlRoot:  os.Getenv("MYSQL_ROOT_PASSWORD"),
	}
	return e, nil
}

func (e *env) Reset() (Env, error) {
	return NewEnv(e.path)
}

func (e *env) DebugLevel() int           { return e.debugLevel }
func (e *env) HttpPort() string          { return e.httpPort }
func (e *env) CertPath() string          { return e.certPath }
func (e *env) KeyPath() string           { return e.keyPath }
func (e *env) MySQLHost() string         { return e.mysqlHost }
func (e *env) MySQLPort() string         { return e.mysqlPort }
func (e *env) MySQLUser() string         { return e.mysqlUser }
func (e *env) MySQLPassword() string     { return e.mysqlPass }
func (e *env) MySQLDatabase() string     { return e.mysqlDB }
func (e *env) MySQLRootPassword() string { return e.mysqlRoot }
