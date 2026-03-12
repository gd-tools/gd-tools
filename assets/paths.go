package assets

import (
	"path/filepath"
)

type Download struct {
	DownloadURL string `json:"download_url"`
	Filename    string `json:"filename"`
	Directory   string `json:"directory"`
	Binary      string `json:"binary"`
	MD5         string `json:"md5"`
	SHA256      string `json:"sha256"`
	SHA512      string `json:"sha512"`
}

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

	all := append([]string{base}, paths...)
	return filepath.Join(all...)
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
	toolsDir := GetToolsDir("data", "apache")
	if len(paths) == 0 {
		return toolsDir
	}
	return filepath.Join(append([]string{toolsDir}, paths...)...)
}

func GetApacheEtcDir(paths ...string) string {
	etcDir := GetEtcDir("apache2")
	if len(paths) == 0 {
		return etcDir
	}
	return filepath.Join(append([]string{etcDir}, paths...)...)
}
