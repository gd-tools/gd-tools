package utils

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func withTempDir(t *testing.T, fn func()) {
	t.Helper()

	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(old)
	}()

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	fn()
}

func TestGenerateSecretModes(t *testing.T) {
	hash, err := GenerateSecret("secret", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("secret")); err != nil {
		t.Fatal("bcrypt mismatch")
	}

	pbkdf, err := GenerateSecret("secret", "pbkdf2")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(pbkdf, "$1$") {
		t.Fatalf("unexpected pbkdf2 format: %s", pbkdf)
	}
}

func TestGenerateSecretInvalidMode(t *testing.T) {
	_, err := GenerateSecret("x", "invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGenerateBcryptEmpty(t *testing.T) {
	_, err := GenerateBcrypt("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGeneratePBKDF2Empty(t *testing.T) {
	_, err := GeneratePBKDF2("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreatePassword(t *testing.T) {
	pw, err := CreatePassword(20)
	if err != nil {
		t.Fatal(err)
	}
	if len(pw) != 20 {
		t.Fatalf("len=%d", len(pw))
	}
	if strings.ContainsAny(pw, "lO") {
		t.Fatal("contains ambiguous chars")
	}
}

func TestLoadSecretsMissingFile(t *testing.T) {
	withTempDir(t, func() {
		list, err := LoadSecrets()
		if err != nil {
			t.Fatal(err)
		}
		if list == nil {
			t.Fatal("expected list")
		}
		if len(list.Secrets) != 0 {
			t.Fatalf("len=%d want 0", len(list.Secrets))
		}
	})
}

func TestSetAndGet(t *testing.T) {
	withTempDir(t, func() {
		list := &SecretList{}

		if err := list.Set("db", "postgres", "in", "out"); err != nil {
			t.Fatal(err)
		}

		loaded, err := LoadSecrets()
		if err != nil {
			t.Fatal(err)
		}

		e := loaded.Get("db", "postgres")
		if e == nil {
			t.Fatal("missing entry")
		}
		if e.Input != "in" || e.Output != "out" {
			t.Fatal("wrong values")
		}
	})
}

func TestSetUpdatesOnlyOnInputChange(t *testing.T) {
	withTempDir(t, func() {
		list := &SecretList{}

		if err := list.Set("db", "postgres", "in1", "out1"); err != nil {
			t.Fatal(err)
		}
		if err := list.Set("db", "postgres", "in1", "out2"); err != nil {
			t.Fatal(err)
		}

		e := list.Get("db", "postgres")
		if e == nil {
			t.Fatal("missing entry")
		}
		if e.Output != "out1" {
			t.Fatal("should not update if input unchanged")
		}

		if err := list.Set("db", "postgres", "in2", "out2"); err != nil {
			t.Fatal(err)
		}

		e = list.Get("db", "postgres")
		if e == nil {
			t.Fatal("missing entry")
		}
		if e.Output != "out2" {
			t.Fatal("should update on input change")
		}
	})
}

func TestSetMailUser(t *testing.T) {
	withTempDir(t, func() {
		list := &SecretList{}

		pw, hash, err := list.SetMailUser("user@example.org", "secret")
		if err != nil {
			t.Fatal(err)
		}

		if pw != "secret" {
			t.Fatal("wrong password")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("secret")); err != nil {
			t.Fatal("bcrypt mismatch")
		}

		e := list.Get(MailUserScope, "user@example.org")
		if e == nil {
			t.Fatal("entry missing")
		}
		if e.Input != "secret" {
			t.Fatal("wrong input")
		}
	})
}

func TestSetMailUserGeneratesPassword(t *testing.T) {
	withTempDir(t, func() {
		list := &SecretList{}

		pw, hash, err := list.SetMailUser("user@example.org", "")
		if err != nil {
			t.Fatal(err)
		}
		if len(pw) != 20 {
			t.Fatalf("len=%d want 20", len(pw))
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)); err != nil {
			t.Fatal("bcrypt mismatch")
		}
	})
}

func TestEnsurePassword(t *testing.T) {
	withTempDir(t, func() {
		pw1, err := EnsurePassword(16, "db", "postgres")
		if err != nil {
			t.Fatal(err)
		}

		pw2, err := EnsurePassword(16, "db", "postgres")
		if err != nil {
			t.Fatal(err)
		}

		if pw1 != pw2 {
			t.Fatal("should reuse existing password")
		}
	})
}

func TestSaveSortsByScopeAndName(t *testing.T) {
	withTempDir(t, func() {
		list := &SecretList{
			Secrets: []Secret{
				{Scope: "mail", Name: "z@example.org"},
				{Scope: "db", Name: "z"},
				{Scope: "db", Name: "a"},
			},
		}

		if err := list.Save(); err != nil {
			t.Fatal(err)
		}

		loaded, err := LoadSecrets()
		if err != nil {
			t.Fatal(err)
		}

		if len(loaded.Secrets) != 3 {
			t.Fatalf("len=%d want 3", len(loaded.Secrets))
		}

		got := []string{
			loaded.Secrets[0].Scope + ":" + loaded.Secrets[0].Name,
			loaded.Secrets[1].Scope + ":" + loaded.Secrets[1].Name,
			loaded.Secrets[2].Scope + ":" + loaded.Secrets[2].Name,
		}
		want := []string{
			"db:a",
			"db:z",
			"mail:z@example.org",
		}

		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("sorted entry %d = %q, want %q", i, got[i], want[i])
			}
		}
	})
}

func TestLoadSecretsInvalidJSON(t *testing.T) {
	withTempDir(t, func() {
		if err := os.WriteFile(SecretsFile, []byte("{invalid"), 0o644); err != nil {
			t.Fatal(err)
		}

		_, err := LoadSecrets()
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
