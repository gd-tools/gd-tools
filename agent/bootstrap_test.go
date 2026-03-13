package agent

import (
	"os"
	"testing"

	"github.com/gd-tools/gd-tools/assets"
)

func TestBootstrapTestFirstRun(t *testing.T) {
	tmp := t.TempDir()

	restore := assets.SetTestRootDir(tmp)
	defer restore()

	req := &Request{}

	if !BootstrapTest(req) {
		t.Fatal("expected BootstrapTest to trigger on first run")
	}
}

func TestBootstrapTestNoChanges(t *testing.T) {
	tmp := t.TempDir()

	restore := assets.SetTestRootDir(tmp)
	defer restore()

	marker := assets.GetRootDir(FirstRunMarker)

	if err := os.WriteFile(marker, []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}

	req := &Request{}

	if BootstrapTest(req) {
		t.Fatal("BootstrapTest should not trigger")
	}
}

func TestBootstrapTestWithChanges(t *testing.T) {
	tmp := t.TempDir()

	restore := assets.SetTestRootDir(tmp)
	defer restore()

	marker := assets.GetRootDir(FirstRunMarker)

	if err := os.WriteFile(marker, []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}

	req := &Request{
		FQDN: "example.com",
	}

	if !BootstrapTest(req) {
		t.Fatal("expected BootstrapTest when request has changes")
	}
}

func TestBootstrapHandlerNil(t *testing.T) {
	if err := BootstrapHandler(nil, nil); err != nil {
		t.Fatal(err)
	}
}
