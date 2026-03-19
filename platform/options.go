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

	LoadJSON func(name string, v any) error
	SaveFile func(name string, data []byte) error
	SaveJSON func(name string, v any) error

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

		LoadJSON: utils.LoadJSON,
		SaveFile: utils.SaveFile,
		SaveJSON: utils.SaveJSON,

		LookupIP:   net.LookupIP,
		DHParams:   utils.DHParams,
		RSAKeyPair: utils.RSAKeyPair,
		RunShell:   utils.RunShell,
		RunCommand: utils.RunCommand,
	}
}

// These paths are fundamental, can be faked for tests.
func (pf *Platform) RootPath(paths ...string) (string, error) {
	if pf == nil {
		return "", fmt.Errorf("RootPath: missing pf")
	}
	if pf.options == nil {
		return "", fmt.Errorf("RootPath: missing pf.options")
	}
	if pf.options.rootDir == "" {
		return "", fmt.Errorf("RootPath: missing pf.options.rootDir")
	}
	return utils.JoinPath(pf.options.rootDir, paths...), nil
}

func (pf *Platform) VarPath(paths ...string) (string, error) {
	if pf == nil {
		return "", fmt.Errorf("VarPath: missing pf")
	}
	if pf.options == nil {
		return "", fmt.Errorf("VarPath: missing pf.options")
	}
	if pf.options.varDir == "" {
		return "", fmt.Errorf("VarPath: missing pf.options.varDir")
	}
	return utils.JoinPath(pf.options.varDir, paths...), nil
}

func (pf *Platform) EtcPath(paths ...string) (string, error) {
	if pf == nil {
		return "", fmt.Errorf("EtcPath: missing pf")
	}
	if pf.options == nil {
		return "", fmt.Errorf("EtcPath: missing pf.options")
	}
	if pf.options.etcDir == "" {
		return "", fmt.Errorf("EtcPath: missing pf.options.etcDir")
	}
	return utils.JoinPath(pf.options.etcDir, paths...), nil
}

func (pf *Platform) BinPath(paths ...string) (string, error) {
	if pf == nil {
		return "", fmt.Errorf("BinPath: missing pf")
	}
	if pf.options == nil {
		return "", fmt.Errorf("BinPath: missing pf.options")
	}
	if pf.options.binDir == "" {
		return "", fmt.Errorf("BinPath: missing pf.options.binir")
	}
	return utils.JoinPath(pf.options.binDir, paths...), nil
}

func (pf *Platform) RunPath(paths ...string) (string, error) {
	if pf == nil {
		return "", fmt.Errorf("RunPath: missing pf")
	}
	if pf.options == nil {
		return "", fmt.Errorf("RunPath: missing pf.options")
	}
	if pf.options.runDir == "" {
		return "", fmt.Errorf("RunPath: missing pf.options.runDir")
	}
	return utils.JoinPath(pf.options.runDir, paths...), nil
}

// These paths are derived from the fundamental paths.
func (pf *Platform) DownloadsPath(paths ...string) (string, error) {
	base, err := pf.RootPath("Downloads")
	if err != nil {
		return "", err
	}
	return utils.JoinPath(base, paths...), nil
}

func (pf *Platform) ToolsPath(paths ...string) (string, error) {
	base, err := pf.VarPath("gd-tools")
	if err != nil {
		return "", err
	}
	return utils.JoinPath(base, paths...), nil
}

func (pf *Platform) DataPath(paths ...string) (string, error) {
	base, err := pf.ToolsPath("data")
	if err != nil {
		return "", err
	}
	return utils.JoinPath(base, paths...), nil
}

func (pf *Platform) CertsPath(paths ...string) (string, error) {
	base, err := pf.DataPath("certs")
	if err != nil {
		return "", err
	}
	return utils.JoinPath(base, paths...), nil
}

func (pf *Platform) ToolsApachePath(paths ...string) (string, error) {
	base, err := pf.DataPath("apache")
	if err != nil {
		return "", err
	}
	return utils.JoinPath(base, paths...), nil
}

func (pf *Platform) LogsPath(paths ...string) (string, error) {
	base, err := pf.ToolsPath("logs")
	if err != nil {
		return "", err
	}
	return utils.JoinPath(base, paths...), nil
}

