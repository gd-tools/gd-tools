package protocol

import (
	"encoding/base64"
	"strings"
)

// RustDesk describes one RustDesk application instance.
type RustDesk struct {
	HostName   string `json:"host_name,omitempty"`
	DomainName string `json:"domain_name,omitempty"`
	Version    string `json:"version"`

	// PrivateB64 contains the private key as base64-encoded binary data.
	PrivateB64 string `json:"private_b64,omitempty"`
	PublicKey  string `json:"public_key,omitempty"`
}

// FQDN returns the fully qualified domain name of the RustDesk instance.
func (rd *RustDesk) FQDN() string {
	if rd == nil {
		return ""
	}
	return MakeFQDN(rd.HostName, rd.DomainName)
}

// GetPrivate decodes the base64-encoded private key.
func (rd *RustDesk) GetPrivate() ([]byte, error) {
	if rd == nil || rd.PrivateB64 == "" {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(rd.PrivateB64)
}

// SetPrivate stores the private key as base64-encoded binary data.
func (rd *RustDesk) SetPrivate(data []byte) {
	if rd == nil {
		return
	}
	if len(data) == 0 {
		rd.PrivateB64 = ""
		return
	}
	rd.PrivateB64 = base64.StdEncoding.EncodeToString(data)
}

// GetPublic returns the normalized public key.
func (rd *RustDesk) GetPublic() string {
	if rd == nil || rd.PublicKey == "" {
		return ""
	}
	return strings.TrimSpace(rd.PublicKey)
}

// SetPublic stores the normalized public key.
func (rd *RustDesk) SetPublic(s string) {
	if rd == nil {
		return
	}
	rd.PublicKey = strings.TrimSpace(s)
}

type RustDeskApp struct {
	RustDesk *RustDesk `json:"rust_desk,omitempty"`
}

// HasRustDeskApp reports whether the request contains RustDesk-related data.
func (req *Request) HasRustDeskApp() bool {
	if req == nil {
		return false
	}
	return req.RustDesk != nil
}
