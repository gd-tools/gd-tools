package identity

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gd-tools/gd-tools/testutil"
	"github.com/urfave/cli/v2"
)

func TestIdentityPrint(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{Command},
	}

	res := testutil.RunCLI(t, app, "identity")

	if res.Err != nil {
		t.Fatalf("command failed: %v", res.Err)
	}

	if !strings.Contains(res.Output, "Company") {
		t.Fatalf("expected table output")
	}
}

func TestIdentityUpdate(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{Command},
	}

	res := testutil.RunCLI(t, app,
		"identity",
		"--company", "Example Ltd",
		"--domain", "example.org",
	)

	if res.Err != nil {
		t.Fatalf("command failed: %v", res.Err)
	}

	data, err := os.ReadFile(filepath.Join(res.Base, "identity.json"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "Example Ltd") {
		t.Fatalf("company not saved")
	}
}
