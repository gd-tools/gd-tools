package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRSAKeyPairCreatesFiles(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen not found")
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir(%q) failed: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatalf("restore Chdir(%q) failed: %v", oldWD, err)
		}
	})

	priv, publ, err := RSAKeyPair("example.org")
	if err != nil {
		t.Fatalf("ProdRSAKeyPair returned error: %v", err)
	}

	if len(priv) == 0 {
		t.Fatal("expected private key data")
	}
	if len(publ) == 0 {
		t.Fatal("expected public key data")
	}

	if _, err := os.Stat(filepath.Join(dir, "root_id_rsa")); err != nil {
		t.Fatalf("missing root_id_rsa: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "root_id_rsa.pub")); err != nil {
		t.Fatalf("missing root_id_rsa.pub: %v", err)
	}
}

func TestRSAKeyPairReadsExistingFiles(t *testing.T) {
	t.Helper()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir(%q) failed: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatalf("restore Chdir(%q) failed: %v", oldWD, err)
		}
	})

	wantPriv := []byte("test-private-key\n")
	wantPubl := []byte("test-public-key\n")

	if err := os.WriteFile("root_id_rsa", wantPriv, 0o600); err != nil {
		t.Fatalf("WriteFile(root_id_rsa) failed: %v", err)
	}
	if err := os.WriteFile("root_id_rsa.pub", wantPubl, 0o644); err != nil {
		t.Fatalf("WriteFile(root_id_rsa.pub) failed: %v", err)
	}

	priv, publ, err := RSAKeyPair("example.org")
	if err != nil {
		t.Fatalf("ProdRSAKeyPair returned error: %v", err)
	}

	if string(priv) != string(wantPriv) {
		t.Fatalf("unexpected private key: got %q want %q", string(priv), string(wantPriv))
	}
	if string(publ) != string(wantPubl) {
		t.Fatalf("unexpected public key: got %q want %q", string(publ), string(wantPubl))
	}
}

func TestRSAKeyPairFailsWithoutPublicKey(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen not found")
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir(%q) failed: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatalf("restore Chdir(%q) failed: %v", oldWD, err)
		}
	})

	if err := os.WriteFile("root_id_rsa", []byte("broken-private-key\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(root_id_rsa) failed: %v", err)
	}

	_, _, err = RSAKeyPair("example.org")
	if err == nil {
		t.Fatal("expected error")
	}
}
