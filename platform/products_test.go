package platform

import "testing"

func testProduct() *Product {
	return &Product{
		Name:      "nextcloud",
		Default:   "28.0.1",
		Directory: "nextcloud",
		Binary:    "occ",
		Marker:    "config/config.php",
		Versions: []Release{
			{
				Number: "28.0.1",
				Series: "stable",
				Download: Download{
					Filename: "nextcloud-28.0.1.zip",
				},
			},
			{
				Number: "27.1.5",
				Series: "oldstable",
				Download: Download{
					Filename: "nextcloud-27.1.5.zip",
				},
			},
		},
	}
}

func TestFindProduct(t *testing.T) {

	pf := &Platform{
		Products: []Product{
			{Name: "nextcloud"},
		},
	}

	pr, err := pf.FindProduct("nextcloud")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pr.Name != "nextcloud" {
		t.Fatalf("unexpected product: %s", pr.Name)
	}
}

func TestFindProductNotFound(t *testing.T) {

	pf := &Platform{}

	_, err := pf.FindProduct("missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetReleaseDefault(t *testing.T) {

	pr := testProduct()

	rel, err := pr.GetRelease("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rel.Number != "28.0.1" {
		t.Fatalf("unexpected release: %s", rel.Number)
	}
}

func TestGetReleaseExplicit(t *testing.T) {

	pr := testProduct()

	rel, err := pr.GetRelease("27.1.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rel.Number != "27.1.5" {
		t.Fatalf("unexpected release: %s", rel.Number)
	}
}

func TestGetReleaseNotFound(t *testing.T) {

	pr := testProduct()

	_, err := pr.GetRelease("99.0")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetReleaseInheritsDefaults(t *testing.T) {

	pr := testProduct()

	rel, err := pr.GetRelease("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rel.Download.Directory != "nextcloud" {
		t.Fatalf("unexpected directory: %s", rel.Download.Directory)
	}

	if rel.Download.Binary != "occ" {
		t.Fatalf("unexpected binary: %s", rel.Download.Binary)
	}

	if rel.Download.Marker != "config/config.php" {
		t.Fatalf("unexpected marker: %s", rel.Download.Marker)
	}
}

func TestReleaseCandidates(t *testing.T) {

	pr := testProduct()

	candidates := pr.ReleaseCandidates()

	if len(candidates) != 2 {
		t.Fatalf("unexpected candidate count: %d", len(candidates))
	}

	if !candidates[0].Default {
		t.Fatal("expected first candidate to be default")
	}
}

func TestProductInfo(t *testing.T) {

	pr := testProduct()

	lines := pr.Info()

	if len(lines) == 0 {
		t.Fatal("expected info output")
	}

	if lines[0] != "Product:    nextcloud" {
		t.Fatalf("unexpected first line: %s", lines[0])
	}
}
