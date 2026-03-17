package identity

import (
	"fmt"
	"io"

	"github.com/gd-tools/gd-tools/config"
	"github.com/gd-tools/gd-tools/utils"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:  "identity",
	Usage: "Prepare or update gd-tools common information",
	Flags: []cli.Flag{
		config.FlagVerbose,
		&cli.StringFlag{
			Name:  "company",
			Usage: "Name of your organisation",
		},
		&cli.StringFlag{
			Name:  "domain",
			Usage: "internet domain of your organisation",
		},
		&cli.StringFlag{
			Name:  "sysadmin",
			Usage: "System administrator email",
		},
		&cli.StringFlag{
			Name:  "help-url",
			Usage: "support URL, e.g. for Webmail",
		},
		&cli.StringFlag{
			Name:  "timezone",
			Usage: "Timezone, e.g. Europe/Berlin",
		},
		&cli.StringFlag{
			Name:  "language",
			Usage: "Language, e.g. de",
		},
		&cli.StringFlag{
			Name:  "region",
			Usage: "Region, e.g. DE",
		},
		&cli.IntFlag{
			Name:  "reg-ttl",
			Usage: "Regular Time-To-Live, cache-time for DNS",
			Value: utils.DefaultRegTTL,
		},
		&cli.StringFlag{
			Name:  "dmarc",
			Usage: "Standard DMARC value",
		},
	},
	Action: Run,
}

func Run(c *cli.Context) error {
	if err := utils.EnsureBaseDir(); err != nil {
		return err
	}

	id, err := utils.EnsureIdentity()
	if err != nil {
		return err
	}

	if !identityHasChanges(c) {
		printIdentity(c.App.Writer, id)
		return nil
	}

	if c.IsSet("company") {
		id.Company = c.String("company")
	}
	if c.IsSet("domain") {
		id.Domain = c.String("domain")
	}
	if c.IsSet("sysadmin") {
		id.SysAdmin = c.String("sysadmin")
	}
	if c.IsSet("help-url") {
		id.HelpURL = c.String("help-url")
	}
	if c.IsSet("timezone") {
		id.TimeZone = c.String("timezone")
	}
	if c.IsSet("language") {
		id.Language = c.String("language")
	}
	if c.IsSet("region") {
		id.Region = c.String("region")
	}
	if c.IsSet("reg-ttl") {
		id.RegTTL = c.Int("reg-ttl")
	}

	if c.IsSet("dmarc") {
		id.DMARC = c.String("dmarc")
	}
	if id.DMARC == "" {
		id.DMARC = utils.DefaultDMARC
	}

	if err := id.Save(); err != nil {
		return err
	}

	return nil
}

func identityHasChanges(c *cli.Context) bool {
	return c.IsSet("company") ||
		c.IsSet("domain") ||
		c.IsSet("sysadmin") ||
		c.IsSet("help-url") ||
		c.IsSet("timezone") ||
		c.IsSet("language") ||
		c.IsSet("region") ||
		c.IsSet("reg-ttl") ||
		c.IsSet("dmarc")
}

func printIdentity(w io.Writer, id *utils.Identity) {
	fmt.Fprintf(w, "%-12s  %s\n", "KEY", "VALUE")
	fmt.Fprintf(w, "%-12s  %s\n", "------------", "------------------------------")

	fmt.Fprintf(w, "%-12s  %s\n", "Company", id.Company)
	fmt.Fprintf(w, "%-12s  %s\n", "Domain", id.Domain)
	fmt.Fprintf(w, "%-12s  %s\n", "SysAdmin", id.SysAdmin)
	fmt.Fprintf(w, "%-12s  %s\n", "HelpURL", id.HelpURL)
	fmt.Fprintf(w, "%-12s  %s\n", "TimeZone", id.TimeZone)
	fmt.Fprintf(w, "%-12s  %s\n", "Language", id.Language)
	fmt.Fprintf(w, "%-12s  %s\n", "Region", id.Region)
	fmt.Fprintf(w, "%-12s  %d\n", "RegTTL", id.RegTTL)
	fmt.Fprintf(w, "%-12s  %s\n", "DMARC", id.DMARC)
}
