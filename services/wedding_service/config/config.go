package config

import (
	"wedding_service/config/env"
	"wedding_service/config/flags"
)

type Config interface {
	load() (Config, error)
	flags.Flags
	env.Env
}

type config struct {
	flags flags.Flags
	env   env.Env
}

func NewConfig() (Config, error) {
	c := &config{}
	return c.load()
}

func (c *config) load() (Config, error) {
	var err error
	c.flags = flags.NewFlags()
	c.env, err = env.NewEnv(".env")
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Flags

func (c *config) FrontendPath() string {
	return c.flags.FrontendPath()
}

func (c *config) HotReloadEnabled() bool {
	return c.flags.HotReloadEnabled()
}

// Env

func (c *config) DebugLevel() int {
	return c.env.DebugLevel()
}

func (c *config) HttpPort() string {
	return c.env.HttpPort()
}

func (c *config) CertPath() string {
	return c.env.CertPath()
}

func (c *config) KeyPath() string {
	return c.env.KeyPath()
}

func (c *config) Reset() (env.Env, error) {
	return c.env.Reset()
}
