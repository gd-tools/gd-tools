package platform

import "testing"

func TestValidateOK(t *testing.T) {
	pf := &Platform{
		baselines: []Baseline{
			{Name: "noble-8.3-2.4", PHP: "8.3"},
		},
		Baseline: &Baseline{
			Name: "noble-8.3-2.4",
			PHP:  "8.3",
		},
		options: &Options{
			rootDir: "/root",
			varDir:  "/var",
			etcDir:  "/etc",
			binDir:  "/usr/local/bin",
			runDir:  "/run",
		},
		Products: []Product{
			{},
		},
	}

	if err := pf.Validate(); err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}
}

func TestValidateNoBaselines(t *testing.T) {
	pf := &Platform{
		Baseline: &Baseline{Name: "noble-8.3-2.4"},
		options:  &Options{},
		Products: []Product{{}},
	}

	err := pf.Validate()
	if err == nil {
		t.Fatal("expected error for missing baselines")
	}
	if got, want := err.Error(), "platform has no baselines"; got != want {
		t.Fatalf("unexpected error: got %q want %q", got, want)
	}
}

func TestValidateNoBaselinePointer(t *testing.T) {
	pf := &Platform{
		baselines: []Baseline{{Name: "noble-8.3-2.4"}},
		options:   &Options{},
		Products:  []Product{{}},
	}

	err := pf.Validate()
	if err == nil {
		t.Fatal("expected error for missing baseline pointer")
	}
	if got, want := err.Error(), "platform has no baseline pointer"; got != want {
		t.Fatalf("unexpected error: got %q want %q", got, want)
	}
}

func TestValidateNoOptions(t *testing.T) {
	pf := &Platform{
		baselines: []Baseline{{Name: "noble-8.3-2.4"}},
		Baseline:  &Baseline{Name: "noble-8.3-2.4"},
		Products:  []Product{{}},
	}

	err := pf.Validate()
	if err == nil {
		t.Fatal("expected error for missing options")
	}
	if got, want := err.Error(), "platform has no options"; got != want {
		t.Fatalf("unexpected error: got %q want %q", got, want)
	}
}

func TestValidateNoProducts(t *testing.T) {
	pf := &Platform{
		baselines: []Baseline{{Name: "noble-8.3-2.4"}},
		Baseline:  &Baseline{Name: "noble-8.3-2.4"},
		options:   &Options{},
	}

	err := pf.Validate()
	if err == nil {
		t.Fatal("expected error for missing products")
	}
	if got, want := err.Error(), "platform has no products"; got != want {
		t.Fatalf("unexpected error: got %q want %q", got, want)
	}
}
