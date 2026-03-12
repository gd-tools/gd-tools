package config

import (
	"fmt"
	"path/filepath"

	"github.com/gd-tools/gd-tools/assets"
	"github.com/gd-tools/gd-tools/agent"
)

func (cfg *Config) DeployPackages(upgrade bool) error {
	cfg.Debug("Enter config/packages.go")

	if err := cfg.PackagesRepos(); err != nil {
		return err
	}

	task := "installing"
	if cfg.CheckRemote("test -r /root/.gd-tools-first-run") {
		task = "checking"
	}
	cfg.Sayf("%s %d packages - please be patient ...", task, len(cfg.Baseline.Packages))

	req := cfg.NewRequest()
	req.Packages = cfg.Baseline.Packages
	req.Upgrade = upgrade
	req.Firewall = cfg.Firewall
	req.UbuntuPro = cfg.UbuntuPro

	if err := req.Send(); err != nil {
		return err
	}

	cfg.PushCerts()

	cfg.Debug("Leave config/packages.go")
	return nil
}

func (cfg *Config) PackagesRepos() error {
	req := cfg.NewRequest()

	for _, name := range cfg.Baseline.Repos {
		keyTmpl := fmt.Sprintf("apt/%s/keys/%s.gpg", cfg.BaselineName, name)
		keyData, err := assets.Render(keyTmpl, nil)
		if err != nil {
			return err
		}
		req.AddFile(&agent.File{
			Task:    "write",
			Path:    assets.GetEtcDir("apt", "keyrings", keyName),
			Content: keyData,
			Mode:    "0644",
		})

		srcName := name + ".sources"
		oldName := name + ".list"
		srcTmpl := filepath.Join("apt", cfg.BaselineName, "sources", srcName)
		srcData, err := assets.Render(srcTmpl, nil)
		if err != nil {
			return err
		}
		req.AddFile(&agent.File{
			Task:    "write",
			Path:    assets.GetEtcDir("apt", "sources.list.d", srcName),
			Content: srcData,
			Mode:    "0644",
		})
		req.AddFile(&agent.File{
			Task: "delete",
			Path: assets.GetEtcDir("apt", "sources.list.d", oldName),
		})
	}

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}
