package config

import (
	"path/filepath"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/email"
	"github.com/gd-tools/gd-tools/releases"
	"github.com/gd-tools/gd-tools/templates"
	"github.com/gd-tools/gd-tools/utils"
)

func (cfg *Config) SieveBefore(paths ...string) string {
	sieve := releases.GetToolsDir("data", "sieve_before")
	if len(paths) == 0 {
		return sieve
	}
	return filepath.Join(append([]string{sieve}, paths...)...)
}

func (cfg *Config) SieveAfter(paths ...string) string {
	sieve := releases.GetToolsDir("data", "sieve_after")
	if len(paths) == 0 {
		return sieve
	}
	return filepath.Join(append([]string{sieve}, paths...)...)
}

func (cfg *Config) Domains() ([]*email.Domain, error) {
	domainList, _, err := email.GetDomains(nil)
	if err != nil {
		return nil, err
	}

	return domainList.Domains, nil
}

func (cfg *Config) DeployDovecot() error {
	cfg.Debug("Enter config/dovecot.go")

	// the Mailer is defined in pkg/config/accounts.go
	mailer, err := LoadMailer()
	if err != nil {
		return err
	}
	cfg.Mailer = mailer

	cfg.CertDir = releases.GetToolsDir("data", "certs", cfg.FQDN())

	cfg.Password, err = utils.FetchPassword(20, "vmail", "db_password")
	if err != nil {
		return err
	}

	if err := cfg.DovecotTables(); err != nil {
		return err
	}

	if err := cfg.ApacheDHparam("dovecot"); err != nil {
		return err
	}

	if err := cfg.DovecotFiles(); err != nil {
		return err
	}

	cfg.AddFirewall("993/tcp")
	if err := cfg.Save(); err != nil {
		return err
	}

	cfg.Debug("Leave config/dovecot.go")
	return nil
}

func (cfg *Config) DovecotTables() error {
	req := cfg.NewRequest()

	tmpl := filepath.Join("dovecot", "create_users.sql")
	stmts, err := templates.SQL(tmpl, cfg.Verbose, cfg)
	if err != nil {
		return err
	}
	entry := agent.MySQL{
		Stmts:   stmts,
		Comment: "create dovecot (vmail) tables",
	}
	req.MySQLs = append(req.MySQLs, &entry)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) DovecotFiles() error {
	req := cfg.NewRequest()

	beforeMkdir := agent.File{
		Task:  "mkdir",
		Path:  cfg.SieveBefore(),
		Mode:  "0755",
		User:  "vmail",
		Group: "vmail",
	}
	req.AddFile(&beforeMkdir)

	afterMkdir := agent.File{
		Task:  "mkdir",
		Path:  cfg.SieveAfter(),
		Mode:  "0755",
		User:  "vmail",
		Group: "vmail",
	}
	req.AddFile(&afterMkdir)

	if cfg.Spambarrier != "" {
		spamName := "10-spambarrier.sieve"
		spamTmpl := filepath.Join("dovecot/sieve_before", spamName)
		spamData, err := templates.Parse(spamTmpl, cfg.Verbose, cfg)
		if err != nil {
			return err
		}
		req.AddFile(&agent.File{
			Task:    "write",
			Path:    cfg.SieveBefore(spamName),
			Content: spamData,
			Mode:    "0644",
			User:    "vmail",
			Group:   "vmail",
			Service: "dovecot",
		})
	}

	forwardName := "20-forward.sieve"
	forwardTmpl := filepath.Join("dovecot/sieve_after", forwardName)
	forwardData, err := templates.Parse(forwardTmpl, cfg.Verbose, cfg)
	if err != nil {
		return err
	}
	req.AddFile(&agent.File{
		Task:    "write",
		Path:    cfg.SieveAfter(forwardName),
		Content: forwardData,
		Mode:    "0644",
		User:    "vmail",
		Group:   "vmail",
		Service: "dovecot",
	})

	files := []string{
		"conf.d/10-auth.conf",
		"conf.d/10-mail.conf",
		"conf.d/10-master.conf",
		"conf.d/10-ssl.conf",
		"conf.d/20-imap-lastlogin.conf",
		"conf.d/20-lmtp.conf",
	}

	for _, name := range files {
		tmpl := filepath.Join("dovecot", name)
		content, err := templates.Parse(tmpl, cfg.Verbose, cfg)
		if err != nil {
			return err
		}

		req.AddFile(&agent.File{
			Task:    "write",
			Path:    releases.GetEtcDir("dovecot", name),
			Content: content,
			Backup:  true,
			Mode:    "0644",
			Service: "dovecot",
		})
	}

	req.AddFirewall("993/tcp")

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}
