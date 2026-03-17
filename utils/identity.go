package utils

import (
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
	var id Identity

	err := LoadJSON(IdentityFile, &id)
	if err != nil {
		if os.IsNotExist(err) {
			id = Identity{
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
		return nil, err
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

	var id Identity
	err := LoadJSON(path, &id)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("missing %s - are we in the correct dir?", path)
		}
		return nil, err
	}

	if id.DMARC == "" {
		id.DMARC = DefaultDMARC
	}

	return &id, nil
}

func (id *Identity) Save() error {
	return SaveJSON(IdentityFile, id)
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
