package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/assets"
)

const (
	WordPressName = "wordpress"
	WordPressFile = WordPressName + ".json"
)

type WordPress struct {
	HostName     string   `json:"host_name"`
	DomainName   string   `json:"domain_name"`
	Version      string   `json:"version"`
	PhpVersion   string   `json:"php_version"`
	WpCliVersion string   `json:"wp_cli_version"`
	Aliases      []string `json:"aliases"`
	Language     string   `json:"language"`
	Region       string   `json:"region"`
	Password     string   `json:"password"`
	Directory    string   `json:"directory"`
	AdminName    string   `json:"admin_name"`
	AdminEmail   string   `json:"admin_email"`
	AdminPswd    string   `json:"admin_pswd"`
	Salt         string   `json:"salt,omitempty"`
}

type WordPressList struct {
	Entries []*WordPress `json:"entries"`
}

func (wp *WordPress) Locale() string {
	return wp.Language + "_" + wp.Region
}

func (wp *WordPress) FQDN() string {
	return wp.HostName + "." + wp.DomainName
}

func (wp *WordPress) IsWWW() bool {
	return wp.HostName == "www"
}

func (wp *WordPress) Name() string {
	if wp.IsWWW() {
		return agent.MakeDBName(wp.DomainName)
	}
	return agent.MakeDBName(wp.FQDN())
}

func (wp *WordPress) ServerAlias() string {
	aliases := []string{
		wp.FQDN(),
	}
	for _, dom := range wp.Aliases {
		aliases = append(aliases, "www."+dom, dom)
	}
	sort.Strings(aliases)
	return strings.Join(aliases, " ")
}

func (wp *WordPress) RootDir() string {
	return assets.GetToolsDir("data", WordPressName, wp.Name())
}

func (wp *WordPress) SocketPath() string {
	name := fmt.Sprintf("php%s-wordpress-%s.sock", wp.PhpVersion, wp.Name())
	return assets.GetRunDir("php", name)
}

func (wp *WordPress) ConfigPath() string {
	return filepath.Join(wp.RootDir(), "config.json")
}

func (wp *WordPress) BaseDir(paths ...string) string {
	baseDir := filepath.Join(wp.RootDir(), wp.Directory)
	if len(paths) == 0 {
		return baseDir
	}
	return filepath.Join(append([]string{baseDir}, paths...)...)
}

func (wp *WordPress) LogsDir(paths ...string) string {
	logsDir := assets.GetToolsDir("logs", WordPressName, wp.Name())
	if len(paths) == 0 {
		return logsDir
	}
	return filepath.Join(append([]string{logsDir}, paths...)...)
}

func (wp *WordPress) VhostPath() string {
	name := fmt.Sprintf("wp-%s.conf", wp.FQDN())
	return assets.GetApacheEtcDir("sites-available", name)
}

func (wp *WordPress) HookPath() string {
	name := fmt.Sprintf("backup-pre-%s-%s", WordPressName, wp.Name())
	return assets.GetToolsDir("data", "hooks", name)
}

func (wp *WordPress) CertDir() string {
	return assets.GetToolsDir("data", "certs", wp.FQDN())
}

func (wp *WordPress) CertificateList() (string, []string) {
	if !wp.IsWWW() {
		return wp.FQDN(), nil
	}

	domList := []string{
		wp.DomainName,
	}

	for _, dom := range wp.Aliases {
		domList = append(domList, "www."+dom, dom)
	}

	return wp.FQDN(), domList
}

func (wp *WordPress) NameList() []string {
	list := []string{
		wp.DomainName,
	}

	if wp.IsWWW() && len(wp.Aliases) > 0 {
		for _, alias := range wp.Aliases {
			list = append(list, alias)
		}
	}

	return list
}

func (wp *WordPress) SaltEntry(num int) string {
	return fmt.Sprintf("%s_%d", wp.Salt, num)
}

func (wp *WordPress) CronPath() string {
	name := WordPressName + "_" + wp.Name()
	return assets.GetEtcDir("cron.d", name)
}

func (wp *WordPress) WP_CLI_Path() string {
	return assets.GetBinDir("wp-" + wp.Name())
}

