package setup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gd-tools/gd-tools/config"
	"github.com/gd-tools/gd-tools/platform"
)

func TestSaveServer(t *testing.T) {
	base := t.TempDir()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	if err := os.Chdir(base); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("known_hosts", []byte("example host key\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{}
	cfg.HostName = "host00"
	cfg.DomainName = "gd-tools.de"
	cfg.BaselineName = platform.DefaultBaseline

	err = saveServer(cfg)
	if err != nil {
		t.Fatal(err)
	}

	serverDir := filepath.Join(base, "host00.gd-tools.de")
	if _, err := os.Stat(serverDir); err != nil {
		t.Fatalf("server dir missing: %v", err)
	}

	cfgPath := filepath.Join(serverDir, "config.json")
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("config.json missing: %v", err)
	}

	knownHostsPath := filepath.Join(serverDir, "known_hosts")
	data, err := os.ReadFile(knownHostsPath)
	if err != nil {
		t.Fatalf("known_hosts missing: %v", err)
	}

	if string(data) != "example host key\n" {
		t.Fatalf("unexpected known_hosts content: %q", string(data))
	}
}

func TestSaveServerRefusesExistingServer(t *testing.T) {
	base := t.TempDir()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	if err := os.Chdir(base); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{}
	cfg.HostName = "host00"
	cfg.DomainName = "gd-tools.de"
	cfg.BaselineName = "noble-8.3-2.4"

	err = saveServer(cfg)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(base)
	if err != nil {
		t.Fatal(err)
	}

	err = saveServer(cfg)
	if err == nil {
		t.Fatal("expected error for existing server")
	}
}
