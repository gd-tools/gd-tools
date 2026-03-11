package config

import (
	"path/filepath"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/releases"
	"github.com/gd-tools/gd-tools/templates"
)

func (cfg *Config) PhpEtcDir(paths ...string) string {
	phpDir := releases.GetEtcDir("php", cfg.Baseline.PHP)
	if len(paths) == 0 {
		return phpDir
	}
	return filepath.Join(append([]string{phpDir}, paths...)...)
}

func (cfg *Config) PhpFpmPoolPath(name string) string {
	return releases.GetEtcDir("php", cfg.Baseline.PHP, "fpm", "pool.d", name+".conf")
}

func (cfg *Config) PhpFpmService() string {
	return "php" + cfg.Baseline.PHP + "-fpm"
}

func (cfg *Config) DeployPHP() error {
	cfg.Debug("Enter config/php.go")

	req := cfg.NewRequest()

	content, err := templates.Parse("php/custom.ini", cfg.Verbose, cfg)
	if err != nil {
		return err
	}

	phpDirs := []string{
		"cli",
		"fpm",
	}

	for _, name := range phpDirs {
		service := "apache2"
		if name == "fpm" {
			service = cfg.PhpFpmService()
		}
		file := agent.File{
			Task:    "write",
			Path:    cfg.PhpEtcDir(name, "conf.d", "60-custom.ini"),
			Content: content,
			Mode:    "0644",
			Service: service,
		}
		req.AddFile(&file)
	}

	if err := req.Send(); err != nil {
		return err
	}

	cfg.Debug("Leave config/php.go")
	return nil
}
