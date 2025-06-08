package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Env *environment

type environment struct {
	DebugLevel int
	Mode       string
	HttpPort   string
	HttpsPort  string
}

func Init() {
	Env = &environment{
		DebugLevel: convertEnvToInt(os.Getenv("DEBUG")),
		Mode:       strings.ToLower(os.Getenv("MODE")),
		HttpPort:   os.Getenv("WEDDING_SERVICE_HTTP_PORT"),
		HttpsPort:  os.Getenv("WEDDING_SERVICE_HTTP_PORT"),
	}
}

func Reset() {
	Init()
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
