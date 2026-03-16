package config

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sort"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/platform"
	"github.com/gd-tools/gd-tools/utils"
	"github.com/urfave/cli/v2"
)

const (
	ConfigFile     = "config.json"
	DefaultTimeout = 10

	RunPrefix   = "[run]"
	DebugPrefix = "[dbg] ##########"
)

type Config struct {
	// Runtime environment: ubuntu + php + dovecot
	Platform *platform.Platform `json:"-"`

	// Concrete baseline for this server
	BaselineName string             `json:"baseline"`
	Baseline     *platform.Baseline `json:"-"`

	// Flags that everyone uses
	Verbose bool `json:"-"`
	Force   bool `json:"-"`
	Delete  bool `json:"-"`
	SkipDNS bool `json:"-"`
	SkipMX  bool `json:"-"`

	// Connection to the production server
	Port    string    `json:"-"`
	Conn    *tls.Conn `json:"-"`
	Timeout int       `json:"-"`

	// Useful application access helpers
	Mailer   *Mailer  `json:"-"`
	DKIMs    []string `json:"-"`
	Postfix  *Postfix `json:"-"`
	Password string   `json:"-"`
	RootDir  string   `json:"-"`
	BaseDir  string   `json:"-"`
	LogsDir  string   `json:"-"`
	CertDir  string   `json:"-"`

	// The user facing environment
	utils.Basics

	// Some server parameters
	HostName   string         `json:"host_name"`   // host part of FQDN
	DomainName string         `json:"domain_name"` // domain part of FQDN
	IPv4Addr   string         `json:"ipv4_addr"`   // must be set manually
	IPv6Addr   string         `json:"ipv6_addr"`   // must be set manually
	SwapSize   string         `json:"swap_size"`   // e.g. 500M or 2G, or 0
	Mounts     []*agent.Mount `json:"mounts"`      // mounted filesystem (which can grow)
	Firewall   []string       `json:"firewall"`    // ufw ports to open, e.g. 22/tcp

	// List of installed host names (to avoid name collisions)
	UsedFQDNs []string `json:"used_fqdns"`

	// System software versions (for lack of a better place)
	Roundcube string `json:"roundcube"`

	// External credentials
	Spambarrier string `json:"spambarrier,omitempty"` // spambarrier.de API key (see https://www.spambarrier.de/)
	UbuntuPro   string `json:"ubuntu_pro,omitempty"`  // Ubuntu Pro Key (see https://ubuntu.com/pro)

	HetznerToken    string `json:"hetzner_token,omitempty"`    // Token for Hetzner Cloud API
	IonosToken      string `json:"ionos_token,omitempty"`      // Token for IONOS Cloud API
	CloudflareToken string `json:"cloudflare_token,omitempty"` // Token for Cloudflare Cloud API
	// add more provider token here
}

var (
	FlagVerbose = &cli.BoolFlag{
		Name:  "verbose",
		Usage: "print extended debugging information",
	}
	FlagForce = &cli.BoolFlag{
		Name:  "force",
		Usage: "enforce action even if not necessary",
	}
	FlagSkipDNS = &cli.BoolFlag{
		Name:  "skip-dns",
		Usage: "do not execute DNS installs or updates",
	}
	FlagSkipMX = &cli.BoolFlag{
		Name:  "skip-mx",
		Usage: "do not touch the domain MX record",
	}
	FlagDelete = &cli.BoolFlag{
		Name:  "delete",
		Usage: "delete the given resource - may the --force be with you",
	}
	FlagPort = &cli.StringFlag{
		Name:  "port",
		Usage: "TCP port for dev/prod communication",
		Value: agent.DefaultPort,
	}
)

func (cfg *Config) Locale() string {
	return cfg.Language + "_" + cfg.Region
}

func (cfg *Config) FQDN() string {
	return cfg.HostName + "." + cfg.DomainName
}

func (cfg *Config) DHParamsPath() string {
	return assets.GetEtcDir("apache2", utils.DHParamsFile)
}

func (cfg *Config) RootUser() string {
	return "root@" + cfg.FQDN()
}

func (cfg *Config) FQDNdot() string {
	return cfg.HostName + "." + cfg.DomainName + "."
}

func (cfg *Config) DotFQDN() string {
	return "." + cfg.HostName + "." + cfg.DomainName
}

func (cfg *Config) RsyncFlags() string {
	if cfg.Verbose {
		return "-avz"
	}
	return "-avzq"
}

