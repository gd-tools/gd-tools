package hello

import (
	"github.com/gd-tools/gd-tools/config"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:  "hello",
	Usage: "Check mTLS communication with production server",
	Flags: []cli.Flag{
		config.FlagVerbose,
		config.FlagPort,
		&cli.IntFlag{
			Name:  "timeout",
			Usage: "connection timeout in seconds (overrides config)",
			Value: 0, // 0 = use config
		},
	},
	Action: Run,
}

func Run(c *cli.Context) error {
	cfg, err := config.ReadConfig(c)
	if err != nil {
		return err
	}
	defer cfg.Close()

	if err := cfg.EnsureCA(); err != nil {
		return err
	}

	if c.IsSet("timeout") {
		cfg.Timeout = c.Int("timeout")
	}

	cfg.Conn, err = agent.ConnectToAgent(cfg.FQDN(), cfg.Timeout, cfg.Verbose)
	if err != nil {
		return nil, err
	}

	req := cfg.NewRequest()
	req.Hello = "Hi, there."

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}
