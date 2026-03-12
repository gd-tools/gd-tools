package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gd-tools/gd-tools/assets"
	"github.com/gd-tools/gd-tools/agent"
)

const (
	OCISName = "ocis"
	OCISFile = OCISName + ".json"
)

type OCIS struct {
	HostName   string `json:"host_name"`
	DomainName string `json:"domain_name"`
	Version    string `json:"version"`
	AdminName  string `json:"admin_name"`
	AdminEmail string `json:"admin_email"`
	Password   string `json:"password"`
	Language   string `json:"language"`
	LogLevel   string `json:"log_level"`
}

func (oc *OCIS) FQDN() string {
	return oc.HostName + "." + oc.DomainName
}

func (oc *OCIS) ExecPath() string {
	return assets.GetBinDir("ocis")
}

func (oc *OCIS) RootDir(paths ...string) string {
	rootDir := assets.GetToolsDir("data", "ocis")
	if len(paths) == 0 {
		return rootDir
	}
	return filepath.Join(append([]string{rootDir}, paths...)...)
}

func (oc *OCIS) BaseDir(paths ...string) string {
	baseDir := oc.RootDir(".ocis")
	if len(paths) == 0 {
		return baseDir
	}
	return filepath.Join(append([]string{baseDir}, paths...)...)
}

func (oc *OCIS) ConfigDir() string {
	return oc.BaseDir("config")
}

func (oc *OCIS) ClientDir(paths ...string) string {
	clientDir := filepath.Join(oc.RootDir(), "client")
	if len(paths) == 0 {
		return clientDir
	}
	return filepath.Join(append([]string{clientDir}, paths...)...)
}

func (oc *OCIS) CertDir() string {
	return assets.GetToolsDir("data", "certs", oc.FQDN())
}

func (oc *OCIS) LogsDir() string {
	return assets.GetToolsDir("logs", "ocis")
}

func LoadOCIS() (*OCIS, error) {
	content, err := os.ReadFile(OCISFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // OCIS is not installed
		}
		return nil, fmt.Errorf("failed to read %s: %w", OCISFile, err)
	}

	var oc OCIS
	if err := json.Unmarshal(content, &oc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", OCISFile, err)
	}

	return &oc, nil
}

func (oc *OCIS) Save() error {
	content, err := json.MarshalIndent(oc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", OCISFile, err)
	}

	existing, err := os.ReadFile(OCISFile)
	if err == nil && bytes.Equal(existing, content) {
		return nil
	}

	if err := os.WriteFile(OCISFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", OCISFile, err)
	}

	return nil
}

func (cfg *Config) DeployOCIS(oc *OCIS) error {
	if oc == nil {
		return fmt.Errorf("missing OCIS pointer")
	}
	cfg.Debug("Enter config/ocis.go")

	if err := cfg.OCISDownload(oc); err != nil {
		return err
	}

	if err := cfg.OCISUser(oc); err != nil {
		return err
	}

	if err := cfg.EnsureIDMSelfSignedCert(); err != nil {
		return err
	}

	if err := cfg.OCISConfig(oc); err != nil {
		return err
	}

	if err := cfg.OCISService(oc); err != nil {
		return err
	}

	if err := cfg.EnsureCertificate(oc.FQDN()); err != nil {
		return err
	}

	if err := cfg.OCISVhost(oc); err != nil {
		return err
	}

	if status, err := cfg.SetCNAME(oc.DomainName, oc.HostName); err != nil {
		return err
	} else if status != "" {
		cfg.Say(status)
	}

	cfg.Debug("Leave config/ocis.go")
	return nil
}

