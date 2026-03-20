package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerInfo(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		Writer: &buf,
	}

	logger.Info("hello")

	out := buf.String()
	if !strings.Contains(out, "[info] hello") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggerInfoList(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		Writer: &buf,
	}

	logger.Info([]string{"one", "two"})

	out := buf.String()
	if !strings.Contains(out, "[info] one") {
		t.Fatalf("missing first line: %q", out)
	}
	if !strings.Contains(out, "[info] two") {
		t.Fatalf("missing second line: %q", out)
	}
}

func TestLoggerInfof(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		Writer: &buf,
	}

	logger.Infof("value=%d", 42)

	out := buf.String()
	if !strings.Contains(out, "[info] value=42") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggerDebugDisabled(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		Writer:  &buf,
		Verbose: false,
	}

	logger.Debug("hello")
	logger.Debugf("value=%d", 42)

	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestLoggerDebugEnabled(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		Writer:  &buf,
		Verbose: true,
	}

	logger.Debug("hello")
	logger.Debugf("value=%d", 42)

	out := buf.String()
	if !strings.Contains(out, "[dbg] hello") {
		t.Fatalf("missing debug line: %q", out)
	}
	if !strings.Contains(out, "[dbg] value=42") {
		t.Fatalf("missing debugf line: %q", out)
	}
}

func TestLoggerDebugList(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		Writer:  &buf,
		Verbose: true,
	}

	logger.Debug([]string{"one", "two"})

	out := buf.String()
	if !strings.Contains(out, "[dbg] one") {
		t.Fatalf("missing first line: %q", out)
	}
	if !strings.Contains(out, "[dbg] two") {
		t.Fatalf("missing second line: %q", out)
	}
}

func TestLoggerError(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		Writer: &buf,
	}

	logger.Error("boom")

	out := buf.String()
	if !strings.Contains(out, "[err] boom") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggerErrorList(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		Writer: &buf,
	}

	logger.Error([]string{"one", "two"})

	out := buf.String()
	if !strings.Contains(out, "[err] one") {
		t.Fatalf("missing first line: %q", out)
	}
	if !strings.Contains(out, "[err] two") {
		t.Fatalf("missing second line: %q", out)
	}
}

func TestLoggerErrorf(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		Writer: &buf,
	}

	logger.Errorf("value=%d", 42)

	out := buf.String()
	if !strings.Contains(out, "[err] value=42") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggerNilWriter(t *testing.T) {
	logger := &Logger{
		Writer: nil,
	}

	logger.Info("hello")
	logger.Infof("value=%d", 1)
	logger.Debug("hidden")
	logger.Debugf("value=%d", 2)
	logger.Error("boom")
	logger.Errorf("value=%d", 3)
}

func TestLoggerNilReceiver(t *testing.T) {
	var logger *Logger

	logger.Info("hello")
	logger.Infof("value=%d", 1)
	logger.Debug("hidden")
	logger.Debugf("value=%d", 2)
	logger.Error("boom")
	logger.Errorf("value=%d", 3)
}
