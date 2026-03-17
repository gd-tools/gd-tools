package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

const (
	IdentityFile = "identity.json"

	DefaultCompany  = "Example GmbH"
	DefaultDomain   = "example.com"
	DefaultTimeZone = "Europe/Berlin"
	DefaultRegTTL   = 3600
	DefaultDMARC    = "v=DMARC1; p=quarantine; pct=100; adkim=s; aspf=s"
)

type Identity struct {
	Company  string `json:"company"`
	Domain   string `json:"domain"`
	SysAdmin string `json:"sys_admin"`
	HelpURL  string `json:"help_url"`
	TimeZone string `json:"time_zone"`
	Language string `json:"language"`
	Region   string `json:"region"`
	RegTTL   int    `json:"reg_ttl"`
	DMARC    string `json:"dmarc"`
}

func (id *Identity) Locale() string {
	return id.Language + "_" + id.Region
}

func (id *Identity) AdminMail() string {
	if id.SysAdmin != "" {
		return id.SysAdmin
	}
	return "admin@" + id.Domain
}

func (id *Identity) SupportURL() string {
	if id.HelpURL != "" {
		return id.HelpURL
	}
	return "https://support." + id.Domain + "/"
}

func (id *Identity) DMARCDomain() string {
	return "_dmarc." + id.Domain
}

func EnsureIdentity() (*Identity, error) {
	content, err := os.ReadFile(IdentityFile)
	if err != nil {
		if os.IsNotExist(err) {
			id := Identity{
				Company:  DefaultCompany,
				Domain:   DefaultDomain,
				SysAdmin: GetSysAdmin(),
				HelpURL:  "https://support." + DefaultDomain + "/",
				TimeZone: GetTimeZone(),
				Language: GetLanguage(),
				Region:   GetRegion(),
				RegTTL:   DefaultRegTTL,
				DMARC:    DefaultDMARC,
			}
			return &id, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", IdentityFile, err)
	}

	var id Identity
	if err := json.Unmarshal(content, &id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", IdentityFile, err)
	}

	if id.Domain == DefaultDomain {
		return nil, fmt.Errorf("%s has default values - please run 'gdt identity'", IdentityFile)
	}

	if id.DMARC == "" {
		id.DMARC = DefaultDMARC
	}

	return &id, nil
}

func FetchIdentity() (*Identity, error) {
	path := filepath.Join("..", IdentityFile)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("missing %s - are we in the correct dir?", path)
	}

	var id Identity
	if err := json.Unmarshal(content, &id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", path, err)
	}

	if id.DMARC == "" {
		id.DMARC = DefaultDMARC
	}

	return &id, nil
}

func (id *Identity) Save() error {
	content, err := json.MarshalIndent(id, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", IdentityFile, err)
	}

	existing, err := os.ReadFile(IdentityFile)
	if err == nil && bytes.Equal(existing, content) {
		return nil
	}

	tmp := IdentityFile + ".tmp"

	if err := os.WriteFile(tmp, content, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", tmp, err)
	}

	if err := os.Rename(tmp, IdentityFile); err != nil {
		return fmt.Errorf("failed to replace %s: %w", IdentityFile, err)
	}

	return nil
}

func GetTimeZone() string {
	target, err := os.Readlink("/etc/localtime")
	if err != nil {
		return DefaultTimeZone
	}

	rel, err := filepath.Rel("/usr/share/zoneinfo", target)
	if err != nil {
		return DefaultTimeZone
	}

	return rel
}

func GetSysAdmin() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		gitConfigPath := filepath.Join(homeDir, ".gitconfig")
		if gitConfig, err := ini.Load(gitConfigPath); err == nil {
			return gitConfig.Section("user").Key("email").String()
		}
	}

	return "admin@" + DefaultDomain
}
