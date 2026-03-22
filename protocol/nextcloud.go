package protocol

import (
	"strings"
)

type Nextcloud struct {
	Name       string `json:"name"`
	Language   string `json:"language"`
	Region     string `json:"region"`
	HostName   string `json:"host_name"`
	DomainName string `json:"domain_name"`
	Version    string `json:"version"`
	PhpVersion string `json:"php_version"`
	ServerFQDN string `json:"server_fqdn"`
	Subdir     string `json:"subdir"`
	Password   string `json:"password"`
	InstanceID string `json:"instance_id"`
	Salt       string `json:"salt"`
	Secret     string `json:"secret"`
	AdminEmail string `json:"admin_email"`

	// Include download for conveniance
	Download *Download `json:"-"`
}

type NextConfig struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (nc *Nextcloud) FQDN() string {
	if nc == nil {
		return ""
	}
	return MakeFQDN(nc.HostName, nc.DomainName)
}

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

type NextcloudApp struct {
	Nextcloud   *Nextcloud    `json:"nextcloud"`
	NextConfigs []*NextConfig `json:"nextconfigs"`
}

func (req *Request) HasNextcloudApp() bool {
	if req == nil {
		return false
	}
	return req.Nextcloud != nil || len(req.NextConfigs) > 0
}
