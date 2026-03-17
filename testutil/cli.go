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

	oldBase := os.Getenv("GD_TOOLS_BASE")

	t.Cleanup(func() {
		_ = os.Chdir(oldwd)
		if oldBase == "" {
			_ = os.Unsetenv("GD_TOOLS_BASE")
		} else {
			_ = os.Setenv("GD_TOOLS_BASE", oldBase)
		}
	})

	if err := os.Chdir(base); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("GD_TOOLS_BASE", base); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	app.Writer = &out
	app.ErrWriter = &errOut

	err = app.Run(append([]string{"gdt"}, args...))

	output := out.String()
	if errOut.Len() > 0 {
		output += errOut.String()
	}

	return CLIResult{
		Output: output,
		Err:    err,
		Base:   base,
	}
}
