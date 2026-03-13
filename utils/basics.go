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
	BasicsName = "basics"
	BasicsFile = BasicsName + ".json"

	DefaultTimeZone = "Europe/Berlin"
	DefaultCompany  = "My Company"
	DefaultDomain   = "example.com"
	DefaultRegTTL   = 3600
	DefaultDMARC    = "v=DMARC1; p=quarantine; pct=100; adkim=s; aspf=s"
)

type Basics struct {
	Company  string `json:"company"`
	Domain   string `json:"domain"`
	SysAdmin string `json:"sys_admin"`
	TimeZone string `json:"time_zone"`
	Language string `json:"language"`
	Region   string `json:"region"`
	RegTTL   int    `json:"reg_ttl"`
	HelpURL  string `json:"help_url"`
	DMARC    string `json:"dmarc"`
}

func (bsc *Basics) Locale() string {
	return bsc.Language + "_" + bsc.Region
}

func (bsc *Basics) AdminMail() string {
	if bsc.SysAdmin != "" {
		return bsc.SysAdmin
	}
	return "admin@" + bsc.Domain
}

func (bsc *Basics) SupportURL() string {
	if bsc.HelpURL != "" {
		return bsc.HelpURL
	}
	return "https://support." + bsc.Domain + "/"
}

func (bsc *Basics) DMARCDomain() string {
	return "_dmarc." + bsc.Domain
}

func EnsureBasics() (*Basics, error) {
	content, err := os.ReadFile(BasicsFile)
	if err != nil {
		if os.IsNotExist(err) {
			bsc := Basics{
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
			return &bsc, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", BasicsFile, err)
	}

	var bsc Basics
	if err := json.Unmarshal(content, &bsc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", BasicsFile, err)
	}

	if bsc.Domain == DefaultDomain {
		return nil, fmt.Errorf("%s has default values - please update", BasicsFile)
	}

	if bsc.DMARC == "" {
		bsc.DMARC = DefaultDMARC
	}

	return &bsc, nil
}

func GetBasics() (*Basics, error) {
	path := filepath.Join("..", BasicsFile)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("missing %s - are we in the correct dir?", path)
	}

	var bsc Basics
	if err := json.Unmarshal(content, &bsc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", path, err)
	}

	if bsc.DMARC == "" {
		bsc.DMARC = DefaultDMARC
	}

	return &bsc, nil
}

func (bsc *Basics) Save() error {
	content, err := json.MarshalIndent(bsc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", BasicsFile, err)
	}

	existing, err := os.ReadFile(BasicsFile)
	if err == nil && bytes.Equal(existing, content) {
		return nil
	}

	tmp := BasicsFile + ".tmp"

	if err := os.WriteFile(tmp, content, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", tmp, err)
	}

	if err := os.Rename(tmp, BasicsFile); err != nil {
		return fmt.Errorf("failed to replace %s: %w", BasicsFile, err)
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
