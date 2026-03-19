package utils

import (
	"os"
	"testing"
)

func TestDHParamsExisting(t *testing.T) {
	dir := t.TempDir()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	expected := []byte("TEST_DH_PARAMS")

	if err := os.WriteFile(DHParamsFile, expected, 0o644); err != nil {
		t.Fatal(err)
	}

	data, err := DHParams(4096)
	if err != nil {
		t.Fatalf("DHParams failed: %v", err)
	}

	if string(data) != string(expected) {
		t.Fatalf("unexpected dhparams content")
	}
}
