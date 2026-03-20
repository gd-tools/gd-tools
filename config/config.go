package config

import (
	"crypto/tls"
	"fmt"

	"github.com/gd-tools/gd-tools/utils"
	"github.com/urfave/cli/v2"
)

const (
	DefaultTimeout = 10
)

// Config contains the persistent server plus runtime-only helpers.
type Config struct {
	// Persistent production server configuration.
	Server

	// Concrete baseline for this server (Ubuntu version etc.).
	Baseline *Baseline `json:"-"`

	// Products - all downloadable assets.
	Products []Product `json:"-"`

	// Common runtime flags
	Verbose bool `json:"-"`
	Force   bool `json:"-"`
	Delete  bool `json:"-"`
	SkipDNS bool `json:"-"`
	SkipMX  bool `json:"-"`
	logger  *utils.Logger

	// Connection to the production server
	Port    string    `json:"-"`
	Conn    *tls.Conn `json:"-"`
	Timeout int       `json:"-"`
}

// ReadConfig loads and initializes a server configuration.
// Platform and Options can be injected for testing.
func ReadConfig(c *cli.Context) (*Config, error) {
	cfg := &Config{}

	// The persistent server configuration must exist, created by 'gdt setup'.
	err := utils.LoadJSON(utils.ConfigFile, &cfg.Server)
	if err != nil {
		return nil, err
	}

	// The baseline for this particular server.
	cfg.Baseline, err = LoadBaseline(cfg.BaselineName)
	if err != nil {
		return nil, err
	}

	if c != nil {
		cfg.Verbose = c.Bool("verbose")
		cfg.Force = c.Bool("force")
		cfg.Delete = c.Bool("delete")
		cfg.SkipDNS = c.Bool("skip-dns")
		cfg.SkipMX = c.Bool("skip-mx")
		cfg.Port = c.String("port")
		// TODO agent.SetAgentPort(cfg.Port)
	}

	return cfg, nil
}

// Save writes the persistent server configuration.
func (cfg *Config) Save() error {
	return utils.SaveJSON(utils.ConfigFile, &cfg.Server)
}

// Logging outout, controlled by verbosity.
func (cfg *Config) Info(args ...any) {
	if cfg != nil && cfg.Logger != nil {
		cfg.Logger.Info(args...)
	}
}

func (cfg *Config) Infof(fmt string, args ...any) {
	if cfg != nil && cfg.Logger != nil {
		cfg.Logger.Infof(fmt, args...)
	}
}

func (cfg *Config) Debug(args ...any) {
	if cfg != nil && cfg.Logger != nil && cfg.Verbose {
		cfg.Logger.Debug(args...)
	}
}

func (cfg *Config) Debugf(fmt string, args ...any) {
	if cfg != nil && cfg.Logger != nil && cfg.Verbose {
		cfg.Logger.Debugf(fmt, args...)
	}
}

func (cfg *Config) Error(args ...any) {
	if cfg != nil && cfg.Logger != nil {
		cfg.Logger.Error(args...)
	}
}

func (cfg *Config) Errorf(fmt string, args ...any) {
	if cfg != nil && cfg.Logger != nil {
		cfg.Logger.Errorf(fmt, args...)
	}
}

func (cfg *Config) RsyncFlags() string {
	if cfg.Verbose {
		return "-avz"
	}
	return "-avzq"
}
