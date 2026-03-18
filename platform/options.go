package platform

import (
	"net"

	"github.com/gd-tools/gd-tools/utils"
)

type Options struct {
	rootDir string
	varDir  string
	etcDir  string
	binDir  string
	runDir  string

	LookupIP   func(host string) ([]net.IP, error)
	DHParams   func(bits int) ([]byte, error)
	RSAKeyPair func(fqdn string) ([]byte, []byte, error)
	RunShell   func(commands []string) ([]byte, error)
	RunCommand func(name string, args ...string) ([]byte, error)
}

func defaultOptions() *Options {
	return &Options{
		rootDir: "/root",
		varDir:  "/var",
		etcDir:  "/etc",
		binDir:  "/usr/local/bin",
		runDir:  "/run",

		LookupIP:   net.LookupIP,
		DHParams:   utils.DHParams,
		RSAKeyPair: utils.RSAKeyPair,
		RunShell:   utils.RunShell,
		RunCommand: utils.RunCommand,
	}
}

// These paths are fundamental, can be faked for tests.
func (pf *Platform) RootPath(paths ...string) string {
	if pf.options == nil || pf.options.rootDir == "" {
		panic("missing options or rootDir")
	}
	return utils.JoinPath(pf.options.rootDir, paths...)
}

func (pf *Platform) VarPath(paths ...string) string {
	if pf.options == nil || pf.options.varDir == "" {
		panic("missing options or varDir")
	}
	return utils.JoinPath(pf.options.varDir, paths...)
}

func (pf *Platform) EtcPath(paths ...string) string {
	if pf.options == nil || pf.options.etcDir == "" {
		panic("missing options or etcDir")
	}
	return utils.JoinPath(pf.options.etcDir, paths...)
}

func (pf *Platform) BinPath(paths ...string) string {
	if pf.options == nil || pf.options.binDir == "" {
		panic("missing options or binDir")
	}
	return utils.JoinPath(pf.options.binDir, paths...)
}

func (pf *Platform) RunPath(paths ...string) string {
	if pf.options == nil || pf.options.runDir == "" {
		panic("missing options or runDir")
	}
	return utils.JoinPath(pf.options.runDir, paths...)
}

// These paths are derived from the fundamental paths.
func (pf *Platform) DownloadsPath(paths ...string) string {
	return utils.JoinPath(pf.RootPath("Downloads"), paths...)
}

func (pf *Platform) ToolsPath(paths ...string) string {
	return utils.JoinPath(pf.VarPath("gd-tools"), paths...)
}

func (pf *Platform) DataPath(paths ...string) string {
	return utils.JoinPath(pf.ToolsPath("data"), paths...)
}

func (pf *Platform) CertsPath(paths ...string) string {
	return utils.JoinPath(pf.DataPath("certs"), paths...)
}

func (pf *Platform) ToolsApachePath(paths ...string) string {
	return utils.JoinPath(pf.DataPath("apache"), paths...)
}

func (pf *Platform) LogsPath(paths ...string) string {
	return utils.JoinPath(pf.ToolsPath("logs"), paths...)
}

func (pf *Platform) EtcApachePath(paths ...string) string {
	return utils.JoinPath(pf.EtcPath("apache2"), paths...)
}

func (pf *Platform) EtcPhpPath(paths ...string) string {
	if pf.Baseline == nil || pf.Baseline.PHP == "" {
		panic("missing Baseline or PHP version")
	}
	return utils.JoinPath(pf.EtcPath("php", pf.Baseline.PHP), paths...)
}

// These are the helper functions with side effects.
func (pf *Platform) RunCommand(name string, args ...string) ([]byte, error) {
	return pf.options.RunCommand(name, args...)
}

func (pf *Platform) LookupIP(host string) ([]net.IP, error) {
	return pf.options.LookupIP(host)
}

func (pf *Platform) RSAKeyPair(fqdn string) ([]byte, []byte, error) {
	return pf.options.RSAKeyPair(fqdn)
}

// RunShell executes a list of shell commands on the local system.
func (pf *Platform) RunShell(commands []string) ([]byte, error) {
	return pf.options.RunShell(commands)
}
