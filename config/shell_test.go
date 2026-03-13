package config

import (
	"os"
	"path/filepath"
	"testing"
)

func fakeCmd(t *testing.T, name, script string) string {
	dir := t.TempDir()
	path := filepath.Join(dir, name)

	content := "#!/bin/sh\n" + script + "\n"
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", dir+":"+oldPath)

	t.Cleanup(func() {
		os.Setenv("PATH", oldPath)
	})

	return path
}

func testShellConfig() *Config {
	return &Config{
		HostName:   "test",
		DomainName: "example.com",
	}
}

func TestCheckRemoteSuccess(t *testing.T) {
	fakeCmd(t, "ssh", "exit 0")

	cfg := testShellConfig()

	if !cfg.CheckRemote("true") {
		t.Fatal("expected CheckRemote to succeed")
	}
}

func TestCheckRemoteFailure(t *testing.T) {
	fakeCmd(t, "ssh", "exit 1")

	cfg := testShellConfig()

	if cfg.CheckRemote("true") {
		t.Fatal("expected CheckRemote to fail")
	}
}

func TestRemoteCmdSuccess(t *testing.T) {
	fakeCmd(t, "ssh", "exit 0")

	cfg := testShellConfig()

	if err := cfg.RemoteCmd("echo ok"); err != nil {
		t.Fatalf("RemoteCmd failed: %v", err)
	}
}

func TestRemoteCmdFailure(t *testing.T) {
	fakeCmd(t, "ssh", "exit 1")

	cfg := testShellConfig()

	if err := cfg.RemoteCmd("echo fail"); err == nil {
		t.Fatal("expected RemoteCmd to fail")
	}
}

func TestRemoteScript(t *testing.T) {
	fakeCmd(t, "ssh", "exit 0")

	cfg := testShellConfig()

	err := cfg.RemoteScript([]string{
		"echo one",
		"echo two",
	})

	if err != nil {
		t.Fatalf("RemoteScript failed: %v", err)
	}
}

func TestRemoteScriptEmpty(t *testing.T) {
	cfg := testShellConfig()

	if err := cfg.RemoteScript([]string{}); err == nil {
		t.Fatal("expected error for empty command list")
	}
}
