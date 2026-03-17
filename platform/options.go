package platform

import (
	"net"
	"path/filepath"
)

type Options struct {
	rootDir string
	varDir  string
	etcDir  string
	binDir  string
	runDir  string

	LookupIP      func(host string) ([]net.IP, error)
	GetRSAKeyPair func(fqdn string) ([]byte, []byte, error)
}

func defaultOptions() *Options {
	return &Options{
		rootDir: "/root",
		varDir:  "/var",
		etcDir:  "/etc",
		binDir:  "/usr/local/bin",
		runDir:  "/run",

		LookupIP:      net.LookupIP,
		GetRSAKeyPair: ProdRSAKeyPair,
	}
}

func (pf *Platform) joinPath(base string, paths ...string) string {
	elems := append([]string{base}, paths...)
	return filepath.Join(elems...)
}

// These paths are fundamental, can be faked for tests.
func (pf *Platform) RootPath(paths ...string) string {
	if pf.options == nil || pf.options.rootDir == "" {
		panic("missing options or rootDir")
	}
	return pf.joinPath(pf.options.rootDir, paths...)
}

func (pf *Platform) VarPath(paths ...string) string {
	if pf.options == nil || pf.options.varDir == "" {
		panic("missing options or varDir")
	}
	return pf.joinPath(pf.options.varDir, paths...)
}

func (pf *Platform) EtcPath(paths ...string) string {
	if pf.options == nil || pf.options.etcDir == "" {
		panic("missing options or etcDir")
	}
	return pf.joinPath(pf.options.etcDir, paths...)
}

func (pf *Platform) BinPath(paths ...string) string {
	if pf.options == nil || pf.options.binDir == "" {
		panic("missing options or binDir")
	}
	return pf.joinPath(pf.options.binDir, paths...)
}

func (pf *Platform) RunPath(paths ...string) string {
	if pf.options == nil || pf.options.runDir == "" {
		panic("missing options or runDir")
	}
	return pf.joinPath(pf.options.runDir, paths...)
}

// These paths are derived from the fundamental paths.
func (pf *Platform) DownloadsPath(paths ...string) string {
	return pf.joinPath(pf.RootPath("Downloads"), paths...)
}

func (pf *Platform) ToolsPath(paths ...string) string {
	return pf.joinPath(pf.VarPath("gd-tools"), paths...)
}

func (pf *Platform) DataPath(paths ...string) string {
	return pf.joinPath(pf.ToolsPath("data"), paths...)
}

func (pf *Platform) CertsPath(paths ...string) string {
	return pf.joinPath(pf.DataPath("certs"), paths...)
}

func (pf *Platform) ToolsApachePath(paths ...string) string {
	return pf.joinPath(pf.DataPath("apache"), paths...)
}

func (pf *Platform) LogsPath(paths ...string) string {
	return pf.joinPath(pf.ToolsPath("logs"), paths...)
}

func (pf *Platform) EtcApachePath(paths ...string) string {
	return pf.joinPath(pf.EtcPath("apache2"), paths...)
}

func (pf *Platform) EtcPhpPath(paths ...string) string {
	if pf.Baseline == nil || pf.Baseline.PHP == "" {
		panic("missing Baseline or PHP version")
	}
	return pf.joinPath(pf.EtcPath("php", pf.Baseline.PHP), paths...)
}
