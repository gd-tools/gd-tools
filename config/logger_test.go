package config

import (
	"fmt"
	"reflect"
	"testing"
)

type testLogger struct {
	info  []string
	debug []string
	errs  []string
}

func (l *testLogger) Info(args ...any) {
	l.info = append(l.info, fmt.Sprint(args...))
}

func (l *testLogger) Infof(format string, args ...any) {
	l.info = append(l.info, fmt.Sprintf(format, args...))
}

func (l *testLogger) Debug(args ...any) {
	l.debug = append(l.debug, fmt.Sprint(args...))
}

func (l *testLogger) Debugf(format string, args ...any) {
	l.debug = append(l.debug, fmt.Sprintf(format, args...))
}

func (l *testLogger) Error(args ...any) {
	l.errs = append(l.errs, fmt.Sprint(args...))
}

func (l *testLogger) Errorf(format string, args ...any) {
	l.errs = append(l.errs, fmt.Sprintf(format, args...))
}

func TestConfigInfo(t *testing.T) {
	log := &testLogger{}
	cfg := &Config{logger: log}

	cfg.Info("hello", " ", "world")

	want := []string{"hello world"}
	if !reflect.DeepEqual(log.info, want) {
		t.Fatalf("unexpected info logs: got %#v want %#v", log.info, want)
	}
}

func TestConfigInfof(t *testing.T) {
	log := &testLogger{}
	cfg := &Config{logger: log}

	cfg.Infof("hello %s", "world")

	want := []string{"hello world"}
	if !reflect.DeepEqual(log.info, want) {
		t.Fatalf("unexpected info logs: got %#v want %#v", log.info, want)
	}
}

func TestConfigDebug_DisabledWhenNotVerbose(t *testing.T) {
	log := &testLogger{}
	cfg := &Config{
		logger:  log,
		Verbose: false,
	}

	cfg.Debug("hidden")
	cfg.Debugf("hidden %d", 1)

	if len(log.debug) != 0 {
		t.Fatalf("expected no debug logs, got %#v", log.debug)
	}
}

func TestConfigDebug_EnabledWhenVerbose(t *testing.T) {
	log := &testLogger{}
	cfg := &Config{
		logger:  log,
		Verbose: true,
	}

	cfg.Debug("hello")
	cfg.Debugf("value=%d", 7)

	want := []string{"hello", "value=7"}
	if !reflect.DeepEqual(log.debug, want) {
		t.Fatalf("unexpected debug logs: got %#v want %#v", log.debug, want)
	}
}

func TestConfigError(t *testing.T) {
	log := &testLogger{}
	cfg := &Config{logger: log}

	cfg.Error("boom")
	cfg.Errorf("err=%d", 5)

	want := []string{"boom", "err=5"}
	if !reflect.DeepEqual(log.errs, want) {
		t.Fatalf("unexpected error logs: got %#v want %#v", log.errs, want)
	}
}

func TestConfigLogger_NilConfigDoesNothing(t *testing.T) {
	var cfg *Config

	cfg.Info("hello")
	cfg.Infof("hello %s", "world")
	cfg.Debug("debug")
	cfg.Debugf("debug %d", 1)
	cfg.Error("boom")
	cfg.Errorf("boom %d", 2)
}

func TestConfigLogger_NilLoggerDoesNothing(t *testing.T) {
	cfg := &Config{Verbose: true}

	cfg.Info("hello")
	cfg.Infof("hello %s", "world")
	cfg.Debug("debug")
	cfg.Debugf("debug %d", 1)
	cfg.Error("boom")
	cfg.Errorf("boom %d", 2)
}

func TestConfigRsyncFlags(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantOut string
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantOut: "-avzq",
		},
		{
			name:    "not verbose",
			cfg:     &Config{Verbose: false},
			wantOut: "-avzq",
		},
		{
			name:    "verbose",
			cfg:     &Config{Verbose: true},
			wantOut: "-avz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.RsyncFlags()
			if got != tt.wantOut {
				t.Fatalf("unexpected rsync flags: got %q want %q", got, tt.wantOut)
			}
		})
	}
}
