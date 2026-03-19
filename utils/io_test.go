package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJSON(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.json")

	err := os.WriteFile(name, []byte("{\n  \"name\": \"gd-tools\"\n}\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	var v struct {
		Name string `json:"name"`
	}

	err = LoadJSON(name, &v)
	if err != nil {
		t.Fatal(err)
	}

	if v.Name != "gd-tools" {
		t.Fatalf("unexpected name: %q", v.Name)
	}
}

func TestSaveJSON(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.json")

	v := struct {
		Name string `json:"name"`
	}{
		Name: "gd-tools",
	}

	err := SaveJSON(name, &v)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	want := "{\n  \"name\": \"gd-tools\"\n}\n"
	if string(data) != want {
		t.Fatalf("unexpected content:\n%s", string(data))
	}
}

func TestSaveFileCreate(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.txt")

	data := []byte("hello\n")

	err := SaveFile(name, data)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(data) {
		t.Fatalf("unexpected content: %q", string(got))
	}
}

func TestSaveFileReplace(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.txt")

	err := os.WriteFile(name, []byte("old\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("new\n")

	err = SaveFile(name, data)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(data) {
		t.Fatalf("unexpected content: %q", string(got))
	}
}

func TestSaveFileEqual(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.txt")

	data := []byte("same\n")

	err := os.WriteFile(name, data, 0o644)
	if err != nil {
		t.Fatal(err)
	}

	err = SaveFile(name, data)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(data) {
		t.Fatalf("unexpected content: %q", string(got))
	}
}
