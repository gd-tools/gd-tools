package config

import (
	"errors"
	"net"
	"strings"
	"testing"
)

func TestConfigLookupIPOverride(t *testing.T) {
	called := false
	want := []net.IP{
		net.ParseIP("192.0.2.10"),
		net.ParseIP("2001:db8::10"),
	}

	cfg := &Config{
		lookupIP: func(host string) ([]net.IP, error) {
			called = true
			if host != "example.test" {
				t.Fatalf("unexpected host: got %q", host)
			}
			return want, nil
		},
	}

	got, err := cfg.LookupIP("example.test")
	if err != nil {
		t.Fatalf("LookupIP returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
	if len(got) != len(want) {
		t.Fatalf("unexpected IP count: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if !got[i].Equal(want[i]) {
			t.Fatalf("unexpected IP at index %d: got %v want %v", i, got[i], want[i])
		}
	}
}

func TestConfigLookupIPFallback(t *testing.T) {
	cfg := &Config{}

	got, err := cfg.LookupIP("localhost")
	if err != nil {
		t.Fatalf("LookupIP returned error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected at least one IP for localhost")
	}
}

func TestConfigDHParamsOverride(t *testing.T) {
	called := false
	want := []byte("dhparams-data")

	cfg := &Config{
		dhParams: func(bits int) ([]byte, error) {
			called = true
			if bits != 2048 {
				t.Fatalf("unexpected bits: got %d", bits)
			}
			return want, nil
		},
	}

	got, err := cfg.DHParams(2048)
	if err != nil {
		t.Fatalf("DHParams returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
	if string(got) != string(want) {
		t.Fatalf("unexpected result: got %q want %q", string(got), string(want))
	}
}

func TestConfigRSAKeyPairOverride(t *testing.T) {
	called := false
	wantPriv := []byte("private-key")
	wantPub := []byte("public-key")

	cfg := &Config{
		rsaKeyPair: func(fqdn string) ([]byte, []byte, error) {
			called = true
			if fqdn != "host.example.test" {
				t.Fatalf("unexpected fqdn: got %q", fqdn)
			}
			return wantPriv, wantPub, nil
		},
	}

	gotPriv, gotPub, err := cfg.RSAKeyPair("host.example.test")
	if err != nil {
		t.Fatalf("RSAKeyPair returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
	if string(gotPriv) != string(wantPriv) {
		t.Fatalf("unexpected private key: got %q want %q", string(gotPriv), string(wantPriv))
	}
	if string(gotPub) != string(wantPub) {
		t.Fatalf("unexpected public key: got %q want %q", string(gotPub), string(wantPub))
	}
}

func TestConfigRunShellOverride(t *testing.T) {
	called := false
	want := []byte("override-output")

	cfg := &Config{
		runShell: func(commands []string) ([]byte, error) {
			called = true
			if len(commands) != 2 {
				t.Fatalf("unexpected command count: got %d", len(commands))
			}
			if commands[0] != "echo one" || commands[1] != "echo two" {
				t.Fatalf("unexpected commands: got %#v", commands)
			}
			return want, nil
		},
	}

	got, err := cfg.RunShell([]string{"echo one", "echo two"})
	if err != nil {
		t.Fatalf("RunShell returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
	if string(got) != string(want) {
		t.Fatalf("unexpected output: got %q want %q", string(got), string(want))
	}
}

func TestConfigRunShellFallback(t *testing.T) {
	cfg := &Config{}

	got, err := cfg.RunShell([]string{"printf 'hello'"})
	if err != nil {
		t.Fatalf("RunShell returned error: %v", err)
	}

	if strings.TrimSpace(string(got)) != "hello" {
		t.Fatalf("unexpected output: got %q", string(got))
	}
}

func TestConfigRunCommandOverride(t *testing.T) {
	called := false
	want := []byte("override-command")

	cfg := &Config{
		runCommand: func(name string, args ...string) ([]byte, error) {
			called = true
			if name != "printf" {
				t.Fatalf("unexpected command name: got %q", name)
			}
			if len(args) != 1 || args[0] != "hello" {
				t.Fatalf("unexpected args: got %#v", args)
			}
			return want, nil
		},
	}

	got, err := cfg.RunCommand("printf", "hello")
	if err != nil {
		t.Fatalf("RunCommand returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
	if string(got) != string(want) {
		t.Fatalf("unexpected output: got %q want %q", string(got), string(want))
	}
}

func TestConfigRunCommandFallback(t *testing.T) {
	cfg := &Config{}

	got, err := cfg.RunCommand("printf", "hello")
	if err != nil {
		t.Fatalf("RunCommand returned error: %v", err)
	}
	if string(got) != "hello" {
		t.Fatalf("unexpected output: got %q want %q", string(got), "hello")
	}
}

func TestConfigFunctionOverrideErrorIsReturned(t *testing.T) {
	want := errors.New("boom")

	cfg := &Config{
		runCommand: func(name string, args ...string) ([]byte, error) {
			return nil, want
		},
	}

	_, err := cfg.RunCommand("false")
	if !errors.Is(err, want) {
		t.Fatalf("unexpected error: got %v want %v", err, want)
	}
}
