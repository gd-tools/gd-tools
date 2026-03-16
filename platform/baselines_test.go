package platform

import "testing"

func TestBaselineInfoMinimal(t *testing.T) {
	bl := &Baseline{
		Name: "noble-8.3-2.4",
	}

	lines := bl.Info()

	if len(lines) != 1 {
		t.Fatalf("unexpected number of lines: got %d want 1", len(lines))
	}
	if got, want := lines[0], "Baseline: noble-8.3-2.4"; got != want {
		t.Fatalf("unexpected line: got %q want %q", got, want)
	}
}

func TestBaselineInfoFull(t *testing.T) {
	bl := &Baseline{
		Name:    "noble-8.3-2.4",
		Ubuntu:  "24.04",
		PHP:     "8.3",
		Dovecot: "2.4",
		Repos: []string{
			"docker",
			"dovecot",
		},
		Packages: []string{
			"apache2",
			"php8.3",
		},
	}

	lines := bl.Info()

	want := []string{
		"Baseline: noble-8.3-2.4",
		"Ubuntu:   24.04",
		"PHP:      8.3",
		"Dovecot:  2.4",
		"Repositories:",
		"  - docker",
		"  - dovecot",
		"Packages:",
		"  - apache2",
		"  - php8.3",
	}

	if len(lines) != len(want) {
		t.Fatalf("unexpected number of lines: got %d want %d", len(lines), len(want))
	}

	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d mismatch: got %q want %q", i, lines[i], want[i])
		}
	}
}
