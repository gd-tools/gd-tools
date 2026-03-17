package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func useTempWorkDir(t *testing.T) string {
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

	return dir
}

func TestLoadUserIDsMissingFile(t *testing.T) {
	useTempWorkDir(t)

	list, err := LoadUserIDs()
	if err != nil {
		t.Fatalf("LoadUserIDs returned error: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %d entries", len(list))
	}
}

func TestSaveAndLoadUserIDs(t *testing.T) {
	useTempWorkDir(t)

	err := SaveUserIDs(UserID{
		Name: "www-data",
		UID:  "33",
		GID:  "33",
	})
	if err != nil {
		t.Fatalf("SaveUserIDs returned error: %v", err)
	}

	list, err := LoadUserIDs()
	if err != nil {
		t.Fatalf("LoadUserIDs returned error: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}

	if list[0].Name != "www-data" {
		t.Fatalf("unexpected name: %q", list[0].Name)
	}
	if list[0].UID != "33" {
		t.Fatalf("unexpected UID: %q", list[0].UID)
	}
	if list[0].GID != "33" {
		t.Fatalf("unexpected GID: %q", list[0].GID)
	}
}

func TestSaveUserIDsUpdatesExistingEntry(t *testing.T) {
	useTempWorkDir(t)

	err := SaveUserIDs(UserID{
		Name: "www-data",
		UID:  "33",
		GID:  "33",
	})
	if err != nil {
		t.Fatalf("first SaveUserIDs returned error: %v", err)
	}

	err = SaveUserIDs(UserID{
		Name: "www-data",
		UID:  "100",
		GID:  "101",
	})
	if err != nil {
		t.Fatalf("second SaveUserIDs returned error: %v", err)
	}

	list, err := LoadUserIDs()
	if err != nil {
		t.Fatalf("LoadUserIDs returned error: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}

	if list[0].UID != "100" {
		t.Fatalf("unexpected UID: %q", list[0].UID)
	}
	if list[0].GID != "101" {
		t.Fatalf("unexpected GID: %q", list[0].GID)
	}
}

func TestSaveUserIDsSortedByName(t *testing.T) {
	useTempWorkDir(t)

	for _, u := range []UserID{
		{Name: "zeta", UID: "3", GID: "3"},
		{Name: "alpha", UID: "1", GID: "1"},
		{Name: "beta", UID: "2", GID: "2"},
	} {
		if err := SaveUserIDs(u); err != nil {
			t.Fatalf("SaveUserIDs(%q) returned error: %v", u.Name, err)
		}
	}

	list, err := LoadUserIDs()
	if err != nil {
		t.Fatalf("LoadUserIDs returned error: %v", err)
	}

	if len(list) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(list))
	}

	want := []string{"alpha", "beta", "zeta"}
	for i := range want {
		if list[i].Name != want[i] {
			t.Fatalf("entry %d: got %q want %q", i, list[i].Name, want[i])
		}
	}
}

func TestGetUserIDFound(t *testing.T) {
	useTempWorkDir(t)

	if err := SaveUserIDs(UserID{Name: "www-data", UID: "33", GID: "33"}); err != nil {
		t.Fatalf("SaveUserIDs returned error: %v", err)
	}

	id, err := GetUserID("www-data")
	if err != nil {
		t.Fatalf("GetUserID returned error: %v", err)
	}
	if id == nil {
		t.Fatal("expected user ID")
	}
	if id.Name != "www-data" {
		t.Fatalf("unexpected name: %q", id.Name)
	}
}

func TestGetUserIDNotFound(t *testing.T) {
	useTempWorkDir(t)

	_, err := GetUserID("missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if got, want := err.Error(), `user "missing" not found`; got != want {
		t.Fatalf("unexpected error: got %q want %q", got, want)
	}
}

func TestLoadUserIDsInvalidJSON(t *testing.T) {
	dir := useTempWorkDir(t)

	path := filepath.Join(dir, UserIDsFile)
	if err := os.WriteFile(path, []byte("{invalid"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	_, err := LoadUserIDs()
	if err == nil {
		t.Fatal("expected error")
	}
}
