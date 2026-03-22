package protocol

import (
	"strings"
)

// Nextcloud describes one Nextcloud application instance.
//
// SystemFQDN is the canonical hostname of the target system.
// It is used for authentication and trust decisions.
type Nextcloud struct {
	Name       string `json:"name"`
	Language   string `json:"language"`
	Region     string `json:"region"`
	HostName   string `json:"host_name"`
	DomainName string `json:"domain_name"`
	Version    string `json:"version"`
	PhpVersion string `json:"php_version"`
	SystemFQDN string `json:"system_fqdn"`
	Subdir     string `json:"subdir,omitempty"`
	Password   string `json:"password"`
	InstanceID string `json:"instance_id"`
	Salt       string `json:"salt"`
	Secret     string `json:"secret"`
	AdminEmail string `json:"admin_email,omitempty"`
}

// NextConfig describes one Nextcloud configuration entry.
type NextConfig struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// FQDN returns the fully qualified domain name of the Nextcloud instance.
func (nc *Nextcloud) FQDN() string {
	if nc == nil {
		return ""
	}
	return MakeFQDN(nc.HostName, nc.DomainName)
}

// SlashSubdir returns Subdir normalized with a leading slash.
// It returns an empty string if no subdir is configured.
func (nc *Nextcloud) SlashSubdir() string {
	if nc == nil {
		return ""
	}
	sub := strings.Trim(nc.Subdir, " /")
	if sub == "" {
		return ""
	}
	return "/" + sub
}

// NextcloudApp contains Nextcloud-specific application data.
type NextcloudApp struct {
	Nextcloud   *Nextcloud    `json:"nextcloud,omitempty"`
	NextConfigs []*NextConfig `json:"next_configs,omitempty"`
}

// HasNextcloudApp reports whether the request contains Nextcloud-related data.
func (req *Request) HasNextcloudApp() bool {
	if req == nil {
		return false
	}
	return req.Nextcloud != nil || len(req.NextConfigs) > 0
}
