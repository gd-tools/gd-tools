package protocol

import (
	"encoding/json"
)

const (
	CurrentVersion = 1
)

// Request is the shared protocol message between dev and prod.
//
// Keep this model explicit and small.
// New application support should prefer composition of existing protocol
// primitives such as files, services, packages, downloads, and database.
//
// Add a new top-level request field only if prod requires distinct semantics
// that cannot be expressed cleanly through existing primitives.
type Request struct {
	Version int `json:"version,omitempty"`

	// Basics: mTLS communication testing
	Hello

	// Basics: bootstrap
	Bootstrap

	// Basics: apt-get Ubuntu packages
	PackageList

	// Basics: user management
	UserList

	// Basics: mounts and filesystem layout
	Filesystem

	// Basics: database startup
	Database

	// General: file handling
	FileList

	// General: downloads from the internet
	DownloadList

	// General: services to (re)start
	ServiceList

	// General: firewall rules
	FirewallList

	// App: Nextcloud
	NextcloudApp

	// App: RustDesk
	RustDeskApp

	// More apps here if they require common knowledge.
}

func (req *Request) String() string {
	if req == nil {
		return ""
	}
	content, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return "failed to marshal Request"
	}
	return string(content)
}

func MakeFQDN(host, domain string) string {
	if host == "" {
		return domain
	}
	if domain == "" {
		return host
	}
	return host + "." + domain
}
