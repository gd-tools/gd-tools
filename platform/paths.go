package platform

import (
	"fmt"
	"path/filepath"
)

type Path struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

const (
	PathRoot      = "root"
	PathVar       = "var"
	PathTools     = "tools"
	PathEtc       = "etc"
	PathBin       = "bin"
	PathRun       = "run"
	PathDownloads = "downloads"
)

// DefaultPaths returns the production directory layout.
func DefaultPaths() []Path {
	return []Path{
		{Name: PathRoot, Value: "/root"},
		{Name: PathVar, Value: "/var"},
		{Name: PathTools, Value: "/var/gd-tools"},
		{Name: PathEtc, Value: "/etc"},
		{Name: PathBin, Value: "/usr/local/bin"},
		{Name: PathRun, Value: "/run"},
		{Name: PathDownloads, Value: "/root/Downloads"},
	}
}

// ClonePaths copies a path slice so tests can modify it safely.
func ClonePaths(src []Path) []Path {
	dst := make([]Path, len(src))
	copy(dst, src)
	return dst
}

func join(base string, paths ...string) string {
	if len(paths) == 0 {
		return base
	}
	return filepath.Join(append([]string{base}, paths...)...)
}

func (pf *Platform) pathValue(name string) (string, error) {
	for _, p := range pf.Paths {
		if p.Name == name {
			return p.Value, nil
		}
	}
	return "", fmt.Errorf("missing platform path: %s", name)
}

func (pf *Platform) MustPath(name string) string {
	v, err := pf.pathValue(name)
	if err != nil {
		panic(err)
	}
	return v
}

func (pf *Platform) RootDir(paths ...string) string {
	return join(pf.MustPath(PathRoot), paths...)
}

func (pf *Platform) VarDir(paths ...string) string {
	return join(pf.MustPath(PathVar), paths...)
}

func (pf *Platform) ToolsDir(paths ...string) string {
	return join(pf.MustPath(PathTools), paths...)
}

func (pf *Platform) DataDir(paths ...string) string {
	return join(pf.ToolsDir("data"), paths...)
}

func (pf *Platform) CertsDir(paths ...string) string {
	return join(pf.DataDir("certs"), paths...)
}

func (pf *Platform) LogsDir(paths ...string) string {
	return join(pf.ToolsDir("logs"), paths...)
}

func (pf *Platform) EtcDir(paths ...string) string {
	return join(pf.MustPath(PathEtc), paths...)
}

func (pf *Platform) BinDir(paths ...string) string {
	return join(pf.MustPath(PathBin), paths...)
}

func (pf *Platform) RunDir(paths ...string) string {
	return join(pf.MustPath(PathRun), paths...)
}

func (pf *Platform) DownloadsDir(paths ...string) string {
	return join(pf.MustPath(PathDownloads), paths...)
}

func (pf *Platform) ApacheToolsDir(paths ...string) string {
	return join(pf.DataDir("apache"), paths...)
}

func (pf *Platform) ApacheEtcDir(paths ...string) string {
	return join(pf.EtcDir("apache2"), paths...)
}
