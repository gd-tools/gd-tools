package deploy

import (
	"fmt"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/config"
	"github.com/urfave/cli/v2"
)

var Describe = `The deploy command is used to update the production system.

Valid arguments are:
  - system (includes bootstrap ... backup)
  - bootstrap
  - packages
  - filesystem
  - database (Mariadb and redis)
  - php      (currently 8.3)
  - apache   (with ACME certs)
  - redirect
  - backup   (Borg)

  - email (includes dovecot ... roundcube)
  - postfix
  - dovecot
  - opendkim
  - accounts
  - roundcube

  - nextcloud (all)
  - nc-<id>
  - ocis (ownCloud Infinite Scale)

  - wordpress (all)
  - wp-<id>

  - rustdesk`

var Command = &cli.Command{
	Name:        "deploy",
	Usage:       "Deploy some or all system and/or project components",
	Description: Describe,
	Flags: []cli.Flag{
		config.FlagVerbose,
		config.FlagForce,
		config.FlagPort,
		config.FlagSkipDNS,
		config.FlagSkipMX,
		&cli.BoolFlag{
			Name:  "upgrade",
			Usage: "update the system (apt-get update/upgrade)",
		},
	},
	ArgsUsage: "[<component>]...",
	BashComplete: func(c *cli.Context) {
		fmt.Println("system")
		fmt.Println("bootstrap")
		fmt.Println("packages")
		fmt.Println("filesystem")
		fmt.Println("php")
		fmt.Println("apache")
		fmt.Println("redirect")
		fmt.Println("backup")

		fmt.Println("email")
		fmt.Println("postfix")
		fmt.Println("dovecot")
		fmt.Println("opendkim")
		fmt.Println("accounts")
		fmt.Println("roundcube")

		// Cloud storage
		ncList, err := agent.LoadNextcloudList(nil)
		if err == nil && len(ncList.Entries) > 0 {
			fmt.Println("nextcloud")
			for _, nc := range ncList.Entries {
				fmt.Printf("nc-%s\n", nc.FQDN())
			}
		}
		if oc, err := config.LoadOCIS(); err == nil && oc != nil {
			fmt.Println("ocis")
		}

		// Content management
		wpList, err := config.LoadWordPressList(nil)
		if err == nil && len(wpList.Entries) > 0 {
			fmt.Println("wordpress")
			for _, wp := range wpList.Entries {
				fmt.Printf("wp-%s\n", wp.FQDN())
			}
		}

		// Other useful apps
		if rd, err := config.LoadRustDesk(); err == nil && rd != nil {
			fmt.Println("rustdesk")
		}
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

	cfg.Conn, err = agent.ConnectToAgent(cfg.FQDN(), cfg.Timeout, cfg.Verbose)
	if err != nil {
		return nil, err
	}

	ncList, err := agent.LoadNextcloudList(nil)
	if err != nil {
		return err
	}

	wpList, err := config.LoadWordPressList(nil)
	if err != nil {
		return err
	}

	var newArgs []string
	for _, arg := range c.Args().Slice() {
		switch arg {
		case "system":
			newArgs = append(newArgs,
				"bootstrap",
				"packages",
				"filesystem",
				"database",
				"php",
				"apache",
				"redirect",
				"backup",
			)
		case "email":
			newArgs = append(newArgs,
				"postfix",
				"dovecot",
				"opendkim",
				"accounts",
				"roundcube",
			)
		default:
			newArgs = append(newArgs, arg)
		}
	}

	for _, arg := range newArgs {
		switch arg {
		// System modules
		case "bootstrap":
			if err := cfg.DeployBootstrap(); err != nil {
				return err
			}
		case "packages":
			if err = cfg.DeployPackages(c.Bool("upgrade")); err != nil {
				return err
			}
		case "filesystem":
			if err := cfg.DeployMounts(); err != nil {
				return err
			}
			if err := cfg.DeployFilesystem(); err != nil {
				return err
			}
		case "database":
			if err := cfg.DeployDatabase(); err != nil {
				return err
			}
		case "php":
			if err := cfg.DeployPHP(); err != nil {
				return err
			}
		case "apache":
			if err := cfg.DeployApache(); err != nil {
				return err
			}
		case "redirect":
			if err := cfg.DeployRedirect(); err != nil {
				return err
			}
		case "backup":
			if err := cfg.DeployBackup(); err != nil {
				return err
			}

		// Email modules
		case "postfix":
			if err := cfg.DeployPostfix(); err != nil {
				return err
			}
		case "dovecot":
			if err := cfg.DeployDovecot(); err != nil {
				return err
			}
		case "opendkim":
			if err := cfg.DeployOpenDKIM(); err != nil {
				return err
			}
		case "accounts":
			if err := cfg.DeployAccounts(); err != nil {
				return err
			}
		case "roundcube":
			if err := cfg.DeployRoundcube(); err != nil {
				return err
			}

		// Cloud storage
		case "ocis":
			oc, err := config.LoadOCIS()
			if err != nil {
				return err
			}
			if oc == nil {
				return fmt.Errorf("OCIS is not installed")
			}
			if err := cfg.DeployOCIS(oc); err != nil {
				return err
			}

		// Other useful apps
		case "rustdesk":
			rd, err := config.LoadRustDesk()
			if err != nil {
				return err
			}
			if rd == nil {
				return fmt.Errorf("RustDesk is not installed")
			}
			if err := cfg.DeployRustDesk(rd); err != nil {
				return err
			}

			// the switch with fixed names ends here
		}

		// Nextcloud instances
		for _, nc := range ncList.Entries {
			if arg == "nextcloud" || arg == nc.FQDN() {
				if err := cfg.DeployNextcloud(nc); err != nil {
					return err
				}
			}
		}

		// WordPress instances
		for _, wp := range wpList.Entries {
			if arg == "wordpress" || arg == wp.FQDN() {
				if err := cfg.DeployWordPress(wp); err != nil {
					return err
				}
			}
		}

		//  unknown arg, ignore silently
	}

	return nil
}