func (pf *Platform) EtcApachePath(paths ...string) (string, error) {
	base, err := pf.EtcPath("apache2")
	if err != nil {
		return "", err
	}
	return utils.JoinPath(base, paths...)
}

func (pf *Platform) EtcPhpPath(paths ...string) (string, error) {
	if pf == nil {
		return "", fmt.Errorf("EtcPhpPath: missing pf")
	}
	if pf.Baseline == nil {
		return "", fmt.Errorf("EtcPhpPath: missing pf.Baseline")
	}
	if pf.Baseline.PHP == "" {
		return "", fmt.Errorf("EtcPhpPath: missing pf.Baseline.PHP")
	}
	base, err := pf.EtcPath("php", pf.Baseline.PHP)
	if err != nil {
		return "", err
	}
	return utils.JoinPath(base, paths...), nil
}

// These are the file load/save functions.
func (pf *Platform) LoadJSON(name string, v any) error {
	if pf == nil {
		return fmt.Errorf("LoadJSON: missing pf")
	}
	if pf.options == nil {
		return fmt.Errorf("LoadJSON: missing pf.options")
	}
	if pf.options.LoadJSON == nil {
		return fmt.Errorf("LoadJSON: missing pf.options.LoadJSON")
	}
	return pf.options.LoadJSON(name, v)
}

func (pf *Platform) SaveFile(name string, data []byte) error {
	if pf == nil {
		return fmt.Errorf("SaveFile: missing pf")
	}
	if pf.options == nil {
		return fmt.Errorf("SaveFile: missing pf.options")
	}
	if pf.options.SaveFile == nil {
		return fmt.Errorf("SaveFile: missing pf.options.SaveFile")
	}
	return pf.options.SaveFile(name, data)
}

func (pf *Platform) SaveJSON(name string, v any) error {
	if pf == nil {
		return fmt.Errorf("SaveJSON: missing pf")
	}
	if pf.options == nil {
		return fmt.Errorf("SaveJSON: missing pf.options")
	}
	if pf.options.SaveJSON == nil {
		return fmt.Errorf("SaveJSON: missing pf.options.SaveJSON")
	}
	return pf.options.SaveJSON(name, v)
}

// These are the helper functions with side effects.
func (pf *Platform) LookupIP(host string) ([]net.IP, error) {
	if pf == nil {
		return nil, fmt.Errorf("LookupIP: missing pf")
	}
	if pf.options == nil {
		return nil, fmt.Errorf("LookupIP: missing pf.options")
	}
	if pf.options.LookupIP == nil {
		return nil, fmt.Errorf("LookupIP: missing pf.options.LookupIP")
	}
	return pf.options.LookupIP(host)
}

func (pf *Platform) RSAKeyPair(fqdn string) ([]byte, []byte, error) {
	if pf == nil {
		return nil, fmt.Errorf("RSAKeyPair: missing pf")
	}
	if pf.options == nil {
		return nil, fmt.Errorf("RSAKeyPair: missing pf.options")
	}
	if pf.options.RSAKeyPair == nil {
		return nil, fmt.Errorf("RSAKeyPair: missing pf.options.RSAKeyPair")
	}
	return pf.options.RSAKeyPair(fqdn)
}

// RunShell executes a list of shell commands on the local system.
func (pf *Platform) RunShell(commands []string) ([]byte, error) {
	if pf == nil {
		return nil, fmt.Errorf("RunShell: missing pf")
	}
	if pf.options == nil {
		return nil, fmt.Errorf("RunShell: missing pf.options")
	}
	if pf.options.RunShell == nil {
		return nil, fmt.Errorf("RunShell: missing pf.options.RunShell")
	}
	return pf.options.RunShell(commands)
}

func (pf *Platform) RunCommand(name string, args ...string) ([]byte, error) {
	if pf == nil {
		return nil, fmt.Errorf("RunCommand: missing pf")
	}
	if pf.options == nil {
		return nil, fmt.Errorf("RunCommand: missing pf.options")
	}
	if pf.options.RunCommand == nil {
		return nil, fmt.Errorf("RunCommand: missing pf.options.RunCommand")
	}
	return pf.options.RunCommand(name, args...)
}
