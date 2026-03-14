package platform

import "testing"

func TestDefaultBaselines(t *testing.T) {
	baselines := DefaultBaselines()
	if len(baselines) == 0 {
		t.Fatal("expected embedded baselines")
	}
}

func TestGetBaseline(t *testing.T) {
	pf := &Platform{
		Baselines: []Baseline{
			{
				Name:    "noble-8.3-2.4",
				Ubuntu:  "24.04 LTS",
				PHP:     "8.3",
				Dovecot: "2.4",
			},
		},
		Products: []Product{
			{
				Name:    "nextcloud",
				Default: "31.0.0",
				Versions: []Release{
					{
						Number: "31.0.0",
						Download: Download{
							Filename: "nextcloud.tar.bz2",
						},
					},
				},
			},
		},
		Paths: DefaultPaths(),
	}

	bl, err := pf.GetBaseline("noble-8.3-2.4")
	if err != nil {
		t.Fatalf("GetBaseline failed: %v", err)
	}
	if bl.Name != "noble-8.3-2.4" {
		t.Fatalf("unexpected baseline name: %q", bl.Name)
	}
	if bl.Ubuntu != "24.04 LTS" {
		t.Fatalf("unexpected ubuntu version: %q", bl.Ubuntu)
	}
	if bl.PHP != "8.3" {
		t.Fatalf("unexpected php version: %q", bl.PHP)
	}
	if bl.Dovecot != "2.4" {
		t.Fatalf("unexpected dovecot version: %q", bl.Dovecot)
	}
}

func TestGetBaselineNotFound(t *testing.T) {
	pf := &Platform{
		Baselines: []Baseline{
			{Name: "noble-8.3-2.4"},
		},
		Products: []Product{
			{
				Name:    "nextcloud",
				Default: "31.0.0",
				Versions: []Release{
					{
						Number: "31.0.0",
						Download: Download{
							Filename: "nextcloud.tar.bz2",
						},
					},
				},
			},
		},
		Paths: DefaultPaths(),
	}

	_, err := pf.GetBaseline("missing")
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
}

func TestBaselineInfo(t *testing.T) {
	bl := &Baseline{
		Name:     "noble-8.3-2.4",
		Ubuntu:   "24.04 LTS",
		PHP:      "8.3",
		Dovecot:  "2.4",
		Repos:    []string{"docker", "dovecot"},
		Packages: []string{"curl", "vim"},
	}

	lines := bl.Info()
	if len(lines) == 0 {
		t.Fatal("expected info lines")
	}

	foundName := false
	foundUbuntu := false
	foundPHP := false
	foundDovecot := false
	foundRepo := false
	foundPackage := false

	for _, line := range lines {
		switch line {
		case "Baseline: noble-8.3-2.4":
			foundName = true
		case "Ubuntu:   24.04 LTS":
			foundUbuntu = true
		case "PHP:      8.3":
			foundPHP = true
		case "Dovecot:  2.4":
			foundDovecot = true
		case "  - docker":
			foundRepo = true
		case "  - curl":
			foundPackage = true
		}
	}

	if !foundName {
		t.Fatal("missing baseline name in info")
	}
	if !foundUbuntu {
		t.Fatal("missing ubuntu in info")
	}
	if !foundPHP {
		t.Fatal("missing php in info")
	}
	if !foundDovecot {
		t.Fatal("missing dovecot in info")
	}
	if !foundRepo {
		t.Fatal("missing repo in info")
	}
	if !foundPackage {
		t.Fatal("missing package in info")
	}
}
