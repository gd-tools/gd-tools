package utils

import (
	"os"
	"testing"
)

func TestGenerateDHParamsExisting(t *testing.T) {
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

	if err := os.WriteFile(DHParamsFile, expected, 0644); err != nil {
		t.Fatal(err)
	}

	data, err := GenerateDHParams(4096)
	if err != nil {
		t.Fatalf("GenerateDHParams failed: %v", err)
	}

	if string(data) != string(expected) {
		t.Fatalf("unexpected dhparams content")
	}
}
