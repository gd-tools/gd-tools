package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gd-tools/gd-tools/config"
	"github.com/gd-tools/gd-tools/server"
	"github.com/gd-tools/gd-tools/utils"
)

// saveServer writes the persistent server configuration.
// It is called from the $(GD_TOOLS_BASE) directory.
// The newly created server will be a direct child path.
func saveServer(cfg *config.Config) error {
	if err := utils.EnsureBaseDir(); err != nil {
		return err
	}

	fqdn := cfg.FQDN()
	configPath := filepath.Join(fqdn, utils.ConfigFile)

	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("server %s exists - will not overwrite", fqdn)
	}

	khContent, khErr := os.ReadFile("known_hosts")

	if err := os.Mkdir(fqdn, 0o755); err != nil {
		return err
	}
	if err := os.Chdir(fqdn); err != nil {
		return err
	}

	if err := cfg.EnsureCA(); err != nil {
		return err
	}

	if khErr == nil {
		if err := os.WriteFile("known_hosts", khContent, 0o600); err != nil {
			return fmt.Errorf("failed to write known_hosts: %w", err)
		}
	}

	if _, _, err := cfg.RSAKeyPair(fqdn); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	if err := os.MkdirAll(config.ACME_Cert_Dir, 0o755); err != nil {
		return err
	}

	return nil
}
