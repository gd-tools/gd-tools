package protocol

import "strings"

// User describes a system user.
type User struct {
	Name    string   `json:"name"`
	Comment string   `json:"comment,omitempty"`
	System  bool     `json:"system,omitempty"`
	Shell   string   `json:"shell,omitempty"`
	HomeDir string   `json:"home_dir,omitempty"`
	Groups  []string `json:"groups,omitempty"`
}

// UserID describes the return type for Response.
type UserID struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
	GID  string `json:"gid"`
}

type UserList struct {
	Users []*User `json:"users,omitempty"`
}

// AddUser adds a user if it is not nil and not already present (by name).
func (req *Request) AddUser(user *User) {
	if req == nil || user == nil {
		return
	}

	name := strings.TrimSpace(user.Name)
	if name == "" {
		return
	}
	user.Name = name

	for _, check := range req.Users {
		if check.Name == user.Name {
			return
		}
	}

	req.Users = append(req.Users, user)
}

// HasUserList reports whether the request contains at least one user entry.
func (req *Request) HasUserList() bool {
	if req == nil {
		return false
	}
	return len(req.Users) > 0
}
