package config

import (
	"gateway_service/config/env"
	"os"
	"strings"
)

type Config interface {
	Hostnames() map[string][]string
	Backends() map[string]string
	HTTPPort() string
	HTTPSPort() string
	env.Env
}

type config struct {
	env       env.Env
	hostnames map[string][]string
	backends  map[string]string
}

func NewConfig() (Config, error) {
	e, err := env.NewEnv()
	if err != nil {
		return nil, err
	}

	backends := parseBackends(os.Getenv("GATEWAY_BACKENDS"))
	if len(backends) == 0 {
		// default: proxy localhost to tinyserver in docker network
		backends = map[string]string{"localhost": "http://tinyserver"}
	}

	hosts := os.Getenv("GATEWAY_HOSTNAMES")
	if strings.TrimSpace(hosts) == "" {
		// default allow both localhost and wedding-go-service hostnames
		hosts = "localhost|wedding-go-service"
	}

	return &config{
		env:       e,
		hostnames: parseHostnames(hosts), // certificate host policy
		backends:  backends,
	}, nil
}

func (c *config) DebugLevel() int {
	return c.env.DebugLevel()
}

func (c *config) Reset() (env.Env, error) {
	return c.env.Reset()
}

func (c *config) ServiceName() string {
	return c.env.ServiceName()
}

func (c *config) CertPath() string {
	return c.env.CertPath()
}

func (c *config) KeyPath() string {
	return c.env.KeyPath()
}

func (c *config) HTTPPort() string {
	return "80"
}

func (c *config) HTTPSPort() string {
	return "443"
}

func (c *config) Hostnames() map[string][]string {
	return c.hostnames
}

func (c *config) Backends() map[string]string {
	return c.backends
}
