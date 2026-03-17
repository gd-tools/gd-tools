package testutil

import (
	"os"
	"strings"
	"testing"
)

func RequireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func RequireError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error")
	}
}

func RequireContains(t *testing.T, text, want string) {
	t.Helper()

	if !strings.Contains(text, want) {
		t.Fatalf("output does not contain %q\n\noutput:\n%s", want, text)
	}
}

func RequireNotContains(t *testing.T, text, unwanted string) {
	t.Helper()

	if strings.Contains(text, unwanted) {
		t.Fatalf("output unexpectedly contains %q\n\noutput:\n%s", unwanted, text)
	}
}

func RequireFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %q to exist: %v", path, err)
	}
}

func RequireFileNotExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected file %q not to exist", path)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat %q failed: %v", path, err)
	}
}