func LoadWordPressList(update *WordPress) (*WordPressList, error) {
	var list WordPressList

	content, err := os.ReadFile(WordPressFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &list, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", WordPressFile, err)
	}

	if err := json.Unmarshal(content, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", WordPressFile, err)
	}

	for index, _ := range list.Entries {
		entry := list.Entries[index]
		if entry.HostName == "" {
			return nil, fmt.Errorf("found WordPress without HostName")
		}
		if entry.DomainName == "" {
			return nil, fmt.Errorf("found WordPress without DomainName")
		}

		if entry.Directory == "" {
			return nil, fmt.Errorf("missing Directory for WordPress %s", entry.FQDN())
		}

		if update != nil && update.Password != "" {
			entry.Password = update.Password
		}
		if entry.Password == "" {
			return nil, fmt.Errorf("missing Password for WordPress %s", entry.FQDN())
		}

		if update != nil && update.AdminName != "" {
			entry.AdminName = update.AdminName
		}
		if entry.AdminName == "" {
			return nil, fmt.Errorf("missing AdminName for WordPress %s", entry.FQDN())
		}

		if update != nil && update.AdminEmail != "" {
			entry.AdminEmail = update.AdminEmail
		}
		if entry.AdminEmail == "" {
			return nil, fmt.Errorf("missing AdminEmail for WordPress %s", entry.FQDN())
		}

		if update != nil && update.AdminPswd != "" {
			entry.AdminPswd = update.AdminPswd
		}
		if entry.AdminPswd == "" {
			return nil, fmt.Errorf("missing AdminPswd for WordPress %s", entry.FQDN())
		}

		if update != nil && update.Salt != "" {
			entry.Salt = update.Salt
		}
		if entry.Salt == "" {
			return nil, fmt.Errorf("missing Salt for WordPress %s", entry.FQDN())
		}
	}

	if err := list.Save(); err != nil {
		return nil, err
	}

	if update != nil {
		return nil, nil
	}

	return &list, nil
}

func (list *WordPressList) Save() error {
	sort.Slice(list.Entries, func(i, j int) bool {
		return list.Entries[i].FQDN() < list.Entries[j].FQDN()
	})

	content, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", WordPressFile, err)
	}

	existing, err := os.ReadFile(WordPressFile)
	if err == nil && bytes.Equal(existing, content) {
		return nil
	}

	if err := os.WriteFile(WordPressFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", WordPressFile, err)
	}

	return nil
}

func (cfg *Config) DeployWordPress(wp *WordPress) error {
	if wp == nil {
		return fmt.Errorf("missing WordPress pointer")
	}
	cfg.Debugf("Enter config/wordpress.go (%s)", wp.FQDN())

	if wp.FQDN() == cfg.FQDN() {
		return fmt.Errorf("cannot use the server name for WordPress")
	}

	if err := cfg.WordPressDownload(wp); err != nil {
		return err
	}

	if err := cfg.WordPressExtract(wp); err != nil {
		return err
	}

	if err := cfg.WordPressLogsDir(wp); err != nil {
		return err
	}

	if err := cfg.WordPress_SQL(wp); err != nil {
		return err
	}

	if err := cfg.WordPressConfig(wp); err != nil {
		return err
	}

	if err := cfg.WordPressBackupHook(wp); err != nil {
		return err
	}

	if err := cfg.WordPressCrontab(wp); err != nil {
		return err
	}

	if err := cfg.WordPress_DNS(wp); err != nil {
		return err
	}

	if err := cfg.WordPressWpCli(wp); err != nil {
		return err
	}

	cfg.Debugf("Leave config/wordpress.go (%s)", wp.FQDN())
	return nil
}

