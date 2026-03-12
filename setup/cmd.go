package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/assets"
	"github.com/gd-tools/gd-tools/config"
	"github.com/gd-tools/gd-tools/utils"
	"github.com/urfave/cli/v2"
)

const (
	DefaultBaseline = "noble-8.3-2.4" // for new servers
)

var Describe = `Initialize a new production server

Creates the initial configuration on the development workstation.
Does not modify the production server (yet).

Calling 'gdt setup' without arguments lists alternatives.

For detailed documentation and usage examples, see:
https://github.com/gd-tools/gd-tools/wiki/10-Setup`

var (
	FlagBaseline = &cli.StringFlag{
		Name:  "baseline",
		Usage: "operating environment for this production server",
		Value: DefaultBaseline,
	}
	FlagHetznerVolume = &cli.StringFlag{
		Name:  "hetzner-volume",
		Usage: "add a Hetzner Cloud Volume for /var/gd-tools",
	}
	FlagRaidDevice = &cli.StringFlag{
		Name:  "raid-device",
		Usage: "add a /dev/mdX RAID device for /var/gd-tools",
	}
	FlagSwapSize = &cli.StringFlag{
		Name:  "swap-size",
		Usage: "e.g. '4G' - create or verify swapfile",
		Value: "0",
	}
	FlagDMARC = &cli.StringFlag{
		Name:  "dmarc",
		Usage: "DMARC value (p=quarantine; pct=100; adkim=s; aspf=s)",
	}
	FlagCompany = &cli.StringFlag{
		Name:  "company",
		Usage: "Company name, used e.g. for Webmail",
	}
	FlagSysAdmin = &cli.StringFlag{
		Name:  "sysadmin",
		Usage: "System Administrator email address",
	}
	FlagHelpURL = &cli.StringFlag{
		Name:  "help-url",
		Usage: "Support URL for this server",
	}
	FlagUbuntuPro = &cli.StringFlag{
		Name:  "ubuntu-pro",
		Usage: "attach Ubuntu Pro subscription using the provided token",
	}
	FlagSpamBarrier = &cli.StringFlag{
		Name:  "spambarrier",
		Usage: "SpamBarrier API key for inbound email",
	}
	FlagHetznerDNS = &cli.StringFlag{
		Name:  "hetzner-dns",
		Usage: "configure Hetzner Cloud DNS API token for declarative DNS management",
	}
	FlagIonosDNS = &cli.StringFlag{
		Name:  "ionos-dns",
		Usage: "configure IONOS DNS API token for declarative DNS management",
	}
	FlagCloudflareDNS = &cli.StringFlag{
		Name:  "cloudflare-dns",
		Usage: "configure Cloudflare DNS API token for declarative DNS management",
	}
)

var Command = &cli.Command{
	Name:        "setup",
	Usage:       "Initialize a new production server",
	Description: Describe,
	Flags: []cli.Flag{
		config.FlagVerbose,
		config.FlagDry,
		FlagBaseline,
		FlagHetznerVolume,
		FlagRaidDevice,
		FlagSwapSize,
		FlagDMARC,
		FlagCompany,
		FlagSysAdmin,
		FlagHelpURL,
		FlagUbuntuPro,
		FlagSpamBarrier,
		FlagHetznerDNS,
		FlagIonosDNS,
		FlagCloudflareDNS,
	},
	ArgsUsage: "<host> <domain>",
	Action:    Run,
}

