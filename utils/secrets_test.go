package utils

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateBcrypt(t *testing.T) {
	hash, err := GenerateBcrypt("secret")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(hash, "$2") {
		t.Fatalf("unexpected bcrypt hash: %s", hash)
	}
}

func TestGeneratePBKDF2(t *testing.T) {
	hash, err := GeneratePBKDF2("secret")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(hash, "$1$") {
		t.Fatalf("unexpected pbkdf2 format: %s", hash)
	}
}

func TestCreatePassword(t *testing.T) {
	pw, err := CreatePassword(20)
	if err != nil {
		t.Fatal(err)
	}

	if len(pw) != 20 {
		t.Fatalf("unexpected password length: %d", len(pw))
	}
}

func TestSecretListSetGet(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)

	os.Chdir(dir)

	list := &SecretList{}

	err := list.Set("example.org", "user1", "in", "out")
	if err != nil {
		t.Fatal(err)
	}

	list2, err := LoadSecrets()
	if err != nil {
		t.Fatal(err)
	}

	entry := list2.Get("example.org", "user1")
	if entry == nil {
		t.Fatal("secret not found")
	}

	if entry.Output != "out" {
		t.Fatalf("unexpected value: %s", entry.Output)
	}
}

func TestFetchPassword(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)

	os.Chdir(dir)

	pw, err := FetchPassword(16, "example.org", "user2")
	if err != nil {
		t.Fatal(err)
	}

	if len(pw) != 16 {
		t.Fatalf("unexpected password length")
	}

	pw2, err := FetchPassword(16, "example.org", "user2")
	if err != nil {
		t.Fatal(err)
	}

	if pw != pw2 {
		t.Fatalf("password should persist")
	}
}
