package config

import (
	"github.com/gd-tools/gd-tools/utils"
)

// Server is the persistent user-facing server configuration.
// It only contains values that belong to config.json and must not
// contain runtime-only helpers, flags, or active connections.
//
// WARNING: this struct may contain sensitive data (API tokens).
// It must never be deployed to production systems.
type Server struct {
	// Baseline selected for this server.
	BaselineName string `json:"baseline"`

	// User-facing identity / ownership data.
	utils.Identity

	// Host identity.
	HostName   string `json:"host_name"`
	DomainName string `json:"domain_name"`

	// Network addresses.
	IPv4Addr string `json:"ipv4_addr"`
	IPv6Addr string `json:"ipv6_addr"`

	// Basic host settings.
	SwapSize string    `json:"swap_size,omitempty"`
	Mounts   MountList `json:"mounts,omitempty"`
	Firewall []string  `json:"firewall,omitempty"`

	// Installed names to avoid collisions.
	UsedFQDNs []string `json:"used_fqdns"`

	// Application / package versions (system modules).
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

func (srv *Server) FQDNdot() string {
	fqdn := srv.FQDN()
	if fqdn == "" {
		return ""
	}
	return fqdn + "."
}

func (srv *Server) DotFQDN() string {
	fqdn := srv.FQDN()
	if fqdn == "" {
		return ""
	}
	return "." + fqdn
}

func (srv *Server) RootUser() string {
	if srv == nil {
		return ""
	}
	fqdn := srv.FQDN()
	if fqdn == "" {
		return "root"
	}
	return "root@" + fqdn
}

func (srv *Server) Locale() string {
	if srv == nil {
		return ""
	}
	if srv.Language == "" {
		return srv.Region
	}
	if srv.Region == "" {
		return srv.Language
	}
	return srv.Language + "_" + srv.Region
}
