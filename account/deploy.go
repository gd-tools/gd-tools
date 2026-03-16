package account

import (
	"fmt"
	"strings"

	"github.com/gd-tools/gd-tools/config"
	"github.com/gd-tools/gd-tools/email"
	"github.com/urfave/cli/v2"
)

var DeployCommand = &cli.Command{
	Name:  "deploy",
	Usage: "update account data on the production server",
	Flags: []cli.Flag{
		config.FlagVerbose,
		config.FlagPort,
		config.FlagSkipDNS,
		config.FlagSkipMX,
		&cli.BoolFlag{
			Name:  "webmail",
			Usage: "deploy webmail in addition to proxies",
		},
		&cli.BoolFlag{
			Name:  "maintenance",
			Usage: "reverse user passwords for admin access",
		},
	},
	ArgsUsage: "[<domain> ...]",
	BashComplete: func(c *cli.Context) {
		domList, _, err := email.GetDomains(nil)
		if err != nil {
			return
		}
		for _, dom := range domList.Domains {
			fmt.Fprintln(c.App.Writer, dom.Name)
			for _, alias := range dom.Aliases {
				fmt.Fprintln(c.App.Writer, alias)
			}
		}
	},
	Action: DeployRun,
}

func DeployRun(c *cli.Context) error {
	cfg, err := config.ReadConfig(c)
	if err != nil {
		return err
	}
	defer cfg.Close()

	if err := cfg.EnsureCA(); err != nil {
		return err
	}

	cfg.Conn, err = agent.ConnectToAgent(cfg.FQDN(), cfg.Timeout, cfg.Verbose)
	if err != nil {
		return nil, err
	}

	sel := make(map[string]bool)
	for _, a := range c.Args().Slice() {
		a = strings.TrimSpace(a)
		if a != "" {
			sel[a] = true
		}
	}

	if err := cfg.DeployAccountMap(sel); err != nil {
		return err
	}

	if c.Bool("webmail") {
		if err := cfg.DeployRoundcubeMap(sel); err != nil {
			return err
		}
	}

	return nil
}
