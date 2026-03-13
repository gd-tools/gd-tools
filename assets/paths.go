package assets

import (
	"path/filepath"
)

type directories struct {
	rootDir      string
	varDir       string
	toolsDir     string
	etcDir       string
	binDir       string
	runDir       string
	downloadsDir string
}

var defaultDirs = directories{
	rootDir:      "/root",
	varDir:       "/var",
	toolsDir:     "/var/gd-tools",
	etcDir:       "/etc",
	binDir:       "/usr/local/bin",
	runDir:       "/run",
	downloadsDir: "/root/Downloads",
}

func join(base string, paths ...string) string {
	if len(paths) == 0 {
		return base
	}
	return filepath.Join(append([]string{base}, paths...)...)
}

func GetRootDir(paths ...string) string {
	return join(defaultDirs.rootDir, paths...)
}

func GetVarDir(paths ...string) string {
	return join(defaultDirs.varDir, paths...)
}

func GetToolsDir(paths ...string) string {
	return join(defaultDirs.toolsDir, paths...)
}

func GetEtcDir(paths ...string) string {
	return join(defaultDirs.etcDir, paths...)
}

func GetBinDir(paths ...string) string {
	return join(defaultDirs.binDir, paths...)
}

func GetRunDir(paths ...string) string {
	return join(defaultDirs.runDir, paths...)
}

func GetDownloadsDir(paths ...string) string {
	return join(defaultDirs.downloadsDir, paths...)
}

func GetApacheToolsDir(paths ...string) string {
	return join(GetToolsDir("data", "apache"), paths...)
}

func GetApacheEtcDir(paths ...string) string {
	return join(GetEtcDir("apache2"), paths...)
}
