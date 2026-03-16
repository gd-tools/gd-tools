package platform

import (
	"net"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := defaultOptions()
	if opts == nil {
		t.Fatal("expected options")
	}
	if opts.rootDir != "/root" {
		t.Fatalf("unexpected rootDir: %q", opts.rootDir)
	}
	if opts.varDir != "/var" {
		t.Fatalf("unexpected varDir: %q", opts.varDir)
	}
	if opts.etcDir != "/etc" {
		t.Fatalf("unexpected etcDir: %q", opts.etcDir)
	}
	if opts.binDir != "/usr/local/bin" {
		t.Fatalf("unexpected binDir: %q", opts.binDir)
	}
	if opts.runDir != "/run" {
		t.Fatalf("unexpected runDir: %q", opts.runDir)
	}
	if opts.LookupIP == nil {
		t.Fatal("expected LookupIP function")
	}
}

func TestJoinPath(t *testing.T) {
	pf := &Platform{}

	got := pf.joinPath("/var", "gd-tools", "data")
	want := "/var/gd-tools/data"

	if got != want {
		t.Fatalf("joinPath mismatch: got %q want %q", got, want)
	}
}

func TestRootPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			rootDir: "/root",
		},
	}

	got := pf.RootPath("Downloads", "file.txt")
	want := "/root/Downloads/file.txt"

	if got != want {
		t.Fatalf("RootPath mismatch: got %q want %q", got, want)
	}
}

func TestDownloadsPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			rootDir: "/root",
		},
	}

	got := pf.DownloadsPath("archive.tar.gz")
	want := "/root/Downloads/archive.tar.gz"

	if got != want {
		t.Fatalf("DownloadsPath mismatch: got %q want %q", got, want)
	}
}

func TestVarPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			varDir: "/var",
		},
	}

	got := pf.VarPath("lib")
	want := "/var/lib"

	if got != want {
		t.Fatalf("VarPath mismatch: got %q want %q", got, want)
	}
}

func TestToolsPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			varDir: "/var",
		},
	}

	got := pf.ToolsPath("data")
	want := "/var/gd-tools/data"

	if got != want {
		t.Fatalf("ToolsPath mismatch: got %q want %q", got, want)
	}
}

func TestDataPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			varDir: "/var",
		},
	}

	got := pf.DataPath("apache")
	want := "/var/gd-tools/data/apache"

	if got != want {
		t.Fatalf("DataPath mismatch: got %q want %q", got, want)
	}
}

func TestCertsPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			varDir: "/var",
		},
	}

	got := pf.CertsPath("example.org")
	want := "/var/gd-tools/data/certs/example.org"

	if got != want {
		t.Fatalf("CertsPath mismatch: got %q want %q", got, want)
	}
}

func TestToolsApachePath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			varDir: "/var",
		},
	}

	got := pf.ToolsApachePath("acme-challenge")
	want := "/var/gd-tools/data/apache/acme-challenge"

	if got != want {
		t.Fatalf("ToolsApachePath mismatch: got %q want %q", got, want)
	}
}

func TestLogsPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			varDir: "/var",
		},
	}

	got := pf.LogsPath("serve.log")
	want := "/var/gd-tools/logs/serve.log"

	if got != want {
		t.Fatalf("LogsPath mismatch: got %q want %q", got, want)
	}
}

func TestEtcPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			etcDir: "/etc",
		},
	}

	got := pf.EtcPath("apache2", "sites-enabled")
	want := "/etc/apache2/sites-enabled"

	if got != want {
		t.Fatalf("EtcPath mismatch: got %q want %q", got, want)
	}
}

func TestEtcApachePath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			etcDir: "/etc",
		},
	}

	got := pf.EtcApachePath("sites-enabled")
	want := "/etc/apache2/sites-enabled"

	if got != want {
		t.Fatalf("EtcApachePath mismatch: got %q want %q", got, want)
	}
}

func TestEtcPhpPath(t *testing.T) {
	pf := &Platform{
		Baseline: &Baseline{
			PHP: "8.3",
		},
		options: &Options{
			etcDir: "/etc",
		},
	}

	got := pf.EtcPhpPath("cli", "php.ini")
	want := "/etc/php/8.3/cli/php.ini"

	if got != want {
		t.Fatalf("EtcPhpPath mismatch: got %q want %q", got, want)
	}
}

func TestBinPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			binDir: "/usr/local/bin",
		},
	}

	got := pf.BinPath("gdt")
	want := "/usr/local/bin/gdt"

	if got != want {
		t.Fatalf("BinPath mismatch: got %q want %q", got, want)
	}
}

func TestRunPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			runDir: "/run",
		},
	}

	got := pf.RunPath("mysqld", "mysqld.sock")
	want := "/run/mysqld/mysqld.sock"

	if got != want {
		t.Fatalf("RunPath mismatch: got %q want %q", got, want)
	}
}

func TestLookupIPCanBeInjected(t *testing.T) {
	called := false

	opts := &Options{
		LookupIP: func(host string) ([]net.IP, error) {
			called = true
			if host != "example.org" {
				t.Fatalf("unexpected host: %q", host)
			}
			return []net.IP{net.ParseIP("127.0.0.1")}, nil
		},
	}

	ips, err := opts.LookupIP("example.org")
	if err != nil {
		t.Fatalf("LookupIP returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected LookupIP to be called")
	}
	if len(ips) != 1 {
		t.Fatalf("unexpected number of IPs: %d", len(ips))
	}
	if got, want := ips[0].String(), "127.0.0.1"; got != want {
		t.Fatalf("unexpected IP: got %q want %q", got, want)
	}
}

func TestRootPathPanicsWithoutOptions(t *testing.T) {
	pf := &Platform{}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		if got, want := r.(string), "missing options or rootDir"; got != want {
			t.Fatalf("unexpected panic: got %q want %q", got, want)
		}
	}()

	_ = pf.RootPath("x")
}

func TestVarPathPanicsWithoutVarDir(t *testing.T) {
	pf := &Platform{
		options: &Options{},
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		if got, want := r.(string), "missing options or varDir"; got != want {
			t.Fatalf("unexpected panic: got %q want %q", got, want)
		}
	}()

	_ = pf.VarPath("x")
}

func TestEtcPhpPathPanicsWithoutBaseline(t *testing.T) {
	pf := &Platform{
		options: &Options{
			etcDir: "/etc",
		},
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		if got, want := r.(string), "missing Baseline or PHP version"; got != want {
			t.Fatalf("unexpected panic: got %q want %q", got, want)
		}
	}()

	_ = pf.EtcPhpPath("cli")
}
