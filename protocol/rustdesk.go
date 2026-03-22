package protocol

import (
	"encoding/base64"
	"strings"
)

type RustDesk struct {
	HostName   string `json:"host_name,omitempty"`
	DomainName string `json:"domain_name,omitempty"`
	Version    string `json:"version"`

	PrivateB64 string `json:"private_b64,omitempty"`
	PublicKey  string `json:"public_key,omitempty"`

	// Include download for conveniance
	Download *Download `json:"-"`
}

func (rd *RustDesk) FQDN() string {
	if rd == nil {
		return ""
	}
	return MakeFQDN(rd.HostName, rd.DomainName)
}

func (rd *RustDesk) GetPrivate() ([]byte, error) {
	if rd.PrivateB64 == "" {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(rd.PrivateB64)
}

func (rd *RustDesk) SetPrivate(data []byte) {
	if len(data) == 0 {
		rd.PrivateB64 = ""
		return
	}
	rd.PrivateB64 = base64.StdEncoding.EncodeToString(data)
}

func (rd *RustDesk) GetPublic() string {
	if rd.PublicKey == "" {
		return ""
	}
	return strings.TrimSpace(rd.PublicKey)
}

func (rd *RustDesk) SetPublic(s string) {
	rd.PublicKey = strings.TrimSpace(s)
}

type RustDeskApp struct {
	RustDesk *RustDesk `json:"rust_desk,omitempty"`
}

func (req *Request) HasRustDeskApp() bool {
	if req == nil {
		return false
	}

	return req.RustDesk != nil
}
