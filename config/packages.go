package config

import (
	"path/filepath"

	"github.com/railduino/gd-tools/agent"
	"github.com/railduino/gd-tools/releases"
	"github.com/railduino/gd-tools/templates"
)

func (cfg *Config) DeployPackages(upgrade bool) error {
	cfg.Debug("Enter config/packages.go")

	catalog, err := releases.Load()
	if err != nil {
		return err
	}
	baseline, err := catalog.GetBaseline(cfg.Baseline)
	if err != nil {
		return err
	}

	if err := cfg.PackagesRepos(baseline); err != nil {
		return err
	}

	task := "installing"
	if cfg.CheckRemote("test -r /root/.gd-tools-first-run") {
		task = "checking"
	}
	cfg.Sayf("%s %d packages - please be patient ...", task, len(baseline.Packages))

	req := cfg.NewRequest()
	req.Packages = baseline.Packages
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

func (cfg *Config) PackagesRepos(bl *releases.Baseline) error {
	req := cfg.NewRequest()

	for _, name := range bl.Repos {
		keyName := name + ".gpg"
		keyTmpl := filepath.Join("apt", bl.Name, "keys", keyName)
		keyData, err := templates.Load(keyTmpl, cfg.Verbose)
		if err != nil {
			return err
		}
		req.AddFile(&agent.File{
			Task:    "write",
			Path:    agent.GetEtcDir("apt", "keyrings", keyName),
			Content: keyData,
			Mode:    "0644",
		})

		srcName := name + ".sources"
		oldName := name + ".list"
		srcTmpl := filepath.Join("apt", bl.Name, "sources", srcName)
		srcData, err := templates.Load(srcTmpl, cfg.Verbose)
		if err != nil {
			return err
		}
		req.AddFile(&agent.File{
			Task:    "write",
			Path:    agent.GetEtcDir("apt", "sources.list.d", srcName),
			Content: srcData,
			Mode:    "0644",
		})
		req.AddFile(&agent.File{
			Task: "delete",
			Path: agent.GetEtcDir("apt", "sources.list.d", oldName),
		})
	}

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}