func ReadConfig(c *cli.Context, pf *platform.Platform) (*Config, error) {
	content, err := os.ReadFile(ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("missing %s - are we in the correct dir?", ConfigFile)
		}
		return nil, fmt.Errorf("failed to read %s: %w", ConfigFile, err)
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", ConfigFile, err)
	}

	// Load the server platform (necessary for recovery: never use "latest")
	if pf == nil {
		pf, err = platform.LoadPlatform()
		if err != nil {
			return nil, err
		}
	}
	cfg.Platform = pf

	// The baseline for this particular server
	cfg.Baseline, err = pf.GetBaseline(cfg.BaselineName)
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
		agent.SetAgentPort(cfg.Port)
	}

	if cfg.HostName == "" {
		return nil, fmt.Errorf("missing HostName in config.json")
	}
	if cfg.DomainName == "" {
		return nil, fmt.Errorf("missing DomainName in config.json")
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}

	// Ensure identity values - they are used in various contexts
	id, err := utils.FetchIdentity()
	if err != nil {
		return nil, err
	}
	if cfg.Company == "" {
		cfg.Company = id.Company
	}
	if cfg.Domain == "" {
		cfg.Domain = id.Domain
	}
	if cfg.SysAdmin == "" {
		cfg.SysAdmin = id.SysAdmin
	}
	if cfg.TimeZone == "" {
		cfg.TimeZone = id.TimeZone
	}
	if cfg.Language == "" {
		cfg.Language = id.Language
	}
	if cfg.Region == "" {
		cfg.Region = id.Region
	}
	if cfg.RegTTL == 0 {
		cfg.RegTTL = id.RegTTL
	}
	if cfg.DMARC == "" {
		cfg.DMARC = id.DMARC
	}

	// Ensure working environment
	if err := os.MkdirAll(ACME_Cert_Dir, 0755); err != nil {
		return nil, err
	}

	if _, _, err := utils.GetRSAKeyPair(cfg.FQDN()); err != nil {
		return nil, err
	}

	cfg.AddFirewall("22/tcp")
	cfg.AddFirewall(agent.GetAgentPort() + "/tcp")

	ips, err := net.LookupIP(cfg.FQDN())
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		cfg.Debugf("IP Address: %s", ip.String())
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			if cfg.IPv4Addr != ip.String() {
				if cfg.IPv4Addr == "" {
					cfg.IPv4Addr = ip.String()
					cfg.Sayf("IPv4Addr defaults to %s", cfg.IPv4Addr)
				} else {
					cfg.Debugf("IPv4Addr was changed to %s", cfg.IPv4Addr)
				}
			}
			break
		}
	}

	for _, ip := range ips {
		if ip.To16() != nil && ip.To4() == nil {
			if cfg.IPv6Addr != ip.String() {
				if cfg.IPv6Addr == "" {
					cfg.IPv6Addr = ip.String()
					cfg.Sayf("IPv6Addr defaults to %s", cfg.IPv6Addr)
				} else {
					cfg.Debugf("IPv6Addr was changed to %s", cfg.IPv6Addr)
				}
			}
			break
		}
	}

	if err := cfg.Save(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) Close() {
	if cfg != nil && cfg.Conn != nil {
		cfg.Conn.Close()
		cfg.Conn = nil
	}
}

func (cfg *Config) NewRequest() *agent.Request {
	return &agent.Request{
		Version:  agent.ProtocolVersion,
		Conn:     cfg.Conn,
		Verbose:  cfg.Verbose,
		Platform: cfg.Platform,
		Language: cfg.Language,
		Region:   cfg.Region,
	}
}

func (cfg *Config) Save() error {
	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", ConfigFile, err)
	}

	existing, err := os.ReadFile(ConfigFile)
	if err == nil && bytes.Equal(existing, content) {
		return nil
	}

	if err := os.WriteFile(ConfigFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", ConfigFile, err)
	}
	cfg.Say("config.json has been updated - please check")

	return nil
}

func (cfg *Config) CheckFQDN(fqdn string, add bool) error {
	for _, check := range cfg.UsedFQDNs {
		if fqdn == check {
			if !add {
				return fmt.Errorf("the name '%s' is already in use", fqdn)
			}
		}
	}

	if add {
		cfg.UsedFQDNs = append(cfg.UsedFQDNs, fqdn)
		sort.Strings(cfg.UsedFQDNs)
		if err := cfg.Save(); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) Say(args ...any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			for _, line := range v {
				fmt.Println(RunPrefix, line)
			}
		default:
			fmt.Println(RunPrefix, fmt.Sprint(arg))
		}
	}
}

func (cfg *Config) Sayf(format string, args ...any) {
	fmt.Println(RunPrefix, fmt.Sprintf(format, args...))
}

func (cfg *Config) Debug(args ...any) {
	if !cfg.Verbose {
		return
	}

	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			for _, line := range v {
				fmt.Println(DebugPrefix, line)
			}
		default:
			fmt.Println(DebugPrefix, fmt.Sprint(arg))
		}
	}
}

func (cfg *Config) Debugf(format string, args ...any) {
	if !cfg.Verbose {
		return
	}

	fmt.Println(DebugPrefix, fmt.Sprintf(format, args...))
}

func (cfg *Config) AddFirewall(port string) {
	for _, open := range cfg.Firewall {
		if port == open {
			return
		}
	}

	cfg.Firewall = append(cfg.Firewall, port)
	sort.Strings(cfg.Firewall)
}