func (cfg *Config) WordPressDownload(wp *WordPress) error {
	req := cfg.NewRequest()

	_, wpRel, err := cfg.Catalog.Get(WordPressName, wp.Version)
	if err != nil {
		return err
	}
	req.Downloads = append(req.Downloads, &wpRel.Download)

	_, cliRel, err := cfg.Catalog.Get("wp-cli", wp.WpCliVersion)
	if err != nil {
		return err
	}
	req.Downloads = append(req.Downloads, &cliRel.Download)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) WordPressExtract(wp *WordPress) error {
	req := cfg.NewRequest()

	_, wpRel, err := cfg.Catalog.Get(WordPressName, wp.Version)
	if err != nil {
		return err
	}

	downloadDir := assets.GetDownloadsDir(wpRel.Download.Filename)
	extract := agent.File{
		Task:   "extract",
		Path:   downloadDir,
		Target: wp.RootDir(),
		Mode:   "0755",
		User:   "root",
		Group:  "root",
	}
	req.AddFile(&extract)

	writableTask := agent.File{
		Task:  "process",
		Path:  wp.BaseDir(),
		User:  "www-data",
		Group: "www-data",
	}
	req.AddFile(&writableTask)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) WordPressLogsDir(wp *WordPress) error {
	req := cfg.NewRequest()

	logsMkdir := agent.File{
		Task:  "mkdir",
		Path:  wp.LogsDir(),
		Mode:  "0755",
		User:  "www-data",
		Group: "www-data",
	}
	req.AddFile(&logsMkdir)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) WordPress_SQL(wp *WordPress) error {
	req := cfg.NewRequest()

	sqlStmts, err := assets.SQL("wordpress/create.sql", wp)
	if err != nil {
		return err
	}

	sql := agent.MySQL{
		Stmts:   sqlStmts,
		Comment: fmt.Sprintf("create wordpress (%s) tables", wp.FQDN()),
	}
	req.MySQLs = append(req.MySQLs, &sql)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) WordPressBackupHook(wp *WordPress) error {
	req := cfg.NewRequest()

	hookContent, err := assets.Render("wordpress/backup", wp)
	if err != nil {
		return err
	}

	hookFile := agent.File{
		Task:    "write",
		Path:    wp.HookPath(),
		Content: hookContent,
		Mode:    "0500",
		User:    "root",
		Group:   "root",
	}
	req.AddFile(&hookFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) WordPressConfig(wp *WordPress) error {
	req := cfg.NewRequest()

	confTmpl, err := assets.Render("wordpress/wp-config.php", wp)
	if err != nil {
		return err
	}

	confFile := agent.File{
		Task:    "write",
		Path:    wp.BaseDir("wp-config.php"),
		Content: confTmpl,
		Mode:    "0644",
		User:    "www-data",
		Group:   "www-data",
		Service: "apache2",
	}
	req.AddFile(&confFile)

	poolTmpl, err := assets.Render("wordpress/php-fpm-pool.conf", wp)
	if err != nil {
		return err
	}

	poolPath := cfg.PhpFpmPoolPath("wordpress-" + wp.Name())
	poolFile := agent.File{
		Task:    "write",
		Path:    poolPath,
		Content: poolTmpl,
		Mode:    "0644",
		User:    "root",
		Group:   "root",
		Service: cfg.PhpFpmService(),
	}
	req.AddFile(&poolFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) WordPressCrontab(wp *WordPress) error {
	req := cfg.NewRequest()

	cronTmpl, err := assets.Render("wordpress/cron.d", wp)
	if err != nil {
		return err
	}

	cronFile := agent.File{
		Task:    "write",
		Path:    wp.CronPath(),
		Content: cronTmpl,
		Mode:    "0644",
		User:    "root",
		Group:   "root",
	}
	req.AddFile(&cronFile)

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) WordPress_DNS(wp *WordPress) error {
	// get certificates for all possible names via DNS-01
	fqdnCert, sanCerts := wp.CertificateList()
	if err := cfg.EnsureCertificate(fqdnCert, sanCerts...); err != nil {
		return err
	}

	// setup vhost with all aliases, www and domain
	if err := cfg.WordPressSetupVhost(wp); err != nil {
		return err
	}

	if cfg.SkipDNS {
		cfg.Sayf("skipping dns-update for %s", wp.FQDN())
		return nil
	}

	for _, name := range wp.NameList() {
		// install all www entries as CNAME records
		// The NameList contains at least the FQDN of the WordPress
		if status, err := cfg.SetCNAME(name, wp.HostName); err != nil {
			return err
		} else if status != "" {
			cfg.Say(status)
		}

		// for non-www servers (e.g. demo.example.com) there is nothing more to do
		if !wp.IsWWW() {
			continue
		}

		// install the domain A and AAAA records (those cannot be CNAME)
		if cfg.IPv4Addr != "" {
			if status, err := cfg.SetA(name, "@", cfg.IPv4Addr); err != nil {
				return err
			} else if status != "" {
				cfg.Say(status)
			}
		}
		if cfg.IPv6Addr != "" {
			if status, err := cfg.SetAAAA(name, "@", cfg.IPv6Addr); err != nil {
				return err
			} else if status != "" {
				cfg.Say(status)
			}
		}
	}

	return nil
}

func (cfg *Config) WordPressSetupVhost(wp *WordPress) error {
	req := cfg.NewRequest()

	vhostTmpl, err := assets.Render("wordpress/vhost.conf", wp)
	if err != nil {
		return err
	}

	vhostFile := agent.File{
		Task:    "write",
		Path:    wp.VhostPath(),
		Content: vhostTmpl,
	}
	req.AddFile(&vhostFile)

	req.AddService("apache2")

	if err := req.Send(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) WordPressWpCli(wp *WordPress) error {
	wpSrc := assets.GetBinDir("gd-wp-cli")
	wpDst := assets.GetBinDir("wp-" + wp.Name())

	if _, err := cfg.LocalCommand(
		"rsync",
		cfg.RsyncFlags(),
		"--chown=root:root",
		"--chmod=0755",
		wpSrc,
		cfg.RootUser()+":"+wpDst,
	); err != nil {
		return fmt.Errorf("failed to install %s: %w", wpDst, err)
	}

	return nil
}
