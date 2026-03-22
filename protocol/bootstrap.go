package protocol

// Bootstrap contains basic server settings that are needed during initial setup.
type Bootstrap struct {
	FQDN     string `json:"fqdn,omitempty"`
	TimeZone string `json:"time_zone,omitempty"`
	SSHPort  int    `json:"ssh_port,omitempty"`
}

// HasBootstrap reports whether any bootstrap setting is present.
func (bstr Bootstrap) HasBootstrap() bool {
	return bstr.FQDN != "" || bstr.TimeZone != "" || bstr.SSHPort != 0
}
