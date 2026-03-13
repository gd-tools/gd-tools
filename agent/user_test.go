package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestUsersTest(t *testing.T) {

	req := &Request{}

	if UsersTest(req) {
		t.Fatal("expected false with no users")
	}

	req.Users = []*User{{Name: "test"}}

	if !UsersTest(req) {
		t.Fatal("expected true when users exist")
	}
}

func TestUserString(t *testing.T) {

	u := &User{
		Name:    "alice",
		Comment: "Alice Example",
		Shell:   "/bin/bash",
		HomeDir: "/home/alice",
		System:  false,
	}

	s := u.String()

	if s == "" {
		t.Fatal("string output empty")
	}
}

func TestSaveAndLoadUserIDs(t *testing.T) {

	dir := t.TempDir()

	old := UserIDsName
	UserIDsName = filepath.Join(dir, "user_ids.json")
	defer func() { UserIDsName = old }()

	u := UserID{
		Name: "alice",
		UID:  "1000",
		GID:  "1000",
	}

	if err := SaveUserIDs(u); err != nil {
		t.Fatal(err)
	}

	list, err := LoadUserIDs()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 user, got %d", len(list))
	}

	if list[0].Name != "alice" {
		t.Fatal("user mismatch")
	}
}

func TestSaveUserIDsMerge(t *testing.T) {

	dir := t.TempDir()

	old := UserIDsName
	UserIDsName = filepath.Join(dir, "user_ids.json")
	defer func() { UserIDsName = old }()

	u1 := UserID{Name: "bob", UID: "1001", GID: "1001"}
	u2 := UserID{Name: "alice", UID: "1000", GID: "1000"}

	if err := SaveUserIDs(u1); err != nil {
		t.Fatal(err)
	}

	if err := SaveUserIDs(u2); err != nil {
		t.Fatal(err)
	}

	list, err := LoadUserIDs()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 users, got %d", len(list))
	}

	if list[0].Name != "alice" {
		t.Fatal("users not sorted")
	}
}

func TestGetUserID(t *testing.T) {

	dir := t.TempDir()

	old := UserIDsName
	UserIDsName = filepath.Join(dir, "user_ids.json")
	defer func() { UserIDsName = old }()

	u := UserID{Name: "carol", UID: "1002", GID: "1002"}

	if err := SaveUserIDs(u); err != nil {
		t.Fatal(err)
	}

	id, err := GetUserID("carol")
	if err != nil {
		t.Fatal(err)
	}

	if id.UID != "1002" {
		t.Fatal("uid mismatch")
	}
}

func TestLoadUserIDsMissingFile(t *testing.T) {

	dir := t.TempDir()

	old := UserIDsName
	UserIDsName = filepath.Join(dir, "user_ids.json")
	defer func() { UserIDsName = old }()

	list, err := LoadUserIDs()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 0 {
		t.Fatal("expected empty list")
	}
}

func TestUserIDsJSONFormat(t *testing.T) {

	dir := t.TempDir()

	old := UserIDsName
	UserIDsName = filepath.Join(dir, "user_ids.json")
	defer func() { UserIDsName = old }()

	u := UserID{Name: "dave", UID: "1003", GID: "1003"}

	if err := SaveUserIDs(u); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(UserIDsName)
	if err != nil {
		t.Fatal(err)
	}

	var parsed []UserID

	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal("invalid json written")
	}
}
