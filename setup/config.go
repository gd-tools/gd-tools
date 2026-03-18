package setup

import (
	"github.com/gd-tools/gd-tools/config"
	"github.com/gd-tools/gd-tools/platform"
	"github.com/gd-tools/gd-tools/server"
	"github.com/gd-tools/gd-tools/utils"
	"github.com/urfave/cli/v2"
)

// buildConfig creates the initial server config from CLI flags and identity defaults.
func buildConfig(c *cli.Context, pf *platform.Platform, id *utils.Identity, host, domain string) config.Config {
	cfg := config.Config{
		Verbose: c.Bool("verbose"),
	}
	cfg.BaselineName = pf.Baseline.Name

	// Values for the operating system environment.
	if c.IsSet("swap-size") {
		cfg.SwapSize = c.String("swap-size")
	}
	addMounts(&cfg, c, pf)

	// Values that are taken from Identity.
	cfg.TimeZone = id.TimeZone
	cfg.Language = id.Language
	cfg.Region = id.Region
	cfg.RegTTL = id.RegTTL

	// Mandatory args for the command.
	cfg.HostName = host
	cfg.DomainName = domain

	// Values with defaults from Identity.
	if c.IsSet("company") {
		cfg.Company = c.String("company")
	} else {
		cfg.Company = id.Company
	}
	if c.IsSet("domain") {
		cfg.Domain = c.String("domain")
	} else {
		cfg.Domain = id.Domain
	}
	if c.IsSet("sysadmin") {
		cfg.SysAdmin = c.String("sysadmin")
	} else {
		cfg.SysAdmin = id.SysAdmin
	}
	if c.IsSet("help-url") {
		cfg.HelpURL = c.String("help-url")
	} else {
		cfg.HelpURL = id.HelpURL
	}
	if c.IsSet("dmarc") {
		cfg.DMARC = c.String("dmarc")
	} else {
		cfg.DMARC = id.DMARC
	}

	// Last but not least: provider credentials.
	cfg.Spambarrier = c.String("spambarrier")
	cfg.UbuntuPro = c.String("ubuntu-pro")
	cfg.HetznerToken = c.String("hetzner-dns")
	cfg.IonosToken = c.String("ionos-dns")
	cfg.CloudflareToken = c.String("cloudflare-dns")

	// Initialize the "unique DNS names" system.
	cfg.UsedFQDNs = []string{cfg.FQDN()}

	return cfg
}

// addMounts adds optional storage mounts from CLI flags.
func addMounts(cfg *config.Config, c *cli.Context, pf *platform.Platform) {
	if volume := c.String("hetzner-volume"); volume != "" {
		cfg.Mounts = append(cfg.Mounts, server.Mount{
			Provider: "hetzner",
			ID:       volume,
			Path:     pf.ToolsPath(),
		})
	}
}
