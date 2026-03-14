package platform

import "testing"

func TestFindProduct(t *testing.T) {
	pf := &Platform{
		Baselines: []Baseline{{Name: "base"}},
		Products: []Product{
			{Name: "nextcloud", Default: "31.0.0"},
		},
		Paths: DefaultPaths(),
	}

	pr, err := pf.FindProduct("nextcloud")
	if err != nil {
		t.Fatalf("FindProduct failed: %v", err)
	}
	if pr.Name != "nextcloud" {
		t.Fatalf("unexpected product: %q", pr.Name)
	}
}

func TestGetReleaseDefaultInheritance(t *testing.T) {
	pr := &Product{
		Name:      "nextcloud",
		Default:   "31.0.0",
		Directory: "nextcloud",
		Binary:    "occ",
		Marker:    "README",
		Versions: []Release{
			{
				Number: "31.0.0",
				Download: Download{
					Filename: "nextcloud-31.0.0.tar.bz2",
				},
			},
		},
	}

	rel, err := pr.GetRelease("")
	if err != nil {
		t.Fatalf("GetRelease failed: %v", err)
	}

	if rel.Number != "31.0.0" {
		t.Fatalf("unexpected release: %q", rel.Number)
	}
	if rel.Download.Directory != "nextcloud" {
		t.Fatalf("unexpected directory: %q", rel.Download.Directory)
	}
	if rel.Download.Binary != "occ" {
		t.Fatalf("unexpected binary: %q", rel.Download.Binary)
	}
	if rel.Download.Marker != "README" {
		t.Fatalf("unexpected marker: %q", rel.Download.Marker)
	}
}

func TestGetProduct(t *testing.T) {
	pf := &Platform{
		Baselines: []Baseline{{Name: "base"}},
		Products: []Product{
			{
				Name:    "nextcloud",
				Default: "31.0.0",
				Versions: []Release{
					{Number: "31.0.0", Download: Download{Filename: "a.tar.bz2"}},
				},
			},
		},
		Paths: DefaultPaths(),
	}

	pr, rel, err := pf.GetProduct("nextcloud", "")
	if err != nil {
		t.Fatalf("GetProduct failed: %v", err)
	}
	if pr.Name != "nextcloud" {
		t.Fatalf("unexpected product: %q", pr.Name)
	}
	if rel.Number != "31.0.0" {
		t.Fatalf("unexpected release: %q", rel.Number)
	}
}

func TestReleaseCandidates(t *testing.T) {
	pr := &Product{
		Name:    "nextcloud",
		Default: "31.0.0",
		Versions: []Release{
			{Number: "30.0.0"},
			{Number: "31.0.0", Series: "31"},
		},
	}

	got := pr.ReleaseCandidates()
	if len(got) != 2 {
		t.Fatalf("unexpected candidate count: %d", len(got))
	}
	if got[1].Value != "31.0.0" {
		t.Fatalf("unexpected value: %q", got[1].Value)
	}
	if !got[1].Default {
		t.Fatal("expected default candidate")
	}
	if got[1].Description == "" {
		t.Fatal("expected description")
	}
}

func TestDefaultRelease(t *testing.T) {
	pr := &Product{
		Name:    "nextcloud",
		Default: "31.0.0",
		Versions: []Release{
			{Number: "31.0.0"},
		},
	}

	rel, err := pr.DefaultRelease()
	if err != nil {
		t.Fatalf("DefaultRelease failed: %v", err)
	}
	if rel.Number != "31.0.0" {
		t.Fatalf("unexpected release: %q", rel.Number)
	}
}

func TestIsDefaultRelease(t *testing.T) {
	pr := &Product{Default: "31.0.0"}

	if !pr.IsDefaultRelease("31.0.0") {
		t.Fatal("expected default release")
	}
	if pr.IsDefaultRelease("30.0.0") {
		t.Fatal("did not expect default release")
	}
}
