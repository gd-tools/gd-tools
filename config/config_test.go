package config

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	dir := t.TempDir()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWD)

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("config.json", []byte(`{
  "baseline": "noble-8.3-2.4",
  "host_name": "host00",
  "domain_name": "example.org",
  "ipv4_addr": "192.0.2.10"
}
`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.BaselineName != "noble-8.3-2.4" {
		t.Fatalf("unexpected baseline: %q", cfg.BaselineName)
	}

	if cfg.HostName != "host00" {
		t.Fatalf("unexpected host name: %q", cfg.HostName)
	}

	if cfg.DomainName != "example.org" {
		t.Fatalf("unexpected domain name: %q", cfg.DomainName)
	}

	if cfg.IPv4Addr != "192.0.2.10" {
		t.Fatalf("unexpected ipv4: %q", cfg.IPv4Addr)
	}

	if cfg.Timeout != DefaultTimeout {
		t.Fatalf("unexpected timeout: got %d want %d", cfg.Timeout, DefaultTimeout)
	}
}

func TestConfigLoad(t *testing.T) {
	dir := t.TempDir()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWD)

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("config.json", []byte(`{
  "baseline": "test-base",
  "host_name": "srv",
  "domain_name": "example.org"
}
`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &Config{}
	err = cfg.Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.BaselineName != "test-base" {
		t.Fatalf("unexpected baseline: %q", cfg.BaselineName)
	}

	if cfg.FQDN() != "srv.example.org" {
		t.Fatalf("unexpected fqdn: %q", cfg.FQDN())
	}

	if cfg.Timeout != DefaultTimeout {
		t.Fatalf("unexpected timeout: got %d want %d", cfg.Timeout, DefaultTimeout)
	}
}

func TestConfigSave(t *testing.T) {
	dir := t.TempDir()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWD)

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &Config{}
	cfg.BaselineName = "noble-8.3-2.4"
	cfg.HostName = "host01"
	cfg.DomainName = "railduino.de"
	cfg.IPv4Addr = "198.51.100.20"

	err = cfg.Save()
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile("config.json")
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("config.json is empty")
	}

	cfg2, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	if cfg2.BaselineName != cfg.BaselineName {
		t.Fatalf("unexpected baseline: got %q want %q", cfg2.BaselineName, cfg.BaselineName)
	}

	if cfg2.FQDN() != "host01.railduino.de" {
		t.Fatalf("unexpected fqdn: %q", cfg2.FQDN())
	}

	if cfg2.Timeout != DefaultTimeout {
		t.Fatalf("unexpected timeout: got %d want %d", cfg2.Timeout, DefaultTimeout)
	}
}

func TestConfigLoadNil(t *testing.T) {
	var cfg *Config

	err := cfg.Load()
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestConfigSaveNil(t *testing.T) {
	var cfg *Config

	err := cfg.Save()
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}
