package platform

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

// NewTestPlatform returns a deterministic platform for tests.
func NewTestPlatform(t testing.TB) *Platform {
	t.Helper()
	return NewTestPlatformWithBaseline(t, DefaultBaseline)
}

// NewTestPlatformWithBaseline returns a deterministic platform for tests
// using the selected baseline.
func NewTestPlatformWithBaseline(t testing.TB, baseline string) *Platform {
	t.Helper()

	base := t.TempDir()

	opts := &Options{
		rootDir: filepath.Join(base, "root"),
		varDir:  filepath.Join(base, "var"),
		etcDir:  filepath.Join(base, "etc"),
		binDir:  filepath.Join(base, "bin"),
		runDir:  filepath.Join(base, "run"),

		LookupIP: func(host string) ([]net.IP, error) {
			return []net.IP{
				net.ParseIP("192.0.2.10"),
				net.ParseIP("2001:db8::10"),
			}, nil
		},
	}

	pf, err := LoadPlatform(baseline, opts)
	if err != nil {
		t.Fatalf("LoadPlatform(%q) failed: %v", baseline, err)
	}

	dirs := []string{
		pf.RootPath(),
		pf.VarPath(),
		pf.EtcPath(),
		pf.BinPath(),
		pf.RunPath(),
		pf.DownloadsPath(),
		pf.ToolsPath(),
		pf.DataPath(),
		pf.CertsPath(),
		pf.ToolsApachePath(),
		pf.LogsPath(),
		pf.EtcApachePath(),
		pf.EtcPhpPath(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) failed: %v", dir, err)
		}
	}

	return pf
}
