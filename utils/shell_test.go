package utils

import (
	"strings"
	"testing"
)

func TestRunCommandSuccess(t *testing.T) {
	out, err := RunCommand("echo", "hello")
	if err != nil {
		t.Fatalf("RunCommand failed: %v", err)
	}

	if strings.TrimSpace(string(out)) != "hello" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestRunCommandFailure(t *testing.T) {
	_, err := RunCommand("false")
	if err == nil {
		t.Fatalf("expected error from failing command")
	}
}

func TestRunShellSuccess(t *testing.T) {
	out, err := RunShell([]string{
		"echo hello",
		"echo world",
	})
	if err != nil {
		t.Fatalf("RunShell failed: %v", err)
	}

	output := string(out)

	if !strings.Contains(output, "hello") {
		t.Fatalf("missing hello in output: %s", output)
	}

	if !strings.Contains(output, "world") {
		t.Fatalf("missing world in output: %s", output)
	}
}

func TestRunShellFailure(t *testing.T) {
	_, err := RunShell([]string{
		"echo ok",
		"false",
	})

	if err == nil {
		t.Fatalf("expected error from failing shell script")
	}
}

func TestStartServiceValidation(t *testing.T) {
	_, err := StartService("")
	if err == nil {
		t.Fatalf("expected error for empty service name")
	}
}
