package config

import (
	"fmt"

	"github.com/gd-tools/gd-tools/platform"
	"github.com/gd-tools/gd-tools/utils"
)

// These paths are fundamental, can be faked for tests.
func (cfg *Config) RootPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		return nil, fmt.Errorf("RootPath: config or platform is nil")
	}
	return cfg.Platform.RootPath(paths...)
}

func (cfg *Config) VarPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		return nil, fmt.Errorf("VarPath: config or platform is nil")
	}
	return cfg.Platform.VarPath(paths...)
}

func (cfg *Config) EtcPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		return nil, fmt.Errorf("EtcPath: config or platform is nil")
	}
	return cfg.Platform.EtcPath(paths...)
}

func (cfg *Config) BinPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		return nil, fmt.Errorf("BinPath: config or platform is nil")
	}
	return cfg.Platform.BinPath(paths...)
}

func (cfg *Config) RunPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		return nil, fmt.Errorf("RunPath: config or platform is nil")
	}
	return cfg.Platform.RunPath(paths...)
}

// These paths are derived from the fundamental paths.
func (cfg *Config) DownloadsPath(paths ...string) string {
	return utils.JoinPath(cfg.RootPath("Downloads"), paths...)
}

func (cfg *Config) ToolsPath(paths ...string) string {
	return utils.JoinPath(cfg.VarPath("gd-tools"), paths...)
}

func (cfg *Config) DataPath(paths ...string) string {
	return utils.JoinPath(cfg.ToolsPath("data"), paths...)
}

func (cfg *Config) CertsPath(paths ...string) string {
	return utils.JoinPath(cfg.DataPath("certs"), paths...)
}

func (cfg *Config) ToolsApachePath(paths ...string) string {
	return utils.JoinPath(cfg.DataPath("apache"), paths...)
}

func (cfg *Config) LogsPath(paths ...string) string {
	return utils.JoinPath(cfg.ToolsPath("logs"), paths...)
}

func (cfg *Config) EtcApachePath(paths ...string) string {
	return utils.JoinPath(cfg.EtcPath("apache2"), paths...)
}

func (cfg *Config) EtcPhpPath(paths ...string) string {
	if cfg.Baseline == nil || cfg.Baseline.PHP == "" {
		panic("missing Baseline or PHP version")
	}
	return utils.JoinPath(cfg.EtcPath("php", cfg.Baseline.PHP), paths...)
}

// Render loads a template from the gdt binary.
func (cfg *Config) Render(name string, data interface{}) ([]byte, error) {
	if cfg == nil {
		return nil, fmt.Errorf("Render: config is nil")
	}

	return platform.Render(name, data)
}

// TODO LookupIP   func(host string) ([]net.IP, error)
// TODO RSAKeyPair func(fqdn string) ([]byte, []byte, error)
// TODO RunShell   func(commands []string) ([]byte, error)

// RunCommand executes a single shell command.
func (cfg *Config) RunCommand(name string, args ...string) ([]byte, error) {
	if cfg == nil || cfg.Platform == nil {
		return nil, fmt.Errorf("RunCommand: config or platform is nil")
	}

	return cfg.Platform.RunCommand(name, args...)
}
