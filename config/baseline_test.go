package config

import (
	"strings"
	"testing"
)

func TestBaselineValidateOK(t *testing.T) {
	bl := &Baseline{
		Name:     "noble-8.3-2.4",
		Ubuntu:   "24.04 LTS",
		PHP:      "8.3",
		Dovecot:  "2.4",
		Packages: []string{"vim", "curl"},
	}

	if err := bl.Validate(); err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}
}

func TestBaselineValidateNil(t *testing.T) {
	var bl *Baseline

	err := bl.Validate()
	if err == nil {
		t.Fatal("expected error for nil baseline")
	}
}

func TestBaselineValidateMissingName(t *testing.T) {
	bl := &Baseline{
		Ubuntu:   "24.04 LTS",
		PHP:      "8.3",
		Dovecot:  "2.4",
		Packages: []string{"vim"},
	}

	err := bl.Validate()
	if err == nil || !strings.Contains(err.Error(), "without name") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBaselineValidateMissingUbuntu(t *testing.T) {
	bl := &Baseline{
		Name:     "noble-8.3-2.4",
		PHP:      "8.3",
		Dovecot:  "2.4",
		Packages: []string{"vim"},
	}

	err := bl.Validate()
	if err == nil || !strings.Contains(err.Error(), "has no Ubuntu") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBaselineValidateMissingPHP(t *testing.T) {
	bl := &Baseline{
		Name:     "noble-8.3-2.4",
		Ubuntu:   "24.04 LTS",
		Dovecot:  "2.4",
		Packages: []string{"vim"},
	}

	err := bl.Validate()
	if err == nil || !strings.Contains(err.Error(), "has no PHP") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBaselineValidateMissingDovecot(t *testing.T) {
	bl := &Baseline{
		Name:     "noble-8.3-2.4",
		Ubuntu:   "24.04 LTS",
		PHP:      "8.3",
		Packages: []string{"vim"},
	}

	err := bl.Validate()
	if err == nil || !strings.Contains(err.Error(), "has no Dovecot") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBaselineValidateMissingPackages(t *testing.T) {
	bl := &Baseline{
		Name:    "noble-8.3-2.4",
		Ubuntu:  "24.04 LTS",
		PHP:     "8.3",
		Dovecot: "2.4",
	}

	err := bl.Validate()
	if err == nil || !strings.Contains(err.Error(), "has no packages") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadBaselineMissingName(t *testing.T) {
	_, err := LoadBaseline("")
	if err == nil || !strings.Contains(err.Error(), "missing baseline name") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBaselineInfo(t *testing.T) {
	bl := &Baseline{
		Name:     "noble-8.3-2.4",
		Ubuntu:   "24.04 LTS",
		PHP:      "8.3",
		Dovecot:  "2.4",
		Repos:    []string{"docker", "dovecot"},
		Packages: []string{"vim", "curl"},
	}

	lines := bl.Info()
	got := strings.Join(lines, "\n")

	wantParts := []string{
		"Baseline: noble-8.3-2.4",
		"Ubuntu:   24.04 LTS",
		"PHP:      8.3",
		"Dovecot:  2.4",
		"Repositories:",
		"  - docker",
		"  - dovecot",
		"Packages:",
		"  - vim",
		"  - curl",
	}

	for _, part := range wantParts {
		if !strings.Contains(got, part) {
			t.Fatalf("Info() output missing %q\n%s", part, got)
		}
	}
}

func TestBaselineInfoNil(t *testing.T) {
	var bl *Baseline

	lines := bl.Info()
	if lines != nil {
		t.Fatalf("expected nil lines, got %#v", lines)
	}
}
