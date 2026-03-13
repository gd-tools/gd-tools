package assets

import (
	"testing"
)

func TestCatalogValidation(t *testing.T) {
	cat, err := LoadCatalog()
	if err != nil {
		t.Fatalf("failed to load catalog: %v", err)
	}

	seenProducts := map[string]bool{}

	for _, pr := range cat.Products {
		if pr.Name == "" {
			t.Fatalf("product has no name")
		}

		if seenProducts[pr.Name] {
			t.Fatalf("duplicate product %q", pr.Name)
		}
		seenProducts[pr.Name] = true

		if pr.Default == "" {
			t.Fatalf("product %q has no default release", pr.Name)
		}

		foundDefault := false
		seenVersions := map[string]bool{}

		for _, rel := range pr.Versions {
			if rel.Number == "" {
				t.Fatalf("release %s has no version", pr.Name)
			}

			if seenVersions[rel.Number] {
				t.Fatalf("duplicate release %s/%s", pr.Name, rel.Number)
			}
			seenVersions[rel.Number] = true

			if rel.Download.DownloadURL == "" {
				t.Fatalf("release %s/%s has no download_url", pr.Name, rel.Number)
			}

			if rel.Download.Filename == "" {
				t.Fatalf("release %s/%s has no filename", pr.Name, rel.Number)
			}

			if rel.Download.MD5 == "" &&
				rel.Download.SHA256 == "" &&
				rel.Download.SHA512 == "" {
				t.Fatalf("release %s/%s has no checksum", pr.Name, rel.Number)
			}

			dir := rel.Download.Directory
			if dir == "" {
				dir = pr.Directory
			}

			bin := rel.Download.Binary
			if bin == "" {
				bin = pr.Binary
			}

			if dir != "" && bin != "" {
				t.Fatalf("release %s/%s defines both directory and binary", pr.Name, rel.Number)
			}

			if dir == "" && bin == "" {
				t.Fatalf("release %s/%s defines neither directory nor binary", pr.Name, rel.Number)
			}

			if rel.Number == pr.Default {
				foundDefault = true
			}
		}

		if !foundDefault {
			t.Fatalf("product %q default release %q not found", pr.Name, pr.Default)
		}
	}
}

func TestCatalogDefaultLookup(t *testing.T) {
	cat, err := LoadCatalog()
	if err != nil {
		t.Fatal(err)
	}

	for _, pr := range cat.Products {
		_, rel, err := cat.Get(pr.Name, "")
		if err != nil {
			t.Fatalf("default lookup failed for %s: %v", pr.Name, err)
		}

		if rel.Number != pr.Default {
			t.Fatalf("default mismatch for %s", pr.Name)
		}
	}
}

func TestCatalogReleaseLookup(t *testing.T) {
	cat, err := LoadCatalog()
	if err != nil {
		t.Fatal(err)
	}

	for _, pr := range cat.Products {
		for _, rel := range pr.Versions {
			_, got, err := cat.Get(pr.Name, rel.Number)
			if err != nil {
				t.Fatalf("lookup failed for %s/%s: %v", pr.Name, rel.Number, err)
			}

			if got.Number != rel.Number {
				t.Fatalf("release mismatch for %s/%s", pr.Name, rel.Number)
			}
		}
	}
}
