package config

import (
	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/php"
	"github.com/gd-tools/gd-tools/templates"
)

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
			service = php.GetPhpFpmService()
		}
		file := agent.File{
			Task:    "write",
			Path:    php.GetPhpEtcDir(name, "conf.d", "60-custom.ini"),
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
