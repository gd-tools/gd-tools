package model

import (
	"github.com/gd-tools/gd-tools/utils"
)

const ConfigFile = "config.json"

// ServerModel is the persistent user-facing server configuration.
// It only contains values that belong to config.json and must not
// contain runtime-only helpers, flags, or active connections.
type Server struct {
	// Baseline selected for this server.
	BaselineName string `json:"baseline"`

	// User-facing identity / ownership data.
	utils.Identity

	// Host identity.
	HostName   string `json:"host_name"`   // host part of FQDN
	DomainName string `json:"domain_name"` // domain part of FQDN

	// Network addresses.
	IPv4Addr string `json:"ipv4_addr"` // must be set manually
	IPv6Addr string `json:"ipv6_addr"` // must be set manually

	// Basic host settings.
	SwapSize string    `json:"swap_size,omitempty"` // e.g. "500M", "2G", or "" for 0
	Mounts   MountList `json:"mounts,omitempty"`
	Firewall []string  `json:"firewall,omitempty"`

	// Installed names to avoid collisions.
	UsedFQDNs []string `json:"used_fqdns"`

	// Application / package versions that are part of the persisted model.
	Roundcube string `json:"roundcube"`

	// External credentials.
	Spambarrier string `json:"spambarrier,omitempty"`
	UbuntuPro   string `json:"ubuntu_pro,omitempty"`

	HetznerToken    string `json:"hetzner_token,omitempty"`
	IonosToken      string `json:"ionos_token,omitempty"`
	CloudflareToken string `json:"cloudflare_token,omitempty"`
}

func (srv *Server) FQDN() string {
	if srv == nil {
		return ""
	}
	if srv.HostName == "" {
		return srv.DomainName
	}
	if srv.DomainName == "" {
		return srv.HostName
	}
	return srv.HostName + "." + srv.DomainName
}
