package config

import (
	"fmt"
	"path/filepath"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/assets"
	"github.com/gd-tools/gd-tools/email"
)

type Mailer struct {
	VmailUID string
	VmailGID string
	MailPath string
}

func LoadMailer() (*Mailer, error) {
	vmailUser, err := agent.GetUserID("vmail")
	if err != nil {
		return nil, err
	}

	mailer := Mailer{
		VmailUID: vmailUser.UID,
		VmailGID: vmailUser.GID,
		MailPath: assets.GetToolsDir("data", "vmail"),
	}

	return &mailer, nil
}

func (cfg *Config) DeployAccounts() error {
	return cfg.DeployAccountMap(nil)
}

func (cfg *Config) DeployAccountMap(sel map[string]bool) error {
	cfg.Debugf("Enter config/accounts.go (%v)", sel)

	mailer, err := LoadMailer()
	if err != nil {
		return err
	}
	cfg.Mailer = mailer

	domainList, _, err := email.GetDomains(sel)
	if err != nil {
		return err
	}

	for _, domain := range domainList.Domains {
		if err := cfg.DeployAccountDomain(domain); err != nil {
			return err
		}
		if err := cfg.DeployAccountProxy(domain); err != nil {
			return err
		}
	}

	emailSANs := email.GetDomainSANs()
	if err := cfg.EnsureCertificate(cfg.FQDN(), emailSANs...); err != nil {
		return err
	}

	cfg.Debug("Leave config/accounts.go")
	return nil
}

