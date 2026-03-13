package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/assets"
)

func (cfg *Config) CheckRemote(task string) bool {
	cmd := exec.Command("ssh", cfg.RootUser(), task)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (cfg *Config) RemoteScript(commands []string) error {
	if len(commands) == 0 {
		return fmt.Errorf("no remote commands provided")
	}

	script := strings.Join(commands, " && ")
	return cfg.RemoteCmd(script)
}

func (cfg *Config) RemoteCmd(command string) error {
	rootUser := cfg.RootUser()
	cfg.Sayf("ssh %s %q", rootUser, command)

	cmd := exec.Command("ssh", rootUser, command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ssh to %s failed: %w\nOutput:\n%s", rootUser, err, output)
	}

	return nil
}

func (cfg *Config) PushCerts() {
	if info, err := os.Stat("acme-certs"); err != nil || !info.IsDir() {
		return
	}

	dataPath := assets.GetToolsDir("data")
	if ok := cfg.CheckRemote("test -d " + dataPath); !ok {
		return
	}
	certPath := assets.GetToolsDir("data", "certs")

	if _, err := agent.RunCommand(
		"rsync",
		cfg.RsyncFlags(),
		"--chown=root:root",
		"--delete",
		"acme-certs/",
		cfg.RootUser()+":"+certPath,
	); err != nil {
		cfg.Sayf("failed to push ACME certs: %s (ignored)", err.Error())
	}
}
