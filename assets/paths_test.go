package assets

import (
	"path/filepath"
	"testing"
)

func withDirs(d directories) func() {
	old := defaultDirs
	defaultDirs = d
	return func() {
		defaultDirs = old
	}
}

func TestPaths(t *testing.T) {
	defer withDirs(directories{
		rootDir:      "/tmp/root",
		varDir:       "/tmp/var",
		toolsDir:     "/tmp/tools",
		etcDir:       "/tmp/etc",
		binDir:       "/tmp/bin",
		runDir:       "/tmp/run",
		downloadsDir: "/tmp/downloads",
	})()

	tests := []struct {
		name string
		fn   func(...string) string
		base string
	}{
		{"root", GetRootDir, "/tmp/root"},
		{"var", GetVarDir, "/tmp/var"},
		{"tools", GetToolsDir, "/tmp/tools"},
		{"etc", GetEtcDir, "/tmp/etc"},
		{"bin", GetBinDir, "/tmp/bin"},
		{"run", GetRunDir, "/tmp/run"},
		{"downloads", GetDownloadsDir, "/tmp/downloads"},
	}

	for _, tt := range tests {
		got := tt.fn()
		if got != tt.base {
			t.Fatalf("%s: expected %s got %s", tt.name, tt.base, got)
		}

		want := filepath.Join(tt.base, "a", "b")
		got = tt.fn("a", "b")

		if got != want {
			t.Fatalf("%s: expected %s got %s", tt.name, want, got)
		}
	}
}

func TestApachePaths(t *testing.T) {
	defer withDirs(directories{
		rootDir:      "/tmp/root",
		varDir:       "/tmp/var",
		toolsDir:     "/tmp/tools",
		etcDir:       "/tmp/etc",
		binDir:       "/tmp/bin",
		runDir:       "/tmp/run",
		downloadsDir: "/tmp/downloads",
	})()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			"apache tools base",
			GetApacheToolsDir(),
			filepath.Join("/tmp/tools", "data", "apache"),
		},
		{
			"apache tools join",
			GetApacheToolsDir("conf"),
			filepath.Join("/tmp/tools", "data", "apache", "conf"),
		},
		{
			"apache etc base",
			GetApacheEtcDir(),
			filepath.Join("/tmp/etc", "apache2"),
		},
		{
			"apache etc join",
			GetApacheEtcDir("sites-enabled"),
			filepath.Join("/tmp/etc", "apache2", "sites-enabled"),
		},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Fatalf("%s: expected %s got %s", tt.name, tt.want, tt.got)
		}
	}
}