func (cfg *Config) DeployAccountDomain(domain *email.Domain) error {
	req := cfg.NewRequest()

	domainDir := assets.GetToolsDir("data", "vmail", domain.Name)
	domainMkdir := agent.File{
		Task:  "mkdir",
		Path:  domainDir,
		Mode:  "0750",
		User:  "vmail",
		Group: "vmail",
	}
	req.AddFile(&domainMkdir)

	req.AddService("postfix")
	req.AddService("dovecot")

	if err := req.Send(); err != nil {
		return err
	}

	if err := cfg.ProxyDirs(domain); err != nil {
		return err
	}
	if err := cfg.ProxyCert(domain); err != nil {
		return err
	}
	if err := cfg.ProxyVhost(domain); err != nil {
		return err
	}

	for _, user := range domain.UserList {
		if err := cfg.EmailUser(user); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) EmailUser(user *email.User) error {
	req := cfg.NewRequest()

	userDir := assets.GetToolsDir("data", "vmail", user.Domain, user.Local)
	userSubDirs := []string{
		userDir,
		filepath.Join(userDir, "Maildir"),
		filepath.Join(userDir, "Maildir", "cur"),
		filepath.Join(userDir, "Maildir", "new"),
		filepath.Join(userDir, "Maildir", "tmp"),
	}
	for _, dir := range userSubDirs {
		userMkdir := agent.File{
			Task:  "mkdir",
			Path:  dir,
			Mode:  "0700",
			User:  "vmail",
			Group: "vmail",
		}
		req.AddFile(&userMkdir)
	}

	sqlStmts, err := assets.SQL("account/add_user.sql", user)
	if err != nil {
		return err
	}
	sqlCmd := agent.MySQL{
		Stmts:   sqlStmts,
		Comment: fmt.Sprintf("add vmail account for %s", user.Email()),
	}
	req.MySQLs = append(req.MySQLs, &sqlCmd)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) ProxyDirs(domain *email.Domain) error {
	req := cfg.NewRequest()

	proxyDir := assets.GetToolsDir("data", "mailproxy", domain.Name)
	data := struct {
		Domain string
		FQDN   string
		Root   string
		MXs    []string
	}{
		Domain: domain.Name,
		FQDN:   cfg.FQDN(),
		Root:   proxyDir,
		MXs:    []string{cfg.FQDN()},
	}
	if cfg.Spambarrier != "" {
		data.MXs = append(data.MXs, "*.spambarrier.de")
	}

	// Create dir for autoconfig.xml
	autoConfigDir := filepath.Join(proxyDir, "mail")
	autoConfigMkdir := agent.File{
		Task:  "mkdir",
		Path:  autoConfigDir,
		Mode:  "0755",
		User:  "root",
		Group: "root",
	}
	req.AddFile(&autoConfigMkdir)

	// create autoconfig.xml
	autoConfigTmpl, err := assets.Render("mailproxy/autoconfig.xml", data)
	if err != nil {
		return err
	}
	autoConfigFile := agent.File{
		Task:    "write",
		Path:    filepath.Join(autoConfigDir, "config-v1.1.xml"),
		Content: autoConfigTmpl,
		Mode:    "0644",
		User:    "root",
		Group:   "root",
	}
	req.AddFile(&autoConfigFile)

	// Create dir for autodiscover.xml
	autoDiscoverDir := filepath.Join(proxyDir, "autodiscover")
	autoDiscoverMkdir := agent.File{
		Task:  "mkdir",
		Path:  autoDiscoverDir,
		Mode:  "0755",
		User:  "root",
		Group: "root",
	}
	req.AddFile(&autoDiscoverMkdir)

	// create autodiscover.xml
	autoDiscoverTmpl, err := assets.Render("mailproxy/autodiscover.xml", data)
	if err != nil {
		return err
	}
	autoDiscoverFile := agent.File{
		Task:    "write",
		Path:    filepath.Join(autoDiscoverDir, "autodiscover.xml"),
		Content: autoDiscoverTmpl,
		Mode:    "0644",
		User:    "root",
		Group:   "root",
	}
	req.AddFile(&autoDiscoverFile)

	// Create dir for and mta-sts.txt
	mtaStsDir := filepath.Join(proxyDir, ".well-known")
	mtaStsMkdir := agent.File{
		Task:  "mkdir",
		Path:  mtaStsDir,
		Mode:  "0755",
		User:  "root",
		Group: "root",
	}
	req.AddFile(&mtaStsMkdir)

	// create mta-sts.txt
	mtaStsTmpl, err := assets.Render("mailproxy/mta-sts.txt", data)
	if err != nil {
		return err
	}
	mtaStsFile := agent.File{
		Task:    "write",
		Path:    filepath.Join(mtaStsDir, "mta-sts.txt"),
		Content: mtaStsTmpl,
		Mode:    "0644",
		User:    "root",
		Group:   "root",
	}
	req.AddFile(&mtaStsFile)

	// Create dir for all log files
	allLogsMkdir := agent.File{
		Task:  "mkdir",
		Path:  assets.GetToolsDir("logs", "mailproxy"),
		Mode:  "0755",
		User:  "root",
		Group: "root",
	}
	req.AddFile(&allLogsMkdir)

	// Create dir for domain log files
	logsMkdir := agent.File{
		Task:  "mkdir",
		Path:  assets.GetToolsDir("logs", "mailproxy", domain.Name),
		Mode:  "0755",
		User:  "www-data",
		Group: "www-data",
	}
	req.AddFile(&logsMkdir)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) ProxyCert(domain *email.Domain) error {
	certName := "autoconfig." + domain.Name
	sanList := []string{
		"autodiscover." + domain.Name,
		"mta-sts." + domain.Name,
	}

	if err := cfg.EnsureCertificate(certName, sanList...); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) ProxyVhost(domain *email.Domain) error {
	req := cfg.NewRequest()

	certName := "autoconfig." + domain.Name
	config := struct {
		Domain   string
		SysAdmin string
		RootDir  string
		CertDir  string
		LogsDir  string
	}{
		Domain:   domain.Name,
		SysAdmin: cfg.SysAdmin,
		RootDir:  assets.GetToolsDir("data", "mailproxy", domain.Name),
		CertDir:  assets.GetToolsDir("data", "certs", certName),
		LogsDir:  assets.GetToolsDir("logs", "mailproxy", domain.Name),
	}

	vHostTmpl, err := assets.Render("mailproxy/vhost.conf", config)
	if err != nil {
		return err
	}
	vHostName := fmt.Sprintf("site-autoconfig.%s.conf", domain.Name)
	vHostPath := assets.GetApacheEtcDir("sites-available", vHostName)
	vHostFile := agent.File{
		Task:    "write",
		Path:    vHostPath,
		Content: vHostTmpl,
		Service: "apache2",
	}
	req.AddFile(&vHostFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) DeployAccountProxy(domain *email.Domain) error {
	cnames := []string{
		"autoconfig",
		"autodiscover",
		"mta-sts",
		"imap",
		"smtp",
	}
	for _, alias := range domain.Aliases {
		cnames = append(cnames, alias)
	}

	for _, cname := range cnames {
		if status, err := cfg.SetCNAME(domain.Name, cname); err != nil {
			return err
		} else if status != "" {
			cfg.Say(status)
		}
	}

	return nil
}
