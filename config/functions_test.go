package config

import (
	"bytes"
	"errors"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfigLoadFile_UsesInjectedFunction(t *testing.T) {
	cfg := &Config{}
	want := []byte("mock-data")

	called := false
	cfg.loadFile = func(name string) ([]byte, error) {
		called = true
		if name != "test.txt" {
			t.Fatalf("unexpected name: %q", name)
		}
		return want, nil
	}

	got, err := cfg.LoadFile("test.txt")
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected loadFile to be called")
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected data: got %q want %q", got, want)
	}
}

func TestConfigLoadFile_Fallback(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.txt")

	want := []byte("hello\n")
	if err := os.WriteFile(path, want, 0600); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{}
	got, err := cfg.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected data: got %q want %q", got, want)
	}
}

func TestConfigLoadJSON_UsesInjectedFunction(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}

	cfg := &Config{}
	called := false

	cfg.loadJSON = func(name string, v any) error {
		called = true
		if name != "test.json" {
			t.Fatalf("unexpected name: %q", name)
		}
		ptr, ok := v.(*sample)
		if !ok {
			t.Fatalf("unexpected target type: %T", v)
		}
		ptr.Name = "mock"
		return nil
	}

	var got sample
	err := cfg.LoadJSON("test.json", &got)
	if err != nil {
		t.Fatalf("LoadJSON returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected loadJSON to be called")
	}
	if got.Name != "mock" {
		t.Fatalf("unexpected value: %+v", got)
	}
}

func TestConfigLoadJSON_Fallback(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}

	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.json")

	if err := os.WriteFile(path, []byte(`{"name":"real"}`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{}
	var got sample

	err := cfg.LoadJSON(path, &got)
	if err != nil {
		t.Fatalf("LoadJSON returned error: %v", err)
	}
	if got.Name != "real" {
		t.Fatalf("unexpected value: %+v", got)
	}
}

func TestConfigSaveFile_UsesInjectedFunction(t *testing.T) {
	cfg := &Config{}
	called := false

	cfg.saveFile = func(name string, data []byte) error {
		called = true
		if name != "out.txt" {
			t.Fatalf("unexpected name: %q", name)
		}
		if string(data) != "abc" {
			t.Fatalf("unexpected data: %q", data)
		}
		return nil
	}

	err := cfg.SaveFile("out.txt", []byte("abc"))
	if err != nil {
		t.Fatalf("SaveFile returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected saveFile to be called")
	}
}

func TestConfigSaveFile_Fallback(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "out.txt")

	cfg := &Config{}
	err := cfg.SaveFile(path, []byte("abc"))
	if err != nil {
		t.Fatalf("SaveFile returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "abc" {
		t.Fatalf("unexpected file content: %q", got)
	}
}

func TestConfigSaveJSON_UsesInjectedFunction(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}

	cfg := &Config{}
	called := false

	cfg.saveJSON = func(name string, v any) error {
		called = true
		if name != "out.json" {
			t.Fatalf("unexpected name: %q", name)
		}
		ptr, ok := v.(*sample)
		if !ok {
			t.Fatalf("unexpected value type: %T", v)
		}
		if ptr.Name != "mock" {
			t.Fatalf("unexpected value: %+v", ptr)
		}
		return nil
	}

	err := cfg.SaveJSON("out.json", &sample{Name: "mock"})
	if err != nil {
		t.Fatalf("SaveJSON returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected saveJSON to be called")
	}
}

func TestConfigSaveJSON_Fallback(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}

	tmp := t.TempDir()
	path := filepath.Join(tmp, "out.json")

	cfg := &Config{}
	err := cfg.SaveJSON(path, &sample{Name: "real"})
	if err != nil {
		t.Fatalf("SaveJSON returned error: %v", err)
	}

	var got sample
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("expected JSON file to be non-empty")
	}

	err = cfg.LoadJSON(path, &got)
	if err != nil {
		t.Fatalf("LoadJSON returned error: %v", err)
	}
	if got.Name != "real" {
		t.Fatalf("unexpected value: %+v", got)
	}
}

