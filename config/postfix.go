package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/assets"
	"github.com/gd-tools/gd-tools/email"
	"github.com/gd-tools/gd-tools/utils"
)

const (
	RoutingName = "routing"
	RoutingFile = RoutingName + ".json"
)

type Postfix struct {
	HostName   string
	DomainName string
	Password   string
	CertDir    string
	MailPath   string
	VmailUID   string
	VmailGID   string
	MilterIn   string
	MilterOut  string
}

type RelayDomain struct {
	Domain string `json:"domain"`
	Relay  string `json:"relay"`
}

type PolicyDomain struct {
	Domain string `json:"domain"`
	Policy string `json:"policy"`
}

type Routing struct {
	RelayDomains  []RelayDomain  `json:"relay_domains"`
	PolicyDomains []PolicyDomain `json:"policy_domains"`
}

func (cfg *Config) DeployPostfix() error {
	cfg.Debug("Enter config/postfix.go")

	// the Mailer is defined in config/accounts.go
	mailer, err := LoadMailer()
	if err != nil {
		return err
	}
	cfg.Mailer = mailer

	password, err := utils.FetchPassword(20, "vmail", "db_password")
	if err != nil {
		return err
	}

	cfg.Postfix = &Postfix{
		HostName:   cfg.FQDN(),
		DomainName: cfg.DomainName,
		Password:   password,
		CertDir:    assets.GetToolsDir("data", "certs", cfg.FQDN()),
		VmailUID:   mailer.VmailUID,
		VmailGID:   mailer.VmailGID,
		MailPath:   mailer.MailPath,
	}

	if cfg.Spambarrier != "" {
		cfg.Postfix.MilterIn = ""
	} else {
		cfg.Postfix.MilterIn = "inet:localhost:11332"
	}

	if err := cfg.PostfixTables(); err != nil {
		return err
	}

	if err := cfg.PostfixSASL(); err != nil {
		return err
	}

	if err := cfg.PostfixMaps(); err != nil {
		return err
	}

	if err := cfg.PostfixMain(); err != nil {
		return err
	}

	if err := cfg.PostfixMaster(); err != nil {
		return err
	}

	cfg.AddFirewall("25/tcp")
	cfg.AddFirewall("465/tcp")
	cfg.AddFirewall("587/tcp")

	if err := cfg.Save(); err != nil {
		return err
	}

	cfg.Debug("Leave config/postfix.go")
	return nil
}

func (cfg *Config) PostfixTables() error {
	req := cfg.NewRequest()

	cfNames := []string{
		"mailbox-domains.cf",
		"mailbox-maps.cf",
		"alias-maps.cf",
	}

	for _, name := range cfNames {
		cfTmpl, err := assets.Render("postfix/"+name, cfg.Postfix)
		if err != nil {
			return err
		}

		cfFile := agent.File{
			Task:    "write",
			Path:    assets.GetEtcDir("postfix", name),
			Content: cfTmpl,
			Mode:    "0644",
			Service: "postfix",
		}
		req.AddFile(&cfFile)
	}

	sqlTmpl, err := assets.SQL("postfix/create_tables.sql", cfg.Postfix)
	if err != nil {
		return err
	}
	sqlCmd := agent.MySQL{
		Stmts:   sqlTmpl,
		Comment: "create postfix (vmail) tables",
	}
	req.MySQLs = append(req.MySQLs, &sqlCmd)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) PostfixSASL() error {
	req := cfg.NewRequest()

	// collect the provider credentials
	var lines []string

	line, err := email.BrevoSASL()
	if err != nil {
		return err
	}
	if line != "" {
		lines = append(lines, line)
	}
	// add more providers here

	sort.Strings(lines)

	saslData := strings.Join(lines, "\n") + "\n"
	file := agent.File{
		Task:    "postmap",
		Path:    assets.GetEtcDir("postfix", "sasl_passwd"),
		Content: []byte(saslData),
		Mode:    "0600",
		Service: "postfix",
	}
	req.AddFile(&file)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) PostfixMaps() error {
	req := cfg.NewRequest()

	content, err := os.ReadFile(RoutingFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", RoutingFile, err)
	}

	var routing Routing
	if err := json.Unmarshal(content, &routing); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", RoutingFile, err)
	}

	var transportLines []string
	for _, transport := range routing.RelayDomains {
		switch transport.Relay {
		case "brevo":
			server, port := email.BrevoTarget()
			if err != nil {
				return err
			}
			line := fmt.Sprintf("%s smtp:[%s]:%d", transport.Domain, server, port)
			transportLines = append(transportLines, line)
		}
		// add more providers here
	}
	sort.Strings(transportLines)

	transportData := strings.Join(transportLines, "\n") + "\n"
	transportFile := agent.File{
		Task:    "postmap",
		Path:    assets.GetEtcDir("postfix", "transport"),
		Content: []byte(transportData),
		Mode:    "0600",
		Service: "postfix",
	}
	req.AddFile(&transportFile)

	var policyLines []string
	for _, policy := range routing.PolicyDomains {
		line := fmt.Sprintf("%s %s", policy.Domain, policy.Policy)
		policyLines = append(policyLines, line)
	}
	sort.Strings(policyLines)

	policyData := strings.Join(policyLines, "\n") + "\n"
	policyFile := agent.File{
		Task:    "postmap",
		Path:    assets.GetEtcDir("postfix", "tls_policy"),
		Content: []byte(policyData),
		Mode:    "0600",
		Service: "postfix",
	}
	req.AddFile(&policyFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) PostfixMain() error {
	req := cfg.NewRequest()

	mainContent, err := assets.Render("postfix/main.cf", cfg.Postfix)
	if err != nil {
		return err
	}

	mainFile := agent.File{
		Task:    "write",
		Path:    assets.GetEtcDir("postfix", "main.cf"),
		Content: mainContent,
		Backup:  true,
		Mode:    "0644",
		Service: "postfix",
	}
	req.AddFile(&mainFile)

	req.AddFirewall("25/tcp")
	req.AddFirewall("465/tcp")
	req.AddFirewall("587/tcp")

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) PostfixMaster() error {
	req := cfg.NewRequest()

	masterTmpl, err := assets.Render("postfix/master.cf", cfg.Postfix)
	if err != nil {
		return err
	}

	masterFile := agent.File{
		Task:    "write",
		Path:    assets.GetEtcDir("postfix", "master.cf"),
		Content: masterTmpl,
		Backup:  true,
		Mode:    "0644",
		Service: "postfix",
	}
	req.AddFile(&masterFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}
