package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
	"sync"
)

var EnvPath = ".env"
var Env *environment

type environment struct {
	DebugLevel int
	Mode       string
	HttpPort   string
	HttpsPort  string
	Hostnames  map[string][]string
	CertPath   string
	KeyPath    string
}

var mutex sync.RWMutex

func Init() error {

	// TODO: make a lock for all getters and setters, and then also test if godotenv in general has race conditions
	mutex.Lock()
	defer mutex.Unlock()
	err := godotenv.Load(EnvPath) // Note: doc says it will not overwrite env vars that already exists, this might give issues
	if err != nil {
		return fmt.Errorf("err loading .env file: %s", err)
	}

	Env = &environment{
		DebugLevel: convertEnvToInt(os.Getenv("DEBUG")),
		Mode:       strings.ToLower(os.Getenv("MODE")),
		HttpPort:   os.Getenv("WEDDING_SERVICE_HTTP_PORT"),
		HttpsPort:  os.Getenv("WEDDING_SERVICE_HTTPS_PORT"),
		Hostnames:  readHostnames(os.Getenv("WEDDING_SERVICE_HOSTNAMES")),
		CertPath:   os.Getenv("SELF_SIGNED_CERT_PATH"),
		KeyPath:    os.Getenv("SELF_SIGNED_KEY_PATH"),
	}

	return nil
}

func Reset() error {
	return Init()
}

const (
	None = iota
	Error
	Warning
	Info
	All
)

func IsDebugInfoEnabled() bool {
	return Env.DebugLevel == Info || Env.DebugLevel == All
}

func IsDebugWarningsEnabled() bool {
	return Env.DebugLevel == Warning || Env.DebugLevel == All
}

func IsDebugErrorsEnabled() bool {
	return Env.DebugLevel == Error || Env.DebugLevel == All
}

func IsDebugDisabled() bool {
	return Env.DebugLevel == None
}

func IsModeDevelopment() bool {
	return Env.Mode == "development"
}

func IsModeProduction() bool {
	return Env.Mode == "production"
}

func IsModeNotSet() bool {
	return Env.Mode == ""
}

func convertEnvToInt(envVar string) int {
	i, err := strconv.Atoi(envVar)
	if err != nil {
		fmt.Printf("DEBUG should be set and be an integer %s \n", err)
		return -1
	}
	return i
}

func readHostnames(hostnames string) (domainAliases map[string][]string) {
	domainAliases = make(map[string][]string)

	// Split hostname groups by semicolon
	hostnameGroups := strings.Split(hostnames, "|")
	for _, group := range hostnameGroups {
		// Split hostname and aliases by colon
		parts := strings.SplitN(group, ":", 2)
		hostname := parts[0]
		var aliases []string
		if len(parts) == 2 {
			aliases = strings.Split(parts[1], ",")
		}
		domainAliases[hostname] = aliases
	}

	return
}
