package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/assets"
)

const (
	CommonNameCA     = "/CN=MyCA"
	CommonNameDev    = "/CN=dev"
	CommonNameProd   = "/CN=prod"
	SerialName       = "CA/serial"
	CaKeyName        = "CA/ca.key"
	CaCrtName        = "CA/ca.crt"
	ClientKeyName    = "CA/client.key"
	ClientCsrName    = "CA/client.csr"
	ClientCrtName    = "CA/client.crt"
	ServerKeyName    = "CA/server.key"
	ServerConfigName = "CA/server.config"
	ServerCsrName    = "CA/server.csr"
	ServerCrtName    = "CA/server.crt"
)

// AddHostToCA ensures a local Certificate Authority (CA) and generates
// a server certificate with Subject Alternative Names (SANs).
//
// The provided fqdn is added to the SAN list if not already present.
// The list is sorted to ensure idempotent certificate configuration.
func (cfg *Config) SetupCA() error {
	if err := os.MkdirAll("CA", 0700); err != nil {
		return fmt.Errorf("failed to mkdir CA: %w", err)
	}

	if ok := checkFile(SerialName); !ok {
		if err := os.WriteFile(SerialName, []byte("1000\n"), 0600); err != nil {
			return fmt.Errorf("failed to write %s: %w", SerialName, err)
		}
	}

	if ok := checkFile(CaKeyName); !ok {
		if _, err := agent.RunCommand(
			"openssl",
			"genrsa",
			"-out", CaKeyName,
			"2048",
		); err != nil {
			return err
		}
	}

	if ok := checkFile(CaCrtName); !ok {
		if _, err := agent.RunCommand(
			"openssl",
			"req", "-x509", "-new", "-nodes",
			"-key", CaKeyName,
			"-subj", CommonNameCA,
			"-days", "3650",
			"-out", CaCrtName,
		); err != nil {
			return err
		}
	}

	if ok := checkFile(ClientKeyName); !ok {
		if _, err := agent.RunCommand(
			"openssl",
			"genrsa",
			"-out", ClientKeyName,
			"2048",
		); err != nil {
			return err
		}
	}

	if ok := checkFile(ClientCsrName); !ok {
		if _, err := agent.RunCommand(
			"openssl",
			"req", "-new",
			"-key", ClientKeyName,
			"-subj", CommonNameDev,
			"-out", ClientCsrName,
		); err != nil {
			return err
		}
	}

	if ok := checkFile(ClientCrtName); !ok {
		if _, err := agent.RunCommand(
			"openssl",
			"x509", "-req",
			"-in", ClientCsrName,
			"-CA", CaCrtName,
			"-CAkey", CaKeyName,
			"-CAserial", SerialName,
			"-out", ClientCrtName,
			"-days", "3650",
		); err != nil {
			return err
		}
	}

	if ok := checkFile(ServerKeyName); !ok {
		if _, err := agent.RunCommand(
			"openssl",
			"genrsa",
			"-out", ServerKeyName,
			"2048",
		); err != nil {
			return err
		}
	}

	hostEntries := []string{
		"DNS.1 = " + cfg.FQDN(),
	}

	data := struct {
		HostEntries []string
	}{
		HostEntries: hostEntries,
	}
	configTmpl, err := assets.Render("system/ca.server.config", data)
	if err != nil {
		return err
	}

	existing, err := os.ReadFile(ServerConfigName)
	if err == nil && bytes.Equal(existing, configTmpl) {
		return nil
	}

	if err := os.WriteFile(ServerConfigName, configTmpl, 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", ServerConfigName, err)
	}

	if _, err := agent.RunCommand(
		"openssl",
		"req", "-new",
		"-key", ServerKeyName,
		"-subj", CommonNameProd,
		"-out", ServerCsrName,
		"-config", ServerConfigName,
	); err != nil {
		return err
	}

	if _, err := agent.RunCommand(
		"openssl",
		"x509", "-req",
		"-in", ServerCsrName,
		"-CA", CaCrtName,
		"-CAkey", CaKeyName,
		"-CAserial", SerialName,
		"-out", ServerCrtName,
		"-days", "3650",
		"-extfile", ServerConfigName,
		"-extensions", "v3_req",
	); err != nil {
		return err
	}

	return nil
}

func checkFile(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.Size() > 0
}
