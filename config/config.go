package config

type Config struct {
	model.Server

	Platform *platform.Platform `json:"-"`

	Baseline *platform.Baseline `json:"-"`

	Verbose bool `json:"-"`
	Force   bool `json:"-"`
	Delete  bool `json:"-"`
	SkipDNS bool `json:"-"`
	SkipMX  bool `json:"-"`

	Conn    *tls.Conn `json:"-"`
	Timeout int       `json:"-"`
}
