package agent

import (
	"fmt"
	"path/filepath"

	"github.com/railduino/gd-tools/php"
)

const (
	PrefixAdmidio      = "ad"
	PrefixBookStack    = "bk"
	PrefixDovecot      = "dc" // system module
	PrefixFirefly      = "fi"
	PrefixImmich       = "im"
	PrefixMediaWiki    = "mw"
	PrefixMinecraft    = "mc"
	PrefixNextcloud    = "nc"
	PrefixOCIS         = "oc"
	PrefixOpenDKIM     = "od" // system module
	PrefixPaperlessNGX = "pl"
	PrefixPostfix      = "pf" // system module
	PrefixRoundcube    = "rc"
	PrefixRustDesk     = "rd"
	PrefixUptimeKuma   = "uk"
	PrefixVaultwarden  = "vw"
	PrefixWordPress    = "wp"
)

const (
	NamingAdmidio      = "admidio"
	NamingBookStack    = "bookstack"
	NamingFirefly      = "firefly"
	NamingImmich       = "immich"
	NamingMediaWiki    = "mediawiki"
	NamingMinecraft    = "minecraft"
	NamingNextcloud    = "nextcloud"
	NamingOCIS         = "ocis"
	NamingPaperlessNGX = "paperless"
	NamingRoundcube    = "roundcube"
	NamingRustDesk     = "rustdesk"
	NamingUptimeKuma   = "uptimekuma"
	NamingVaultwarden  = "vaultwarden"
	NamingWordPress    = "wordpress"
)

type NamingScheme struct {
	Short string
	Name  string
}

var namingSchemes = map[string]*NamingScheme{
	NamingAdmidio:      {Short: PrefixAdmidio, Name: NamingAdmidio},
	NamingBookStack:    {Short: PrefixBookStack, Name: NamingBookStack},
	NamingFirefly:      {Short: PrefixFirefly, Name: NamingFirefly},
	NamingImmich:       {Short: PrefixImmich, Name: NamingImmich},
	NamingMediaWiki:    {Short: PrefixMediaWiki, Name: NamingMediaWiki},
	NamingMinecraft:    {Short: PrefixMinecraft, Name: NamingMinecraft},
	NamingNextcloud:    {Short: PrefixNextcloud, Name: NamingNextcloud},
	NamingOCIS:         {Short: PrefixOCIS, Name: NamingOCIS},
	NamingPaperlessNGX: {Short: PrefixPaperlessNGX, Name: NamingPaperlessNGX},
	NamingRoundcube:    {Short: PrefixRoundcube, Name: NamingRoundcube},
	NamingRustDesk:     {Short: PrefixRustDesk, Name: NamingRustDesk},
	NamingUptimeKuma:   {Short: PrefixUptimeKuma, Name: NamingUptimeKuma},
	NamingVaultwarden:  {Short: PrefixVaultwarden, Name: NamingVaultwarden},
	NamingWordPress:    {Short: PrefixWordPress, Name: NamingWordPress},
}

// GetNamingScheme returns the configured naming scheme.
func GetNamingScheme(name string) *NamingScheme {
	ns, ok := namingSchemes[name]
	if !ok {
		panic("unknown naming scheme: " + name)
	}
	return ns
}

// Prefix returns "<short>-<name>".
func (ns *NamingScheme) Prefix(name string) string {
	return fmt.Sprintf("%s-%s", ns.Short, name)
}

// ApacheFile returns the Apache config filename.
func (ns *NamingScheme) ApacheFile(name string) string {
	return fmt.Sprintf("site-%s-%s.conf", ns.Short, name)
}

// ApachePath returns the Apache config path.
func (ns *NamingScheme) ApachePath(name string) string {
	return GetApacheEtcDir("sites-available", ns.ApacheFile(name))
}

// PhpFpmPoolFile returns the PHP-FPM pool filename.
func (ns *NamingScheme) PhpFpmPoolFile(name string) string {
	return fmt.Sprintf("%s-%s.conf", ns.Short, name)
}

// PhpFpmPoolPath returns the PHP-FPM pool config path.
func (ns *NamingScheme) PhpFpmPoolPath(name string) string {
	return php.GetPhpFpmPoolPath(0, ns.PhpFpmPoolFile(name))
}

// SystemdUnit returns the systemd unit name.
func (ns *NamingScheme) SystemdUnit(name string) string {
	return fmt.Sprintf("%s-%s.service", ns.Short, name)
}

// SQLiteFile returns the SQLite database filename.
func (ns *NamingScheme) SQLiteFile(name string) string {
	return fmt.Sprintf("%s-%s.sqlite", ns.Short, name)
}

// LogFile returns the log filename.
func (ns *NamingScheme) LogFile(name string) string {
	return fmt.Sprintf("%s-%s.log", ns.Short, name)
}

// CronName returns the cron.d filename.
func (ns *NamingScheme) CronName(name string) string {
	return fmt.Sprintf("%s_%s", ns.Name, name)
}

// CronPath returns the cron.d path.
func (ns *NamingScheme) CronPath(name string) string {
	return GetEtcDir("cron.d", ns.CronName(name))
}

// CertPath returns the certificate directory.
func (ns *NamingScheme) CertPath(fqdn string) string {
	return GetToolsDir("data", "certs", fqdn)
}

// HookName returns a hook filename.
func (ns *NamingScheme) HookName(action, name string) string {
	return fmt.Sprintf("%s-%s-%s", action, ns.Name, name)
}

// HookPath returns the hook path.
func (ns *NamingScheme) HookPath(action, name string) string {
	return GetToolsDir("data", "hooks", ns.HookName(action, name))
}

// LogsDir returns the logs directory with optional subpaths.
func (ns *NamingScheme) LogsDir(name string, paths ...string) string {
	base := GetToolsDir("logs", ns.Name, name)
	if len(paths) == 0 {
		return base
	}
	return filepath.Join(append([]string{base}, paths...)...)
}

// DataDir returns the data directory with optional subpaths.
func (ns *NamingScheme) DataDir(name string, paths ...string) string {
	base := GetToolsDir("data", ns.Name, name)
	if len(paths) == 0 {
		return base
	}
	return filepath.Join(append([]string{base}, paths...)...)
}

// CacheDir returns the cache directory with optional subpaths.
func (ns *NamingScheme) CacheDir(name string, paths ...string) string {
	base := GetToolsDir("cache", ns.Name, name)
	if len(paths) == 0 {
		return base
	}
	return filepath.Join(append([]string{base}, paths...)...)
}

// RuntimeDir returns the runtime directory with optional subpaths.
func (ns *NamingScheme) RuntimeDir(name string, paths ...string) string {
	base := GetToolsDir("run", ns.Name, name)
	if len(paths) == 0 {
		return base
	}
	return filepath.Join(append([]string{base}, paths...)...)
}

// SocketPath returns a socket path inside the runtime directory.
func (ns *NamingScheme) SocketPath(name, file string) string {
	return filepath.Join(ns.RuntimeDir(name), file)
}
