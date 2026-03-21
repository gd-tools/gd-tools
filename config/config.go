package config

import (
	"crypto"
	"crypto/tls"
	"net"
	"os"

	"github.com/gd-tools/gd-tools/utils"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/urfave/cli/v2"
)

const (
	DefaultTimeout = 10
)

// Config contains the persistent server configuration plus runtime-only helpers.
type Config struct {
	// Persistent production server configuration.
	Server

	// Concrete baseline for this server (Ubuntu version etc.).
	Baseline *Baseline

	// Products contains all downloadable assets.
	Products []Product

	// Common runtime flags.
	Verbose bool
	Force   bool
	Delete  bool
	SkipDNS bool
	SkipMX  bool

	// Logger used for Info, Debug and Error output.
	logger Logger

	// Connection to the production server.
	Port    string
	Conn    *tls.Conn
	Timeout int

	// File handling helpers.
	mkdirAll func(path string, perm os.FileMode) error
	setenv   func(key, value string) error
	unsetenv func(key string) error
	loadFile func(name string) ([]byte, error)
	loadJSON func(name string, v any) error
	saveFile func(name string, data []byte) error
	saveJSON func(name string, v any) error

	// Template rendering from the gdt binary.
	render     func(name string, data any) ([]byte, error)
	renderJSON func(name string, data any, v any) error {
	renderSQL  func(name string, data any) ([]string, error) {
	renderList func(name, comment string, data any) ([]string, error) {

	// Function calls with side effects.
	lookupIP   func(host string) ([]net.IP, error)
	dhParams   func(bits int) ([]byte, error)
	rsaKeyPair func(fqdn string) ([]byte, []byte, error)
	runShell   func(commands []string) ([]byte, error)
	runCommand func(name string, args ...string) ([]byte, error)

	// ACME helpers.
	getPrivateKey            func(path string) (crypto.PrivateKey, error)
	getCloudflareCertificate func(domains []string, email string, key crypto.PrivateKey) (*certificate.Resource, error)
	getHetznerCertificate    func(domains []string, email string, key crypto.PrivateKey) (*certificate.Resource, error)
	getIonosCertificate      func(domains []string, email string, key crypto.PrivateKey) (*certificate.Resource, error)
}

// ReadConfig loads and initializes a server configuration.
func ReadConfig(c *cli.Context) (*Config, error) {
	cfg := &Config{}

	// The persistent server configuration must exist, created by "gdt setup".
	err := utils.LoadJSON(utils.ConfigFile, &cfg.Server)
	if err != nil {
		return nil, err
	}

	// The baseline for this particular server.
	cfg.Baseline, err = LoadBaseline(cfg.BaselineName)
	if err != nil {
		return nil, err
	}

	if c != nil {
		cfg.Verbose = c.Bool("verbose")
		cfg.Force = c.Bool("force")
		cfg.Delete = c.Bool("delete")
		cfg.SkipDNS = c.Bool("skip-dns")
		cfg.SkipMX = c.Bool("skip-mx")
		cfg.Port = c.String("port")
		// TODO agent.SetAgentPort(cfg.Port)
	}

	return cfg, nil
}

// Save writes the persistent server configuration.
func (cfg *Config) Save() error {
	return cfg.SaveJSON(utils.ConfigFile, &cfg.Server)
}
