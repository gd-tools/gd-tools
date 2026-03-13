package hello

import (
	"testing"
)

func TestCommandDefinition(t *testing.T) {
	if Command.Name != "hello" {
		t.Fatalf("unexpected command name: %s", Command.Name)
	}

	if Command.Action == nil {
		t.Fatal("command action must not be nil")
	}

	if len(Command.Flags) == 0 {
		t.Fatal("expected flags to be defined")
	}
}
