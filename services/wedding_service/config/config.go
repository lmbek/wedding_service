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

func (c *config) HttpPort() string {
	return c.env.HttpPort()
}

func (c *config) HttpsPort() string {
	return c.env.HttpsPort()
}

func (c *config) MySQLHost() string         { return c.env.MySQLHost() }
func (c *config) MySQLPort() string         { return c.env.MySQLPort() }
func (c *config) MySQLUser() string         { return c.env.MySQLUser() }
func (c *config) MySQLPassword() string     { return c.env.MySQLPassword() }
func (c *config) MySQLDatabase() string     { return c.env.MySQLDatabase() }
func (c *config) MySQLRootPassword() string { return c.env.MySQLRootPassword() }
