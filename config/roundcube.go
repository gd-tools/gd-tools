package config

import (
	"fmt"
	"path/filepath"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/assets"
	"github.com/gd-tools/gd-tools/email"
	"github.com/gd-tools/gd-tools/utils"
)

type Roundcube struct {
	DomainName string
	Name       string
	PhpVersion string
	FQDN       string
	SysAdmin   string
	Locale     string
	Password   string
	DesKey     string
	Directory  string
	Download   *assets.Download
}

func (rc *Roundcube) WebMail() string {
	return "webmail." + rc.DomainName
}

func (rc *Roundcube) RootDir() string {
	return assets.GetToolsDir("data", "roundcube", rc.Name)
}

func (rc *Roundcube) BaseDir() string {
	return filepath.Join(rc.RootDir(), rc.Directory)
}

func (rc *Roundcube) CertDir() string {
	return assets.GetToolsDir("data", "certs", rc.WebMail())
}

func (rc *Roundcube) LogsDir() string {
	return assets.GetToolsDir("logs", "roundcube", rc.Name)
}

func (rc *Roundcube) SocketPath() string {
	name := fmt.Sprintf("php%s-roundcube-%s.sock", rc.PhpVersion, rc.Name)
	return assets.GetRunDir("php", name)
}

func (cfg *Config) DeployRoundcube() error {
	domainList, _, err := email.GetDomains(nil)
	if err != nil {
		return err
	}

	for _, domain := range domainList.Domains {
		if err := cfg.DeployRoundcubeDomain(domain); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) DeployRoundcubeMap(sel map[string]bool) error {
	domainList, _, err := email.GetDomains(sel)
	if err != nil {
		return err
	}

	for _, domain := range domainList.Domains {
		if err := cfg.DeployRoundcubeDomain(domain); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) DeployRoundcubeDomain(domain *email.Domain) error {
	cfg.Debugf("Enter config/roundcube.go (%s)", domain.Name)

	_, rcRel, err := cfg.Catalog.Get("roundcube", cfg.Roundcube)
	if err != nil {
		return err
	}

	if rcRel.Download.Directory == "" {
		return fmt.Errorf("missing Directory in Roundcube download")
	}

	rc := &Roundcube{
		DomainName: domain.Name,
		Name:       agent.MakeDBName(domain.Name),
		PhpVersion: cfg.Baseline.PHP,
		FQDN:       cfg.FQDN(), // for Postfix / Dovecot access
		SysAdmin:   cfg.SysAdmin,
		Locale:     cfg.Locale(),
		Directory:  rcRel.Download.Directory,
		Download:   &rcRel.Download,
	}

	rc.Password, err = utils.FetchPassword(20, "vmail", "db_password")
	if err != nil {
		return err
	}

	rc.DesKey, err = utils.FetchPassword(24, "roundcube-"+rc.Name, "des_key")
	if err != nil {
		return err
	}

	if err := cfg.RoundcubeDownload(rc); err != nil {
		return err
	}

	if err := cfg.RoundcubeExtract(rc); err != nil {
		return err
	}

	// TODO (later) replace logo, favicon in skins/elastic/images

	if err := cfg.RoundcubeDataDirs(rc); err != nil {
		return err
	}

	if err := cfg.RoundcubeConfigFiles(rc); err != nil {
		return err
	}

	if err := cfg.RoundcubeSQL(rc); err != nil {
		return err
	}

	if err := cfg.RoundcubeHideInstaller(rc); err != nil {
		return err
	}

	if err := cfg.EnsureCertificate(rc.WebMail()); err != nil {
		return err
	}

	if err := cfg.RoundcubeSetupVhost(rc); err != nil {
		return err
	}

	if status, err := cfg.SetCNAME(domain.Name, "webmail"); err != nil {
		return err
	} else if status != "" {
		cfg.Say(status)
	}

	cfg.Debugf("Leave config/roundcube.go (%s)", domain.Name)
	return nil
}

func (cfg *Config) RoundcubeDownload(rc *Roundcube) error {
	req := cfg.NewRequest()

	req.Downloads = append(req.Downloads, rc.Download)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) RoundcubeExtract(rc *Roundcube) error {
	req := cfg.NewRequest()

	extract := agent.File{
		Task:   "extract",
		Path:   assets.GetDownloadsDir(rc.Download.Filename),
		Target: rc.RootDir(),
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

func (cfg *Config) RoundcubeDataDirs(rc *Roundcube) error {
	req := cfg.NewRequest()

	dataDirs := []string{
		rc.LogsDir(),
		filepath.Join(rc.BaseDir(), "temp"),
		filepath.Join(rc.BaseDir(), "upload"),
	}
	for _, dir := range dataDirs {
		subDir := agent.File{
			Task:  "mkdir",
			Path:  dir,
			Mode:  "0750",
			User:  "www-data",
			Group: "www-data",
		}
		req.AddFile(&subDir)
	}

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) RoundcubeSQL(rc *Roundcube) error {
	req := cfg.NewRequest()

	cfgTmpl := filepath.Join("roundcube", "create_config.sql")
	stmts, err := assets.SQL(cfgTmpl, rc)
	if err != nil {
		return err
	}

	entry := agent.MySQL{
		Stmts:   stmts,
		Comment: "create roundcube (vmail) tables",
	}

	loader := agent.MySQL{
		DbName:  "rc_" + rc.Name,
		DbPath:  filepath.Join(rc.BaseDir(), "SQL", "mysql.initial.sql"),
		Comment: "mysql.initial.sql",
	}
	req.MySQLs = append(req.MySQLs, &entry, &loader)

	req.AddService("apache2")
	req.AddService("dovecot")

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) RoundcubeConfigFiles(rc *Roundcube) error {
	req := cfg.NewRequest()

	configContent, err := assets.Render("roundcube/config.inc.php", rc)
	if err != nil {
		return err
	}
	configFile := agent.File{
		Task:    "write",
		Path:    filepath.Join(rc.BaseDir(), "config", "config.inc.php"),
		Content: configContent,
		Mode:    "0644",
	}
	req.AddFile(&configFile)

	passwordContent, err := assets.Render("roundcube/password.inc.php", rc)
	if err != nil {
		return err
	}
	passwordFile := agent.File{
		Task:    "write",
		Path:    filepath.Join(rc.BaseDir(), "config", "password.inc.php"),
		Content: passwordContent,
		Mode:    "0644",
	}
	req.AddFile(&passwordFile)

	poolContent, err := assets.Render("roundcube/php-fpm-pool.conf", rc)
	if err != nil {
		return err
	}
	poolPath := cfg.PhpFpmPoolPath("roundcube-" + rc.Name)
	poolFile := agent.File{
		Task:    "write",
		Path:    poolPath,
		Content: poolContent,
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

func (cfg *Config) RoundcubeHideInstaller(rc *Roundcube) error {
	req := cfg.NewRequest()

	// Hiding means removing access for www-data:www-data
	installerTask := agent.File{
		Task:  "process",
		Path:  filepath.Join(rc.BaseDir(), "installer"),
		Mode:  "0700",
		User:  "root",
		Group: "root",
	}
	req.AddFile(&installerTask)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) RoundcubeSetupVhost(rc *Roundcube) error {
	req := cfg.NewRequest()

	vhostTmpl, err := assets.Render("roundcube/vhost.conf", rc)
	if err != nil {
		return err
	}

	vhostName := fmt.Sprintf("roundcube-%s.conf", rc.WebMail())
	vhostPath := assets.GetApacheEtcDir("sites-available", vhostName)
	vhostFile := agent.File{
		Task:    "write",
		Path:    vhostPath,
		Content: vhostTmpl,
		Service: "apache2",
	}
	req.AddFile(&vhostFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}
