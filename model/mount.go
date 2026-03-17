package model

// Mount describes a storage volume managed by gd-tools.
type Mount struct {
	Provider string `json:"provider"`          // e.g. "hetzner"
	ID       string `json:"id"`                // e.g. "123456789"
	Path     string `json:"path"`              // e.g. "/var/gd-tools"
	FSType   string `json:"fs_type,omitempty"` // e.g. "ext4"
	Options  string `json:"options,omitempty"` // e.g. "defaults,noatime"
}

// MountList is a list of mounts.
type MountList []Mount

// HasMount reports whether a mount for the given directory exists.
func (mounts MountList) HasMount(path string) bool {
	for _, mount := range mounts {
		if mount.Path == path {
			return true
		}
	}
	return false
}