func Run(c *cli.Context) error {
	catalog, err := assets.LoadCatalog()
	if err != nil {
		return err
	}

	if c.NArg() == 0 {
		fmt.Println("Avaliable baselines:")
		for _, bl := range catalog.Baselines {
			fmt.Printf("Baseline ....: %s\n", bl.Name)
			fmt.Printf("  Ubuntu ....: %s\n", bl.Ubuntu)
			fmt.Printf("  PHP .......: %s\n", bl.PHP)
			fmt.Printf("  Dovecot ...: %s\n", bl.Dovecot)
			fmt.Println()
		}
		return nil
	}

	baseline, err := catalog.GetBaseline(c.String("baseline"))
	if err != nil {
		return err
	}

	basics, err := utils.ReadBasics()
	if err != nil {
		return err
	}

	if c.NArg() != 2 {
		return fmt.Errorf("expected arguments: <host> <domain>")
	}
	host := c.Args().Get(0)
	domain := c.Args().Get(1)

	cfg := config.Config{
		Verbose:         c.Bool("verbose"),
		Dry:             c.Bool("dry"),
		BaselineName:    baseline.Name,
		TimeZone:        basics.TimeZone,
		Language:        basics.Language,
		Region:          basics.Region,
		RegTTL:          basics.RegTTL,
		HostName:        host,
		DomainName:      domain,
		SwapSize:        c.String("swap-size"),
		DMARC:           c.String("dmarc"),
		SysAdmin:        c.String("sysadmin"),
		Company:         c.String("company"),
		HelpURL:         c.String("help-url"),
		Spambarrier:     c.String("spambarrier"),
		UbuntuPro:       c.String("ubuntu-pro"),
		HetznerToken:    c.String("hetzner-dns"),
		IonosToken:      c.String("ionos-dns"),
		CloudflareToken: c.String("cloudflare-dns"),
	}

	fqdn := cfg.FQDN()
	cfg.UsedFQDNs = []string{fqdn}

	configPath := filepath.Join(fqdn, config.ConfigName)

	if cfg.DMARC == "" {
		cfg.DMARC = basics.DMARC
	}
	if cfg.Company == "" {
		cfg.Company = basics.Company
	}
	if cfg.SysAdmin == "" {
		cfg.SysAdmin = basics.SysAdmin
	}
	if cfg.HelpURL == "" {
		cfg.HelpURL = basics.HelpURL
	}

	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("Server %s exists - will not overwrite", fqdn)
	}

	reservedNames := []string{
		"autoconfig",
		"autodiscover",
		"mta-sts",
		"imap",
		"smtp",
		"vmail",
		"webmail",
		"www",
	}
	for _, name := range reservedNames {
		if host == name {
			return fmt.Errorf("hostname '%s' is reserved", name)
		}
	}

	// read default postfix-routing and known_hosts
	routing, err := os.ReadFile(config.RoutingName)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", config.RoutingName, err)
	}
	khContent, khErr := os.ReadFile("known_hosts")

	// check for filesystems to be mounted
	// N.B. mounts given here are mutually exclusive
	if volume := c.String("hetzner-volume"); volume != "" {
		mount := agent.Mount{
			Provider: "Hetzner",
			ID:       volume,
			Dir:      assets.GetToolsDir(""),
		}
		cfg.Mounts = append(cfg.Mounts, &mount)
	} else if device := c.String("raid-device"); device != "" {
		mount := agent.Mount{
			Provider: "RAID",
			ID:       device,
			Dir:      assets.GetToolsDir(""),
		}
		cfg.Mounts = append(cfg.Mounts, &mount)
	}

	if cfg.Dry {
		cfg2 := cfg
		cfg2.Spambarrier = "***"
		cfg2.UbuntuPro = "***"
		cfg2.HetznerToken = "***"
		cfg2.IonosToken = "***"
		cfg2.CloudflareToken = "***"

		content, err := json.MarshalIndent(cfg2, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal %s: %w", fqdn, err)
		}
		cfg.Sayf("Config: >%s<", string(content))

		return nil
	}

	// From here on, we operate inside the host directory on the development workstation.
	if err := os.Mkdir(fqdn, 0755); err != nil {
		return err
	}
	if err := os.Chdir(fqdn); err != nil {
		return err
	}

	if err := cfg.SetupCA(); err != nil {
		return err
	}

	if err := os.WriteFile(config.RoutingName, routing, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", config.RoutingName, err)
	}

	if khErr == nil {
		if err := os.WriteFile("known_hosts", khContent, 0600); err != nil {
			return fmt.Errorf("failed to write known_hosts: %w", err)
		}
	}

	if _, _, err := utils.GetRSAKeyPair(fqdn); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	if err := os.MkdirAll(config.ACME_Cert_Dir, 0755); err != nil {
		return err
	}

	return nil
}
