package platform

import "testing"

func TestLoadPlatformWithPaths(t *testing.T) {
	paths := []Path{
		{Name: PathRoot, Value: "/troot"},
		{Name: PathVar, Value: "/tvar"},
		{Name: PathTools, Value: "/tvar/gd-tools"},
		{Name: PathEtc, Value: "/tetc"},
		{Name: PathBin, Value: "/tbin"},
		{Name: PathRun, Value: "/trun"},
		{Name: PathDownloads, Value: "/tdownloads"},
	}

	pf, err := LoadPlatformWithPaths(paths)
	if err != nil {
		t.Fatalf("LoadPlatformWithPaths failed: %v", err)
	}

	if pf == nil {
		t.Fatal("expected platform")
	}
	if len(pf.Paths) != len(paths) {
		t.Fatalf("unexpected path count: got %d want %d", len(pf.Paths), len(paths))
	}
	if pf.Paths[0].Value != "/troot" {
		t.Fatalf("unexpected first path value: %q", pf.Paths[0].Value)
	}
}

func TestLoadPlatformWithPathsClonesInput(t *testing.T) {
	paths := []Path{
		{Name: PathRoot, Value: "/troot"},
		{Name: PathVar, Value: "/tvar"},
		{Name: PathTools, Value: "/tvar/gd-tools"},
		{Name: PathEtc, Value: "/tetc"},
		{Name: PathBin, Value: "/tbin"},
		{Name: PathRun, Value: "/trun"},
		{Name: PathDownloads, Value: "/tdownloads"},
	}

	pf, err := LoadPlatformWithPaths(paths)
	if err != nil {
		t.Fatalf("LoadPlatformWithPaths failed: %v", err)
	}

	paths[0].Value = "/changed"
	if pf.Paths[0].Value != "/troot" {
		t.Fatal("platform paths must not be affected by caller changes")
	}
}

func TestValidateOK(t *testing.T) {
	pf := &Platform{
		Baselines: []Baseline{
			{Name: "noble-8.3-2.4"},
		},
		Products: []Product{
			{Name: "nextcloud", Default: "31.0.0"},
		},
		Paths: DefaultPaths(),
	}

	if err := pf.Validate(); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
}

func TestValidateNil(t *testing.T) {
	var pf *Platform

	if err := pf.Validate(); err == nil {
		t.Fatal("expected error for nil platform")
	}
}

func TestValidateMissingBaselines(t *testing.T) {
	pf := &Platform{
		Products: []Product{
			{Name: "nextcloud", Default: "31.0.0"},
		},
		Paths: DefaultPaths(),
	}

	if err := pf.Validate(); err == nil {
		t.Fatal("expected error for missing baselines")
	}
}

func TestValidateMissingProducts(t *testing.T) {
	pf := &Platform{
		Baselines: []Baseline{
			{Name: "noble-8.3-2.4"},
		},
		Paths: DefaultPaths(),
	}

	if err := pf.Validate(); err == nil {
		t.Fatal("expected error for missing products")
	}
}

func TestValidateMissingPaths(t *testing.T) {
	pf := &Platform{
		Baselines: []Baseline{
			{Name: "noble-8.3-2.4"},
		},
		Products: []Product{
			{Name: "nextcloud", Default: "31.0.0"},
		},
	}

	if err := pf.Validate(); err == nil {
		t.Fatal("expected error for missing paths")
	}
}

func TestMustPath(t *testing.T) {
	pf, err := LoadPlatformWithPaths([]Path{
		{Name: PathRoot, Value: "/troot"},
		{Name: PathVar, Value: "/tvar"},
		{Name: PathTools, Value: "/tvar/gd-tools"},
		{Name: PathEtc, Value: "/tetc"},
		{Name: PathBin, Value: "/tbin"},
		{Name: PathRun, Value: "/trun"},
		{Name: PathDownloads, Value: "/tdownloads"},
	})
	if err != nil {
		t.Fatalf("LoadPlatformWithPaths failed: %v", err)
	}

	if got := pf.MustPath(PathTools); got != "/tvar/gd-tools" {
		t.Fatalf("unexpected tools path: %q", got)
	}
}
