package agent

import (
	"strings"
	"testing"
)

func TestRunShellSuccess(t *testing.T) {
	out, err := RunShell([]string{
		"echo hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "hello") {
		t.Fatal("unexpected output")
	}
}

func TestRunShellFailure(t *testing.T) {
	_, err := RunShell([]string{
		"false",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "shell failed") {
		t.Fatal("unexpected error")
	}
}

func TestRunShellEmpty(t *testing.T) {
	_, err := RunShell([]string{})
	if err == nil {
		t.Fatal("expected error for empty commands")
	}
}

func TestRunCommandSuccess(t *testing.T) {
	out, err := RunCommand("echo", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "hello") {
		t.Fatal("unexpected output")
	}
}

func TestRunCommandFailure(t *testing.T) {
	_, err := RunCommand("false")
	if err == nil {
		t.Fatal("expected command failure")
	}
}

func TestStartServiceFailure(t *testing.T) {
	_, err := StartService("this-service-does-not-exist")
	if err == nil {
		t.Fatal("expected error for invalid service")
	}
}
