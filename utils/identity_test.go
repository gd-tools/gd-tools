package utils

import (
	"os"
	"testing"
)

func TestLocale(t *testing.T) {
	id := Identity{
		Language: "de",
		Region:   "DE",
	}
	if id.Locale() != "de_DE" {
		t.Fatalf("unexpected locale: %s", id.Locale())
	}
}

func TestHelpers(t *testing.T) {
	id := Identity{
		Domain: "example.org",
	}

	if id.AdminMail() != "admin@example.org" {
		t.Fatalf("unexpected admin mail")
	}

	if id.SupportURL() != "https://support.example.org/" {
		t.Fatalf("unexpected support url")
	}

	if id.DMARCDomain() != "_dmarc.example.org" {
		t.Fatalf("unexpected dmarc domain")
	}
}

func TestEnsureAndSaveIdentity(t *testing.T) {
	dir := t.TempDir()

	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	id, err := EnsureIdentity()
	if err != nil {
		t.Fatalf("ensure identity failed: %v", err)
	}

	if id.Domain != DefaultDomain {
		t.Fatalf("unexpected default domain")
	}

	id.Domain = "example.org"
	id.Company = "Example Org"

	if err := id.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	id, err = EnsureIdentity()
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}

	if id.Domain != "example.org" {
		t.Fatalf("domain mismatch after reload")
	}
}