func TestConfigLookupIP_UsesInjectedFunction(t *testing.T) {
	cfg := &Config{}
	want := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("::1"),
	}

	called := false
	cfg.lookupIP = func(host string) ([]net.IP, error) {
		called = true
		if host != "example.test" {
			t.Fatalf("unexpected host: %q", host)
		}
		return want, nil
	}

	got, err := cfg.LookupIP("example.test")
	if err != nil {
		t.Fatalf("LookupIP returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected lookupIP to be called")
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected result: got %v want %v", got, want)
	}
}

func TestConfigDHParams_UsesInjectedFunction(t *testing.T) {
	cfg := &Config{}
	want := []byte("dhparams")

	called := false
	cfg.dhParams = func(bits int) ([]byte, error) {
		called = true
		if bits != 2048 {
			t.Fatalf("unexpected bits: %d", bits)
		}
		return want, nil
	}

	got, err := cfg.DHParams(2048)
	if err != nil {
		t.Fatalf("DHParams returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected dhParams to be called")
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected result: got %q want %q", got, want)
	}
}

func TestConfigRSAKeyPair_UsesInjectedFunction(t *testing.T) {
	cfg := &Config{}
	wantPub := []byte("pub")
	wantKey := []byte("key")

	called := false
	cfg.rsaKeyPair = func(fqdn string) ([]byte, []byte, error) {
		called = true
		if fqdn != "prod.example.test" {
			t.Fatalf("unexpected fqdn: %q", fqdn)
		}
		return wantPub, wantKey, nil
	}

	gotPub, gotKey, err := cfg.RSAKeyPair("prod.example.test")
	if err != nil {
		t.Fatalf("RSAKeyPair returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected rsaKeyPair to be called")
	}
	if !bytes.Equal(gotPub, wantPub) {
		t.Fatalf("unexpected public key: got %q want %q", gotPub, wantPub)
	}
	if !bytes.Equal(gotKey, wantKey) {
		t.Fatalf("unexpected private key: got %q want %q", gotKey, wantKey)
	}
}

func TestConfigRunShell_UsesInjectedFunction(t *testing.T) {
	cfg := &Config{}
	want := []byte("ok")

	called := false
	cfg.runShell = func(commands []string) ([]byte, error) {
		called = true
		if !reflect.DeepEqual(commands, []string{"echo a", "echo b"}) {
			t.Fatalf("unexpected commands: %#v", commands)
		}
		return want, nil
	}

	got, err := cfg.RunShell([]string{"echo a", "echo b"})
	if err != nil {
		t.Fatalf("RunShell returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected runShell to be called")
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected result: got %q want %q", got, want)
	}
}

func TestConfigRunCommand_UsesInjectedFunction(t *testing.T) {
	cfg := &Config{}
	want := []byte("ok")

	called := false
	cfg.runCommand = func(name string, args ...string) ([]byte, error) {
		called = true
		if name != "openssl" {
			t.Fatalf("unexpected command: %q", name)
		}
		if !reflect.DeepEqual(args, []string{"version"}) {
			t.Fatalf("unexpected args: %#v", args)
		}
		return want, nil
	}

	got, err := cfg.RunCommand("openssl", "version")
	if err != nil {
		t.Fatalf("RunCommand returned error: %v", err)
	}
	if !called {
		t.Fatal("expected injected runCommand to be called")
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected result: got %q want %q", got, want)
	}
}

func TestConfigRunCommand_PropagatesError(t *testing.T) {
	cfg := &Config{}
	wantErr := errors.New("boom")

	cfg.runCommand = func(name string, args ...string) ([]byte, error) {
		return nil, wantErr
	}

	_, err := cfg.RunCommand("openssl", "version")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
