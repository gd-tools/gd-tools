package protocol

import "testing"

func TestRequestAddUser(t *testing.T) {
	req := &Request{}

	user := &User{Name: "nginx"}
	req.AddUser(user)

	if len(req.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(req.Users))
	}
	if req.Users[0].Name != "nginx" {
		t.Fatalf("unexpected user name: %q", req.Users[0].Name)
	}
}

func TestRequestAddUserDedup(t *testing.T) {
	req := &Request{}

	req.AddUser(&User{Name: "nginx"})
	req.AddUser(&User{Name: "nginx"})

	if len(req.Users) != 1 {
		t.Fatalf("expected deduplicated user list, got %d", len(req.Users))
	}
}

func TestRequestAddUserTrim(t *testing.T) {
	req := &Request{}

	req.AddUser(&User{Name: "  nginx  "})

	if len(req.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(req.Users))
	}
	if req.Users[0].Name != "nginx" {
		t.Fatalf("expected trimmed name, got %q", req.Users[0].Name)
	}
}

func TestRequestAddUserIgnoreEmpty(t *testing.T) {
	req := &Request{}

	req.AddUser(&User{Name: ""})

	if len(req.Users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(req.Users))
	}
}

func TestRequestAddUserNilReceiver(t *testing.T) {
	var req *Request
	req.AddUser(&User{Name: "nginx"})
}

func TestRequestHasUserList(t *testing.T) {
	if (&Request{}).HasUserList() {
		t.Fatalf("expected false for empty request")
	}

	req := &Request{}
	req.AddUser(&User{Name: "nginx"})

	if !req.HasUserList() {
		t.Fatalf("expected true after adding user")
	}
}
