package config

import (
	"crypto/tls"
	"fmt"

	"github.com/gd-tools/gd-tools/platform"
	"github.com/gd-tools/gd-tools/server"
	"github.com/gd-tools/gd-tools/utils"
	"github.com/urfave/cli/v2"
)

const (
	DefaultTimeout = 10
)

// Config contains the persistent server plus runtime-only helpers.
type Config struct {
	server.Server

	// Runtime environment (contains required Options).
	Platform *platform.Platform `json:"-"`

	// Concrete baseline for this server (Ubuntu version etc.).
	Baseline *platform.Baseline `json:"-"`

	// Common runtime flags
	Verbose bool `json:"-"`
	Force   bool `json:"-"`
	Delete  bool `json:"-"`
	SkipDNS bool `json:"-"`
	SkipMX  bool `json:"-"`

	// Connection to the production server
	Port    string    `json:"-"`
	Conn    *tls.Conn `json:"-"`
	Timeout int       `json:"-"`
}

// ReadConfig loads and initializes a server configuration.
// Platform and Options can be injected for testing.
func ReadConfig(c *cli.Context, pf *platform.Platform, opts *platform.Options) (*Config, error) {
	cfg := &Config{}

	// The persistent server configuration must exist, created by 'gdt setup'.
	err := utils.LoadJSON(utils.ConfigFile, &cfg.Server)
	if err != nil {
		return nil, err
	}

	// Load the server platform (necessary for recovery: never use "latest").
	// LoadPlatform ensures valid options, no need to check again.
	if pf == nil {
		pf, err = platform.LoadPlatform(cfg.BaselineName, opts)
		if err != nil {
			return nil, err
		}
	}
	cfg.Platform = pf

	// The baseline for this particular server.
	if err := pf.LoadBaselines(cfg.BaselineName); err != nil {
		return nil, err
	}
	cfg.Baseline = pf.Baseline

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
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	return utils.SaveJSON(utils.ConfigFile, &cfg.Server)
}

func (cfg *Config) RsyncFlags() string {
	if cfg.Verbose {
		return "-avz"
	}
	return "-avzq"
}
