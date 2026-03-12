package main

import (
	"fmt"
	"os"

	"github.com/gd-tools/gd-tools/account"
	"github.com/gd-tools/gd-tools/backup"
	"github.com/gd-tools/gd-tools/bash"
	"github.com/gd-tools/gd-tools/basics"
	"github.com/gd-tools/gd-tools/brevo"
	"github.com/gd-tools/gd-tools/cert"
	"github.com/gd-tools/gd-tools/deploy"
	"github.com/gd-tools/gd-tools/domain"
	"github.com/gd-tools/gd-tools/hello"
	"github.com/gd-tools/gd-tools/install"
	"github.com/gd-tools/gd-tools/login"
	//"github.com/gd-tools/gd-tools/mediawiki"
	"github.com/gd-tools/gd-tools/nextcloud"
	"github.com/gd-tools/gd-tools/ocis"
	"github.com/gd-tools/gd-tools/redirect"
	"github.com/gd-tools/gd-tools/release"
	"github.com/gd-tools/gd-tools/rustdesk"
	"github.com/gd-tools/gd-tools/setup"
	"github.com/gd-tools/gd-tools/ssh"
	"github.com/gd-tools/gd-tools/status"
	"github.com/gd-tools/gd-tools/sync"
	"github.com/gd-tools/gd-tools/update"
	//"github.com/gd-tools/gd-tools/wbce"
	"github.com/gd-tools/gd-tools/wordpress"
	"github.com/urfave/cli/v2"
)

var (
	version string // will be loaded via -ldflags
)

func main() {
	app := &cli.App{
		Name:     "gdt",
		Version:  fmt.Sprintf("%s", version),
		Usage:    "gdt (short for gd-tools or go-deployment-tools) helps managing Ubuntu/Debian servers",
		Commands: getCommands(),
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		EnableBashCompletion: true,
		BashComplete: func(c *cli.Context) {
			fmt.Println("account")
			fmt.Println("backup")
			fmt.Println("bash")
			fmt.Println("basics")
			fmt.Println("brevo")
			fmt.Println("cert")
			fmt.Println("deploy")
			fmt.Println("domain")
			fmt.Println("hello")
			fmt.Println("install")
			fmt.Println("login")
			fmt.Println("nextcloud")
			fmt.Println("ocis")
			fmt.Println("redirect")
			fmt.Println("release")
			fmt.Println("rustdesk")
			fmt.Println("setup")
			fmt.Println("ssh")
			fmt.Println("status")
			fmt.Println("sync")
			fmt.Println("update")
			//fmt.Println("wbce")
			fmt.Println("wordpress")
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getCommands() []*cli.Command {
	return []*cli.Command{
		account.Command,
		backup.Command,
		bash.Command,
		basics.Command,
		brevo.Command,
		cert.Command,
		deploy.Command,
		domain.Command,
		hello.Command,
		install.Command,
		login.Command,
		//mediawiki.Command,
		nextcloud.Command,
		ocis.Command,
		redirect.Command,
		release.Command,
		rustdesk.Command,
		setup.Command,
		ssh.Command,
		status.Command,
		sync.Command,
		update.Command,
		//wbce.Command,
		wordpress.Command,
	}
}
