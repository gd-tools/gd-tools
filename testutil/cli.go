package testutil

import (
	"bytes"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

type CLIResult struct {
	Output string
	Err    error
	Base   string
}

func RunCLI(t *testing.T, app *cli.App, args ...string) CLIResult {
	t.Helper()

	base := t.TempDir()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	if err := os.Chdir(base); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("GD_TOOLS_BASE", base); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	app.Writer = &out

	err = app.Run(append([]string{"gdt"}, args...))

	return CLIResult{
		Output: out.String(),
		Err:    err,
		Base:   base,
	}
}
