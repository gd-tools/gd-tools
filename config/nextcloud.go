package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/gd-tools/gd-tools/assets"
	"github.com/gd-tools/gd-tools/agent"
)

func (cfg *Config) DeployNextcloud(nc *agent.Nextcloud) error {
	if nc == nil {
		return fmt.Errorf("missing Nextcloud pointer")
	}
	cfg.Debugf("Enter config/nextcloud.go (%s)", nc.FQDN())

	_, ncRel, err := cfg.Catalog.Get("nextcloud", nc.Version)
	if err != nil {
		return err
	}
	if ncRel.Download.Directory == "" {
		return fmt.Errorf("missing Directory in Nextcloud download")
	}
	nc.Download = &ncRel.Download

	if err := cfg.NextcloudDownload(nc); err != nil {
		return err
	}
	if err := cfg.NextcloudExtract(nc); err != nil {
		return err
	}

	if err := cfg.NextcloudMkDirs(nc); err != nil {
		return err
	}

	if err := cfg.NextcloudSQL(nc); err != nil {
		return err
	}

	if err := cfg.NextcloudBackupHook(nc); err != nil {
		return err
	}

	if err := cfg.NextcloudSetupConfig(nc, "maintenance:install"); err != nil {
		return err
	}

	for _, entry := range nc.GetConfigList() {
		if err := cfg.NextcloudSetupConfig(nc, entry.Key); err != nil {
			return err
		}
	}

	if err := cfg.NextcloudCronJob(nc); err != nil {
		return err
	}

	if err := cfg.EnsureCertificate(nc.FQDN()); err != nil {
		return err
	}

	if err := cfg.NextcloudSetupPool(nc); err != nil {
		return err
	}

	if err := cfg.NextcloudVhost(nc); err != nil {
		return err
	}

	if status, err := cfg.SetCNAME(nc.DomainName, nc.HostName); err != nil {
		return err
	} else if status != "" {
		cfg.Say(status)
	}

	// if all went well, install the occ-<name> command
	occSrc := assets.GetBinDir("gd-occ")
	occDst := assets.GetBinDir("occ-" + nc.Name)
	if _, err := cfg.LocalCommand(
		"rsync",
		cfg.RsyncFlags(),
		"--chown=root:root",
		"--chmod=0500",
		occSrc,
		cfg.RootUser()+":"+occDst,
	); err != nil {
		return fmt.Errorf("failed to install %s: %w", occDst, err)
	}

	cfg.Debugf("Leave config/nextcloud.go (%s)", nc.FQDN())
	return nil
}

func (cfg *Config) NextcloudDownload(nc *agent.Nextcloud) error {
	req := cfg.NewRequest()

	req.Downloads = append(req.Downloads, nc.Download)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) NextcloudExtract(nc *agent.Nextcloud) error {
	req := cfg.NewRequest()

	extract := agent.File{
		Task:   "extract",
		Path:   assets.GetDownloadsDir(nc.Download.Filename),
		Target: nc.RootDir(),
		Mode:   "0755",
		User:   "root",
		Group:  "root",
	}
	req.AddFile(&extract)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) NextcloudMkDirs(nc *agent.Nextcloud) error {
	req := cfg.NewRequest()

	dataMkdir := agent.File{
		Task:  "mkdir",
		Path:  nc.DataDir(""),
		Mode:  "0700",
		User:  "www-data",
		Group: "www-data",
	}
	req.AddFile(&dataMkdir)

	logsMkdir := agent.File{
		Task:  "mkdir",
		Path:  nc.LogsDir(),
		Mode:  "0755",
		User:  "root",
		Group: "root",
	}
	req.AddFile(&logsMkdir)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) NextcloudSQL(nc *agent.Nextcloud) error {
	req := cfg.NewRequest()

	sqlTmpl := filepath.Join("nextcloud", "create.sql")
	sqlStmts, err := assets.SQL(sqlTmpl, nc)
	if err != nil {
		return err
	}

	sql := agent.MySQL{
		Stmts:   sqlStmts,
		Comment: fmt.Sprintf("create nextcloud (%s) tables", nc.Name),
	}
	req.MySQLs = append(req.MySQLs, &sql)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) NextcloudBackupHook(nc *agent.Nextcloud) error {
	req := cfg.NewRequest()

	hookTmpl, err := assets.Render("nextcloud/backup", nc)
	if err != nil {
		return err
	}

	hookFile := agent.File{
		Task:    "write",
		Path:    nc.HookPath(),
		Mode:    "0500",
		Content: hookTmpl,
		User:    "root",
		Group:   "root",
	}
	req.AddFile(&hookFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) NextcloudSetupConfig(nc *agent.Nextcloud, key string) error {
	req := cfg.NewRequest()

	req.Nextcloud = nc
	req.NextConf = key

	if key == "maintenance:install" {
		configTmpl, err := json.MarshalIndent(nc, "", "  ")
		if err != nil {
			return err
		}
		configFile := agent.File{
			Task:    "write",
			Path:    nc.ConfigPath(),
			Content: configTmpl,
			Mode:    "0600",
			User:    "root",
			Group:   "root",
		}
		req.AddFile(&configFile)
	}

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) NextcloudCronJob(nc *agent.Nextcloud) error {
	req := cfg.NewRequest()

	cronTmpl, err := assets.Render("nextcloud/cron.d", nc)
	if err != nil {
		return err
	}

	cronFile := agent.File{
		Task:    "write",
		Path:    nc.CronPath(),
		Content: cronTmpl,
		Mode:    "0644",
		User:    "root",
		Group:   "root",
	}
	req.AddFile(&cronFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) NextcloudSetupPool(nc *agent.Nextcloud) error {
	req := cfg.NewRequest()

	poolTmpl, err := assets.Render("nextcloud/php-fpm-pool.conf", nc)
	if err != nil {
		return err
	}

	poolPath := cfg.PhpFpmPoolPath(nc.Name)
	poolFile := agent.File{
		Task:    "write",
		Path:    poolPath,
		Content: poolTmpl,
		Mode:    "0644",
		User:    "root",
		Group:   "root",
		Service: cfg.PhpFpmService(),
	}
	req.AddFile(&poolFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) NextcloudVhost(nc *agent.Nextcloud) error {
	req := cfg.NewRequest()

	vhostTmpl, err := assets.Render("nextcloud/vhost.conf", nc)
	if err != nil {
		return err
	}

	vhostFile := agent.File{
		Task:    "write",
		Path:    nc.VhostPath(),
		Content: vhostTmpl,
	}
	req.AddFile(&vhostFile)
	req.AddService("apache2")

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}
