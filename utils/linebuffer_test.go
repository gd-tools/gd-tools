package utils

import (
	"strings"
	"testing"
)

func TestLineBufferAddAndLines(t *testing.T) {
	var lb LineBuffer

	lb.Add("one")
	lb.Add("two")

	lines := lb.Lines()

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	if lines[0] != "one" || lines[1] != "two" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestLineBufferAddf(t *testing.T) {
	var lb LineBuffer

	lb.Addf("hello %s", "world")

	lines := lb.Lines()

	if len(lines) != 1 {
		t.Fatalf("expected 1 line")
	}

	if lines[0] != "hello world" {
		t.Fatalf("unexpected line: %s", lines[0])
	}
}

func TestLineBufferText(t *testing.T) {
	var lb LineBuffer

	lb.Add("alpha")
	lb.Add("beta")

	text := lb.Text()

	expected := "alpha\nbeta\n"

	if text != expected {
		t.Fatalf("unexpected text:\n%s", text)
	}
}

func TestLineBufferEmpty(t *testing.T) {
	var lb LineBuffer

	text := lb.Text()

	if !strings.HasSuffix(text, "\n") {
		t.Fatalf("text should end with newline")
	}
}
