package platform

import "testing"

func TestDefaultPathsNotEmpty(t *testing.T) {
	paths := DefaultPaths()
	if len(paths) == 0 {
		t.Fatal("expected default paths")
	}
}

func TestClonePaths(t *testing.T) {
	src := DefaultPaths()
	dst := ClonePaths(src)

	if len(dst) != len(src) {
		t.Fatalf("clone length mismatch: got %d want %d", len(dst), len(src))
	}

	dst[0].Value = "/tmp/changed"
	if src[0].Value == dst[0].Value {
		t.Fatal("clone must not modify source slice")
	}
}

func TestPlatformDirMethods(t *testing.T) {
	pf, err := LoadPlatformWithPaths([]Path{
		{Name: PathRoot, Value: "/xroot"},
		{Name: PathVar, Value: "/xvar"},
		{Name: PathTools, Value: "/xvar/gd-tools"},
		{Name: PathEtc, Value: "/xetc"},
		{Name: PathBin, Value: "/xbin"},
		{Name: PathRun, Value: "/xrun"},
		{Name: PathDownloads, Value: "/xdownloads"},
	})
	if err != nil {
		t.Fatalf("LoadPlatformWithPaths failed: %v", err)
	}

	if got := pf.RootDir("abc"); got != "/xroot/abc" {
		t.Fatalf("RootDir mismatch: got %q", got)
	}
	if got := pf.VarDir("lib"); got != "/xvar/lib" {
		t.Fatalf("VarDir mismatch: got %q", got)
	}
	if got := pf.ToolsDir("data"); got != "/xvar/gd-tools/data" {
		t.Fatalf("ToolsDir mismatch: got %q", got)
	}
	if got := pf.DataDir("apache"); got != "/xvar/gd-tools/data/apache" {
		t.Fatalf("DataDir mismatch: got %q", got)
	}
	if got := pf.CertsDir("site.pem"); got != "/xvar/gd-tools/data/certs/site.pem" {
		t.Fatalf("CertsDir mismatch: got %q", got)
	}
	if got := pf.LogsDir("app.log"); got != "/xvar/gd-tools/logs/app.log" {
		t.Fatalf("LogsDir mismatch: got %q", got)
	}
	if got := pf.EtcDir("apache2"); got != "/xetc/apache2" {
		t.Fatalf("EtcDir mismatch: got %q", got)
	}
	if got := pf.BinDir("nextcloud"); got != "/xbin/nextcloud" {
		t.Fatalf("BinDir mismatch: got %q", got)
	}
	if got := pf.RunDir("test.pid"); got != "/xrun/test.pid" {
		t.Fatalf("RunDir mismatch: got %q", got)
	}
	if got := pf.DownloadsDir("a.tar.gz"); got != "/xdownloads/a.tar.gz" {
		t.Fatalf("DownloadsDir mismatch: got %q", got)
	}
	if got := pf.ApacheToolsDir("conf"); got != "/xvar/gd-tools/data/apache/conf" {
		t.Fatalf("ApacheToolsDir mismatch: got %q", got)
	}
	if got := pf.ApacheEtcDir("mods-enabled"); got != "/xetc/apache2/mods-enabled" {
		t.Fatalf("ApacheEtcDir mismatch: got %q", got)
	}
}