func (cfg *Config) OCISDownload(oc *OCIS) error {
	req := cfg.NewRequest()

	_, ocRel, err := cfg.Catalog.Get("ocis", oc.Version)
	if err != nil {
		return err
	}
	if ocRel.Download.Binary == "" {
		return fmt.Errorf("missing Binary in OCIS download")
	}

	req.Downloads = append(req.Downloads, &ocRel.Download)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) OCISUser(oc *OCIS) error {
	req := cfg.NewRequest()

	ocisUser := agent.User{
		Name:    "ocis",
		Comment: "OCIS Server User",
		System:  true,
		HomeDir: oc.RootDir(),
		Shell:   "/bin/bash",
	}
	req.Users = append(req.Users, &ocisUser)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) OCISConfig(oc *OCIS) error {
	req := cfg.NewRequest()

	rootMkdir := agent.File{
		Task:  "mkdir",
		Path:  oc.RootDir(),
		Mode:  "0750",
		User:  "ocis",
		Group: "ocis",
	}
	req.AddFile(&rootMkdir)

	baseMkdir := agent.File{
		Task:  "mkdir",
		Path:  oc.BaseDir(),
		Mode:  "0750",
		User:  "ocis",
		Group: "ocis",
	}
	req.AddFile(&baseMkdir)

	configMkdir := agent.File{
		Task:  "mkdir",
		Path:  oc.ConfigDir(),
		Mode:  "0750",
		User:  "ocis",
		Group: "ocis",
	}
	req.AddFile(&configMkdir)

	envTmpl, err := assets.Render("ocis/ocis.env", oc)
	if err != nil {
		return err
	}
	envPath := filepath.Join(oc.ConfigDir(), "ocis.env")
	envFile := agent.File{
		Task:    "write",
		Path:    envPath,
		Content: envTmpl,
		Mode:    "0640",
		User:    "ocis",
		Group:   "ocis",
	}
	req.AddFile(&envFile)

	idmMkdir := agent.File{
		Task:  "mkdir",
		Path:  oc.BaseDir("idm"),
		Mode:  "0750",
		User:  "ocis",
		Group: "ocis",
	}
	req.AddFile(&idmMkdir)

	crtName := "ocis-ldap.crt"
	crtData, err := os.ReadFile(crtName)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", crtName, err)
	}
	crtFile := agent.File{
		Task:    "write",
		Path:    oc.BaseDir("idm", "ldap.crt"),
		Content: crtData,
		Mode:    "0640",
		User:    "ocis",
		Group:   "ocis",
	}
	req.AddFile(&crtFile)

	keyName := "ocis-ldap.key"
	keyData, err := os.ReadFile(keyName)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", keyName, err)
	}
	keyFile := agent.File{
		Task:    "write",
		Path:    oc.BaseDir("idm", "ldap.key"),
		Content: keyData,
		Mode:    "0600",
		User:    "ocis",
		Group:   "ocis",
	}
	req.AddFile(&keyFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) OCISService(oc *OCIS) error {
	req := cfg.NewRequest()

	svcTmpl, err := assets.Render("ocis/ocis.service", oc)
	if err != nil {
		return err
	}

	file := agent.File{
		Task:    "write",
		Path:    assets.GetEtcDir("systemd", "system", "ocis.service"),
		Content: svcTmpl,
		Mode:    "0644",
		Service: "ocis",
	}
	req.AddFile(&file)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) OCISVhost(oc *OCIS) error {
	req := cfg.NewRequest()

	logsMkdir := agent.File{
		Task:  "mkdir",
		Path:  oc.LogsDir(),
		Mode:  "0750",
		User:  "www-data",
		Group: "www-data",
	}
	req.AddFile(&logsMkdir)

	vhostTmpl, err := assets.Render("ocis/vhost.conf", oc)
	if err != nil {
		return err
	}
	vhostName := fmt.Sprintf("ocis-%s.conf", oc.FQDN())
	vhostPath := assets.GetApacheEtcDir("sites-available", vhostName)
	vhostFile := agent.File{
		Task:    "write",
		Path:    vhostPath,
		Content: vhostTmpl,
		Service: "apache2",
	}
	req.AddFile(&vhostFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}
