package config

import (
	"crypto/tls"
	"fmt"

	"github.com/gd-tools/gd-tools/model"
	"github.com/gd-tools/gd-tools/platform"
	"github.com/gd-tools/gd-tools/utils"
)

const (
	DefaultTimeout = 10

	RunPrefix   = "[run]"
	DebugPrefix = "[dbg] ##########"
)

// Config contains the persistent server model plus runtime-only helpers.
type Config struct {
	model.Server

	// Runtime environment
	Platform *platform.Platform `json:"-"`

	// Concrete baseline for this server
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

// ReadConfig loads config.json into a new Config.
func ReadConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.Load(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Load reads the persistent server configuration.
func (cfg *Config) Load() error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if err := utils.LoadJSON(model.ConfigFile, &cfg.Server); err != nil {
		return err
	}

	cfg.initRuntime()
	return nil
}

// Save writes the persistent server configuration.
func (cfg *Config) Save() error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	return utils.SaveJSON(model.ConfigFile, &cfg.Server)
}

func (cfg *Config) initRuntime() {
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}
}
