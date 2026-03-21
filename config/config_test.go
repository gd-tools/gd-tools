package config

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestConfigSave_UsesInjectedSaveJSON(t *testing.T) {
	cfg := &Config{}

	called := false
	cfg.saveJSON = func(name string, v any) error {
		called = true
		if name == "" {
			t.Fatal("expected filename")
		}
		if v != &cfg.Server {
			t.Fatalf("unexpected value: got %T want *Server", v)
		}
		return nil
	}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected saveJSON to be called")
	}
}

func TestConfigSave_PropagatesError(t *testing.T) {
	cfg := &Config{}

	wantErr := errTest("boom")
	cfg.saveJSON = func(name string, v any) error {
		return wantErr
	}

	err := cfg.Save()
	if err == nil {
		t.Fatal("expected error")
	}
	if err != wantErr {
		t.Fatalf("unexpected error: got %v want %v", err, wantErr)
	}
}

func TestReadConfig_SetsFlags(t *testing.T) {
	app := &cli.App{}

	set := flagSet(t, map[string]string{
		"verbose":  "true",
		"force":    "true",
		"delete":   "true",
		"skip-dns": "true",
		"skip-mx":  "true",
		"port":     "1234",
	})

	ctx := cli.NewContext(app, set, nil)

	cfg, err := ReadConfig(ctx)
	if err != nil {
		t.Fatalf("ReadConfig returned error: %v", err)
	}

	if !cfg.Verbose {
		t.Fatal("expected Verbose=true")
	}
	if !cfg.Force {
		t.Fatal("expected Force=true")
	}
	if !cfg.Delete {
		t.Fatal("expected Delete=true")
	}
	if !cfg.SkipDNS {
		t.Fatal("expected SkipDNS=true")
	}
	if !cfg.SkipMX {
		t.Fatal("expected SkipMX=true")
	}
	if cfg.Port != "1234" {
		t.Fatalf("unexpected Port: got %q want %q", cfg.Port, "1234")
	}
}

func TestReadConfig_NilContext(t *testing.T) {
	cfg, err := ReadConfig(nil)
	if err != nil {
		t.Fatalf("ReadConfig returned error: %v", err)
	}

	// defaults should remain false/empty
	if cfg.Verbose || cfg.Force || cfg.Delete || cfg.SkipDNS || cfg.SkipMX {
		t.Fatal("expected all flags to be false")
	}
	if cfg.Port != "" {
		t.Fatalf("expected empty port, got %q", cfg.Port)
	}
}
