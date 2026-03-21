package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigMkdirAllOverride(t *testing.T) {
	called := false
	cfg := &Config{
		mkdirAll: func(path string, perm os.FileMode) error {
			called = true
			if path != "/tmp/testdir" {
				t.Fatalf("unexpected path: got %q", path)
			}
			if perm != 0o755 {
				t.Fatalf("unexpected perm: got %v", perm)
			}
			return nil
		},
	}

	err := cfg.MkdirAll("/tmp/testdir", 0o755)
	if err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
}

func TestConfigMkdirAllFallback(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "a", "b", "c")

	cfg := &Config{}
	if err := cfg.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}

	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("Stat returned error: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected created path to be a directory")
	}
}

func TestConfigSetenvOverride(t *testing.T) {
	called := false
	cfg := &Config{
		setenv: func(key, value string) error {
			called = true
			if key != "GD_TOOLS_TEST_KEY" {
				t.Fatalf("unexpected key: got %q", key)
			}
			if value != "test-value" {
				t.Fatalf("unexpected value: got %q", value)
			}
			return nil
		},
	}

	err := cfg.Setenv("GD_TOOLS_TEST_KEY", "test-value")
	if err != nil {
		t.Fatalf("Setenv returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
}

func TestConfigSetenvFallback(t *testing.T) {
	const key = "GD_TOOLS_TEST_SETENV"

	cfg := &Config{}
	if err := cfg.Setenv(key, "value-123"); err != nil {
		t.Fatalf("Setenv returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Unsetenv(key)
	})

	got := os.Getenv(key)
	if got != "value-123" {
		t.Fatalf("unexpected env value: got %q want %q", got, "value-123")
	}
}

func TestConfigUnsetenvOverride(t *testing.T) {
	called := false
	cfg := &Config{
		unsetenv: func(key string) error {
			called = true
			if key != "GD_TOOLS_TEST_KEY" {
				t.Fatalf("unexpected key: got %q", key)
			}
			return nil
		},
	}

	err := cfg.Unsetenv("GD_TOOLS_TEST_KEY")
	if err != nil {
		t.Fatalf("Unsetenv returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
}

func TestConfigUnsetenvFallback(t *testing.T) {
	const key = "GD_TOOLS_TEST_UNSETENV"

	if err := os.Setenv(key, "value-123"); err != nil {
		t.Fatalf("Setenv preparation failed: %v", err)
	}

	cfg := &Config{}
	if err := cfg.Unsetenv(key); err != nil {
		t.Fatalf("Unsetenv returned error: %v", err)
	}

	_, ok := os.LookupEnv(key)
	if ok {
		t.Fatal("expected environment variable to be removed")
	}
}

func TestConfigLoadFileOverride(t *testing.T) {
	called := false
	cfg := &Config{
		loadFile: func(name string) ([]byte, error) {
			called = true
			if name != "test.txt" {
				t.Fatalf("unexpected name: got %q", name)
			}
			return []byte("hello"), nil
		},
	}

	data, err := cfg.LoadFile("test.txt")
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected data: got %q", string(data))
	}
}

func TestConfigLoadFileFallback(t *testing.T) {
	base := t.TempDir()
	name := filepath.Join(base, "test.txt")

	if err := os.WriteFile(name, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	cfg := &Config{}
	data, err := cfg.LoadFile(name)
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	if string(data) != "hello world" {
		t.Fatalf("unexpected data: got %q", string(data))
	}
}

func TestConfigLoadJSONOverride(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}

	called := false
	cfg := &Config{
		loadJSON: func(name string, v any) error {
			called = true
			if name != "test.json" {
				t.Fatalf("unexpected name: got %q", name)
			}
			dst, ok := v.(*sample)
			if !ok {
				t.Fatalf("unexpected target type: %T", v)
			}
			dst.Name = "override"
			return nil
		},
	}

	var got sample
	err := cfg.LoadJSON("test.json", &got)
	if err != nil {
		t.Fatalf("LoadJSON returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
	if got.Name != "override" {
		t.Fatalf("unexpected value: got %q", got.Name)
	}
}

func TestConfigLoadJSONFallback(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}

	base := t.TempDir()
	name := filepath.Join(base, "test.json")

	if err := os.WriteFile(name, []byte(`{"name":"fallback"}`), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	cfg := &Config{}
	var got sample
	if err := cfg.LoadJSON(name, &got); err != nil {
		t.Fatalf("LoadJSON returned error: %v", err)
	}
	if got.Name != "fallback" {
		t.Fatalf("unexpected value: got %q", got.Name)
	}
}

func TestConfigSaveFileOverride(t *testing.T) {
	called := false
	cfg := &Config{
		saveFile: func(name string, data []byte) error {
			called = true
			if name != "test.txt" {
				t.Fatalf("unexpected name: got %q", name)
			}
			if string(data) != "payload" {
				t.Fatalf("unexpected data: got %q", string(data))
			}
			return nil
		},
	}

	err := cfg.SaveFile("test.txt", []byte("payload"))
	if err != nil {
		t.Fatalf("SaveFile returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
}

func TestConfigSaveFileFallback(t *testing.T) {
	base := t.TempDir()
	name := filepath.Join(base, "test.txt")

	cfg := &Config{}
	if err := cfg.SaveFile(name, []byte("payload")); err != nil {
		t.Fatalf("SaveFile returned error: %v", err)
	}

	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) != "payload" {
		t.Fatalf("unexpected data: got %q", string(data))
	}
}

func TestConfigSaveJSONOverride(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}

	called := false
	cfg := &Config{
		saveJSON: func(name string, v any) error {
			called = true
			if name != "test.json" {
				t.Fatalf("unexpected name: got %q", name)
			}
			got, ok := v.(*sample)
			if !ok {
				t.Fatalf("unexpected value type: %T", v)
			}
			if got.Name != "payload" {
				t.Fatalf("unexpected struct value: got %q", got.Name)
			}
			return nil
		},
	}

	err := cfg.SaveJSON("test.json", &sample{Name: "payload"})
	if err != nil {
		t.Fatalf("SaveJSON returned error: %v", err)
	}
	if !called {
		t.Fatal("expected override to be called")
	}
}

func TestConfigSaveJSONFallback(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}

	base := t.TempDir()
	name := filepath.Join(base, "test.json")

	cfg := &Config{}
	if err := cfg.SaveJSON(name, &sample{Name: "payload"}); err != nil {
		t.Fatalf("SaveJSON returned error: %v", err)
	}

	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) == "" {
		t.Fatal("expected JSON file to be written")
	}

	var got sample
	if err := cfg.LoadJSON(name, &got); err != nil {
		t.Fatalf("LoadJSON returned error: %v", err)
	}
	if got.Name != "payload" {
		t.Fatalf("unexpected value: got %q", got.Name)
	}
}

func TestConfigOverrideErrorIsReturned(t *testing.T) {
	want := errors.New("boom")

	cfg := &Config{
		loadFile: func(name string) ([]byte, error) {
			return nil, want
		},
	}

	_, err := cfg.LoadFile("x")
	if !errors.Is(err, want) {
		t.Fatalf("unexpected error: got %v want %v", err, want)
	}
}
