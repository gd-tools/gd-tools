package config

import (
	"fmt"
	"net"

	"github.com/gd-tools/gd-tools/platform"
	"github.com/gd-tools/gd-tools/utils"
)

// These paths are fundamental, can be faked for tests.
func (cfg *Config) RootPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		panic("RootPath: config or platform is nil")
	}
	return cfg.Platform.RootPath(paths...)
}

func (cfg *Config) VarPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		panic("VarPath: config or platform is nil")
	}
	return cfg.Platform.VarPath(paths...)
}

func (cfg *Config) EtcPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		panic("EtcPath: config or platform is nil")
	}
	return cfg.Platform.EtcPath(paths...)
}

func (cfg *Config) BinPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		panic("BinPath: config or platform is nil")
	}
	return cfg.Platform.BinPath(paths...)
}

func (cfg *Config) RunPath(paths ...string) string {
	if cfg == nil || cfg.Platform == nil {
		panic("RunPath: config or platform is nil")
	}
	return cfg.Platform.RunPath(paths...)
}

// These paths are derived from the fundamental paths.
func (cfg *Config) DownloadsPath(paths ...string) string {
	return utils.JoinPath(cfg.Platform.RootPath("Downloads"), paths...)
}

func (cfg *Config) ToolsPath(paths ...string) string {
	return utils.JoinPath(cfg.Platform.VarPath("gd-tools"), paths...)
}

func (cfg *Config) DataPath(paths ...string) string {
	return utils.JoinPath(cfg.Platform.ToolsPath("data"), paths...)
}

func (cfg *Config) CertsPath(paths ...string) string {
	return utils.JoinPath(cfg.Platform.DataPath("certs"), paths...)
}

func (cfg *Config) ToolsApachePath(paths ...string) string {
	return utils.JoinPath(cfg.Platform.DataPath("apache"), paths...)
}

func (cfg *Config) LogsPath(paths ...string) string {
	return utils.JoinPath(cfg.Platform.ToolsPath("logs"), paths...)
}

func (cfg *Config) EtcApachePath(paths ...string) string {
	return utils.JoinPath(cfg.Platform.EtcPath("apache2"), paths...)
}

func (cfg *Config) EtcPhpPath(paths ...string) string {
	if cfg.Baseline == nil || cfg.Baseline.PHP == "" {
		panic("missing Baseline or PHP version")
	}
	return utils.JoinPath(cfg.Platform.EtcPath("php", cfg.Baseline.PHP), paths...)
}

// Render loads a template from the gdt binary.
func (cfg *Config) Render(name string, data interface{}) ([]byte, error) {
	if cfg == nil {
		return nil, fmt.Errorf("Render: config is nil")
	}

	return platform.Render(name, data)
}

func (cfg *Config) LookupIP(host string) ([]net.IP, error) {
	if cfg == nil || cfg.Platform == nil {
		return nil, fmt.Errorf("LookupIP: config or platform is nil")
	}
	return cfg.Platform.LookupIP(host)
}

func (cfg *Config) RSAKeyPair(fqdn string) ([]byte, []byte, error) {
	if cfg == nil || cfg.Platform == nil {
		return nil, nil, fmt.Errorf("RSAKeyPair: config or platform is nil")
	}
	return cfg.Platform.RSAKeyPair(fqdn)
}

// RunShell executes a list of shell commands on the local system.
func (cfg *Config) RunShell(commands []string) ([]byte, error) {
	if cfg == nil || cfg.Platform == nil {
		return nil, fmt.Errorf("RunShell: config or platform is nil")
	}
	return cfg.Platform.RunShell(commands)
}

// RunCommand executes a single shell command.
func (cfg *Config) RunCommand(name string, args ...string) ([]byte, error) {
	if cfg == nil || cfg.Platform == nil {
		return nil, fmt.Errorf("RunCommand: config or platform is nil")
	}
	return cfg.Platform.RunCommand(name, args...)
}
