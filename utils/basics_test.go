package utils

import (
	"os"
	"testing"
)

func TestLocale(t *testing.T) {
	bsc := Basics{
		Language: "de",
		Region:   "DE",
	}
	if bsc.Locale() != "de_DE" {
		t.Fatalf("unexpected locale: %s", bsc.Locale())
	}
}

func TestHelpers(t *testing.T) {
	bsc := Basics{
		Domain: "example.org",
	}

	if bsc.AdminMail() != "admin@example.org" {
		t.Fatalf("unexpected admin mail")
	}

	if bsc.SupportURL() != "https://support.example.org/" {
		t.Fatalf("unexpected support url")
	}

	if bsc.DMARCDomain() != "_dmarc.example.org" {
		t.Fatalf("unexpected dmarc domain")
	}
}

func TestEnsureAndSaveBasics(t *testing.T) {
	dir := t.TempDir()

	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	bsc, err := EnsureBasics()
	if err != nil {
		t.Fatalf("ensure basics failed: %v", err)
	}

	if bsc.Domain != DefaultDomain {
		t.Fatalf("unexpected default domain")
	}

	bsc.Domain = "example.org"
	bsc.Company = "Example Org"

	if err := bsc.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	bsc2, err := EnsureBasics()
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}

	if bsc2.Domain != "example.org" {
		t.Fatalf("domain mismatch after reload")
	}
}
