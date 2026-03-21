package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gd-tools/gd-tools/config"
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

	if _, err := cfg.LoadFile(configPath); err == nil {
		return fmt.Errorf("server %s exists - will not overwrite", fqdn)
	}

	if err := cfg.MkdirAll(fqdn, 0o755); err != nil {
		return err
	}
	if err := cfg.Chdir(fqdn); err != nil {
		return err
	}

	if err := cfg.EnsureCA(); err != nil {
		return err
	}

	if _, _, err := cfg.RSAKeyPair(fqdn); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	if err := cfg.MkdirAll(config.ACMECertDir, 0755); err != nil {
		return err
	}

	return nil
}
