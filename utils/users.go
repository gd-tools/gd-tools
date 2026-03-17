package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

const (
	UserIDsFile = "user_ids.json"
)

type User struct {
	Name    string   `json:"name"`
	Comment string   `json:"comment,omitempty"`
	System  bool     `json:"system,omitempty"`
	Shell   string   `json:"shell,omitempty"`
	HomeDir string   `json:"home_dir,omitempty"`
	Groups  []string `json:"groups,omitempty"`
}

type UserID struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
	GID  string `json:"gid"`
}

func LoadUserIDs() ([]UserID, error) {
	var list []UserID

	content, err := os.ReadFile(UserIDsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return list, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", UserIDsFile, err)
	}

	if err := json.Unmarshal(content, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", UserIDsFile, err)
	}

	return list, nil
}

func GetUserID(name string) (*UserID, error) {
	list, err := LoadUserIDs()
	if err != nil {
		return nil, err
	}

	for i := range list {
		if list[i].Name == name {
			return &list[i], nil
		}
	}

	return nil, fmt.Errorf("user %q not found", name)
}

func SaveUserIDs(newUser UserID) error {
	oldList, err := LoadUserIDs()
	if err != nil {
		return err
	}

	userMap := make(map[string]UserID)
	for _, u := range oldList {
		userMap[u.Name] = u
	}
	userMap[newUser.Name] = newUser

	var names []string
	for name := range userMap {
		names = append(names, name)
	}
	sort.Strings(names)

	var newList []UserID
	for _, name := range names {
		newList = append(newList, userMap[name])
	}

	content, err := json.MarshalIndent(newList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", UserIDsFile, err)
	}
	if err := os.WriteFile(UserIDsFile, content, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", UserIDsFile, err)
	}

	return nil
}
