package protocol

// Mount describes a storage volume managed by gd-tools.
type Mount struct {
	Provider string `json:"provider"`          // e.g. "hetzner"
	ID       string `json:"id"`                // e.g. "123456789"
	Path     string `json:"path"`              // e.g. "/var/gd-tools"
	FSType   string `json:"fs_type,omitempty"` // e.g. "ext4"
	Options  string `json:"options,omitempty"` // e.g. "defaults,noatime"
}

// Directory describes a directory and its access and ownership.
type Directory struct {
	Path  string `json:"path"`
	Mode  string `json:"mode,omitempty"`
	User  string `json:"user,omitempty"`
	Group string `json:"group,omitempty"`
}

// Filesystem describes mount points and directories.
type Filesystem struct {
	Mounts      []*Mount     `json:"mounts,omitempty"`
	Directories []*Directory `json:"directories,omitempty"`
}

// AddMount adds a mount point to the request.
func (req *Request) AddMount(mount *Mount) {
	if req == nil || mount == nil {
		return
	}
	req.Mounts = append(req.Mounts, mount)
}

// AddDirectory adds a directory creation task to the request.
func (req *Request) AddDirectory(directory *Directory) {
	if req == nil || directory == nil {
		return
	}
	req.Directories = append(req.Directories, directory)
}

// HasFilesystem reports whether the request contains mount or directory tasks.
func (req *Request) HasFilesystem() bool {
	if req == nil {
		return false
	}
	return len(req.Mounts) > 0 || len(req.Directories) > 0
}
