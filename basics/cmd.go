package basics

import (
	"fmt"
	"io"

	"github.com/gd-tools/gd-tools/config"
	"github.com/gd-tools/gd-tools/utils"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:  "basics",
	Usage: "Prepare or update gd-tools base directory",
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

	bsc, err := utils.EnsureBasics()
	if err != nil {
		return err
	}

	if !basicsHasChanges(c) {
		printBasics(c.App.Writer, bsc)
		return nil
	}

	if c.IsSet("company") {
		bsc.Company = c.String("company")
	}
	if c.IsSet("domain") {
		bsc.Domain = c.String("domain")
	}
	if c.IsSet("sysadmin") {
		bsc.SysAdmin = c.String("sysadmin")
	}
	if c.IsSet("help-url") {
		bsc.HelpURL = c.String("help-url")
	}
	if c.IsSet("timezone") {
		bsc.TimeZone = c.String("timezone")
	}
	if c.IsSet("language") {
		bsc.Language = c.String("language")
	}
	if c.IsSet("region") {
		bsc.Region = c.String("region")
	}
	if c.IsSet("reg-ttl") {
		bsc.RegTTL = c.Int("reg-ttl")
	}

	if c.IsSet("dmarc") {
		bsc.DMARC = c.String("dmarc")
	}
	if bsc.DMARC == "" {
		bsc.DMARC = utils.DefaultDMARC
	}

	if err := bsc.Save(); err != nil {
		return err
	}

	return nil
}

func basicsHasChanges(c *cli.Context) bool {
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

func printBasics(w io.Writer, bsc *utils.Basics) {
	fmt.Fprintf(w, "%-12s  %s\n", "KEY", "VALUE")
	fmt.Fprintf(w, "%-12s  %s\n", "------------", "------------------------------")

	fmt.Fprintf(w, "%-12s  %s\n", "Company", bsc.Company)
	fmt.Fprintf(w, "%-12s  %s\n", "Domain", bsc.Domain)
	fmt.Fprintf(w, "%-12s  %s\n", "SysAdmin", bsc.SysAdmin)
	fmt.Fprintf(w, "%-12s  %s\n", "HelpURL", bsc.HelpURL)
	fmt.Fprintf(w, "%-12s  %s\n", "TimeZone", bsc.TimeZone)
	fmt.Fprintf(w, "%-12s  %s\n", "Language", bsc.Language)
	fmt.Fprintf(w, "%-12s  %s\n", "Region", bsc.Region)
	fmt.Fprintf(w, "%-12s  %d\n", "RegTTL", bsc.RegTTL)
	fmt.Fprintf(w, "%-12s  %s\n", "DMARC", bsc.DMARC)
}
