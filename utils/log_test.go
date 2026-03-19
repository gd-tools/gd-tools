package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogSay(t *testing.T) {
	var buf bytes.Buffer

	log := &Log{
		Writer: &buf,
	}

	log.Say("hello")

	out := buf.String()
	if !strings.Contains(out, "[run] hello") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLogSayList(t *testing.T) {
	var buf bytes.Buffer

	log := &Log{
		Writer: &buf,
	}

	log.Say([]string{"one", "two"})

	out := buf.String()
	if !strings.Contains(out, "[run] one") {
		t.Fatalf("missing first line: %q", out)
	}
	if !strings.Contains(out, "[run] two") {
		t.Fatalf("missing second line: %q", out)
	}
}

func TestLogSayf(t *testing.T) {
	var buf bytes.Buffer

	log := &Log{
		Writer: &buf,
	}

	log.Sayf("value=%d", 42)

	out := buf.String()
	if !strings.Contains(out, "[run] value=42") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLogDebugDisabled(t *testing.T) {
	var buf bytes.Buffer

	log := &Log{
		Writer:  &buf,
		Verbose: false,
	}

	log.Debug("hello")
	log.Debugf("value=%d", 42)

	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestLogDebugEnabled(t *testing.T) {
	var buf bytes.Buffer

	log := &Log{
		Writer:  &buf,
		Verbose: true,
	}

	log.Debug("hello")
	log.Debugf("value=%d", 42)

	out := buf.String()
	if !strings.Contains(out, "[dbg] ################### hello") {
		t.Fatalf("missing debug line: %q", out)
	}
	if !strings.Contains(out, "[dbg] ################### value=42") {
		t.Fatalf("missing debugf line: %q", out)
	}
}

func TestLogDebugList(t *testing.T) {
	var buf bytes.Buffer

	log := &Log{
		Writer:  &buf,
		Verbose: true,
	}

	log.Debug([]string{"one", "two"})

	out := buf.String()
	if !strings.Contains(out, "[dbg] ################### one") {
		t.Fatalf("missing first line: %q", out)
	}
	if !strings.Contains(out, "[dbg] ################### two") {
		t.Fatalf("missing second line: %q", out)
	}
}

func TestLogNilWriter(t *testing.T) {
	log := &Log{
		Writer: nil,
	}

	log.Say("hello")
	log.Sayf("value=%d", 1)
	log.Debug("hidden")
	log.Debugf("value=%d", 2)
}

func TestLogNilReceiver(t *testing.T) {
	var log *Log

	log.Say("hello")
	log.Sayf("value=%d", 1)
	log.Debug("hidden")
	log.Debugf("value=%d", 2)
}
