package wordpress

import (
	"fmt"
	"sort"

	"github.com/gd-tools/gd-tools/config"
	"github.com/gd-tools/gd-tools/email"
	"github.com/gd-tools/gd-tools/platform"
	"github.com/gd-tools/gd-tools/utils"
	"github.com/urfave/cli/v2"
)

var Describe = `The wordpress command prepares a WordPress instance for deployment.

WordPress is a free and open source blogging tool and
content management system (CMS) based on PHP and MySQL.

Homepage ...: https://wordpress.org/
Releases ...: https://wordpress.org/download/releases/
Installation: https://developer.wordpress.org/advanced-administration/before-install/howto-install/

In case the host part is 'www' the domain itself will also be
registered with Lets Encrypt, and a non-WWW scheme will be used.`

var Command = &cli.Command{
	Name:        "wordpress",
	Usage:       "Prepare a WordPress instance for deployment",
	Description: Describe,
	Flags: []cli.Flag{
		config.FlagVerbose,
		&cli.BoolFlag{
			Name:  "update",
			Usage: "update an existing WordPress instance",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "password for MySQL access",
		},
		&cli.StringFlag{
			Name:  "salt",
			Usage: "base for various secret values",
		},
		&cli.StringFlag{
			Name:  "admin-name",
			Usage: "name of WordPress admin",
		},
		&cli.StringFlag{
			Name:  "admin-email",
			Usage: "email address of WordPress admin",
		},
		&cli.StringFlag{
			Name:  "admin-pswd",
			Usage: "password of WordPress admin",
		},
		&cli.StringSliceFlag{
			Name:  "alias",
			Usage: "alias domain (only if host is 'www')",
		},
	},
	ArgsUsage: "<host> <domain>",
	Action:    Run,
}

func Run(c *cli.Context) error {
	cfg, err := config.ReadConfig(c)
	if err != nil {
		return err
	}

	if c.NArg() != 2 {
		return fmt.Errorf("Usage: gdt wordpress <host> <domain>")
	}
	host := c.Args().Get(0)
	domain := c.Args().Get(1)
	aliases := c.StringSlice("alias")
	sort.Strings(aliases)

	pf, err := platform.LoadPlatform()
	if err != nil {
		return err
	}

	wpProd, wpRel, err := pf.GetProduct("wordpress", "")
	if err != nil {
		return err
	}

	cliProd, _, err := pf.GetProduct("wp-cli", "")
	if err != nil {
		return err
	}
	wp := config.WordPress{
		HostName:     host,
		DomainName:   domain,
		Version:      wpProd.Default,
		WpCliVersion: cliProd.Default,
		Language:     cfg.Language,
		Region:       cfg.Region,
		Password:     c.String("password"),
		Directory:    wpRel.Download.Directory,
		Salt:         c.String("salt"),
		AdminName:    c.String("admin-name"),
		AdminEmail:   c.String("admin-email"),
		AdminPswd:    c.String("admin-pswd"),
		Aliases:      aliases,
	}

	if c.Bool("update") {
		if _, err := config.LoadWordPressList(&wp); err != nil {
			return err
		}
		cfg.Sayf("WordPress '%s' was updated", wp.FQDN())
		return nil
	}

	if err := cfg.CheckFQDN(wp.FQDN(), false); err != nil {
		return err
	}

	list, err := config.LoadWordPressList(nil)
	if err != nil {
		return err
	}

	if wp.Password == "" {
		if wp.Password, err = utils.CreatePassword(20); err != nil {
			return err
		}
	}

	if wp.Salt = c.String("salt"); wp.Salt == "" {
		if wp.Salt, err = utils.CreatePassword(30); err != nil {
			return err
		}
	}

	adminEmail := c.String("admin-email")
	if adminEmail == "" {
		adminEmail = "admin@" + domain
	}
	adminUser, err := email.MakeUser(adminEmail)
	if err != nil {
		return err
	}
	wp.AdminEmail = adminUser.Email()

	secrets, err := utils.LoadSecrets()
	if err != nil {
		return err
	}
	_, _, err = secrets.SetMailUser(adminUser.Email(), "")
	if err != nil {
		return err
	}

	if wp.AdminPswd == "" {
		if wp.AdminPswd, err = utils.CreatePassword(20); err != nil {
			return err
		}
	}

	list.Entries = append(list.Entries, &wp)

	if err := list.Save(); err != nil {
		return err
	}

	if err := cfg.CheckFQDN(wp.FQDN(), true); err != nil {
		return err
	}

	return nil
}
