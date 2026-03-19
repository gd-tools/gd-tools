package config

import (
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	host string
}

func (c *testConfig) FQDN() string {
	return c.host
}

func TestEnsureCA(t *testing.T) {
	dir := t.TempDir()

	old, _ := os.Getwd()
	defer os.Chdir(old)

	os.Chdir(dir)

	cfg := &Config{}
	cfg.HostName = "example"
	cfg.DomainName = "test"

	if err := cfg.EnsureCA(); err != nil {
		t.Fatal(err)
	}

	files := []string{
		SerialName,
		CaKeyName,
		CaCrtName,
		ClientKeyName,
		ClientCsrName,
		ClientCrtName,
		ServerKeyName,
		ServerConfigName,
		ServerCsrName,
		ServerCrtName,
	}

	for _, f := range files {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Fatalf("missing file %s", f)
		}
	}
}

func TestEnsureCAIdempotent(t *testing.T) {
	dir := t.TempDir()

	old, _ := os.Getwd()
	defer os.Chdir(old)

	os.Chdir(dir)

	cfg := &Config{}
	cfg.HostName = "example"
	cfg.DomainName = "test"

	if err := cfg.EnsureCA(); err != nil {
		t.Fatal(err)
	}

	info1, _ := os.Stat(CaCrtName)

	if err := cfg.EnsureCA(); err != nil {
		t.Fatal(err)
	}

	info2, _ := os.Stat(CaCrtName)

	if info1.ModTime() != info2.ModTime() {
		t.Fatal("CA certificate should not change")
	}
}
