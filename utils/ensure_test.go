package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureBaseDir(t *testing.T) {
	base := t.TempDir()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	os.Setenv("GD_TOOLS_BASE", base)

	if err := os.Chdir(base); err != nil {
		t.Fatal(err)
	}

	if err := EnsureBaseDir(); err != nil {
		t.Fatalf("expected base dir to be valid: %v", err)
	}

	sub := filepath.Join(base, "host1")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(sub); err != nil {
		t.Fatal(err)
	}

	if err := EnsureBaseDir(); err != nil {
		t.Fatalf("expected subdir to be valid: %v", err)
	}
}

func TestEnsureBaseDirOutside(t *testing.T) {
	base := t.TempDir()
	outside := t.TempDir()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	os.Setenv("GD_TOOLS_BASE", base)

	if err := os.Chdir(outside); err != nil {
		t.Fatal(err)
	}

	if err := EnsureBaseDir(); err == nil {
		t.Fatalf("expected error outside base dir")
	}
}

func TestEnsureHostDir(t *testing.T) {
	dir := t.TempDir()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("config.json", []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := EnsureHostDir(); err != nil {
		t.Fatalf("expected host dir to be valid: %v", err)
	}
}

func TestEnsureBaseOrHostDir(t *testing.T) {
	base := t.TempDir()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	os.Setenv("GD_TOOLS_BASE", base)

	host := filepath.Join(base, "host1")
	if err := os.Mkdir(host, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(host, "config.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(host); err != nil {
		t.Fatal(err)
	}

	if err := EnsureBaseOrHostDir(); err != nil {
		t.Fatalf("expected host dir to be valid: %v", err)
	}
}
