package protocol

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

func (req *Request) HasUserList() bool {
	if req == nil {
		return false
	}
	return len(req.Users) > 0
}
