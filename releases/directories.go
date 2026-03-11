package releases

import (
	"os"
	"path/filepath"
)

var (
	RootDir  string
	VarDir   string
	ToolsDir string
	EtcDir   string
	BinDir   string
)

// RootDir is the home of user root
func SetRootDir(path string) {
	if path == "" {
		path = os.Getenv("GD_TOOLS_ROOT_DIR")
	}
	if path == "" {
		path = "/root"
	}
	RootDir = path
}

func GetRootDir(paths ...string) string {
	SetRootDir("")
	if len(paths) == 0 {
		return RootDir
	}
	return filepath.Join(append([]string{RootDir}, paths...)...)
}

// VarDir is the parent of the mounted gd-tools filesystem
func SetVarDir(path string) {
	if path == "" {
		path = os.Getenv("GD_TOOLS_VAR_DIR")
	}
	if path == "" {
		path = "/var"
	}
	VarDir = path
}

func GetVarDir(paths ...string) string {
	SetVarDir("")
	if len(paths) == 0 {
		return VarDir
	}
	return filepath.Join(append([]string{VarDir}, paths...)...)
}

// ToolsDir is the home of all data to be saved (and not constructed)
func SetToolsDir(path string) {
	if path == "" {
		path = os.Getenv("GD_TOOLS_TOOLS_DIR")
	}
	if path == "" {
		path = "/var/gd-tools"
	}
	ToolsDir = path
}

func GetToolsDir(paths ...string) string {
	SetToolsDir("")
	if len(paths) == 0 {
		return ToolsDir
	}
	return filepath.Join(append([]string{ToolsDir}, paths...)...)
}

// EtcDir is the system /etc dir, can be changed for testing
func SetEtcDir(path string) {
	if path == "" {
		path = os.Getenv("GD_TOOLS_ETC_DIR")
	}
	if path == "" {
		path = "/etc"
	}
	EtcDir = path
}

func GetEtcDir(paths ...string) string {
	SetEtcDir("")
	if len(paths) == 0 {
		return EtcDir
	}
	return filepath.Join(append([]string{EtcDir}, paths...)...)
}

// BinDir is the local binary dir, can be changed for testing
func SetBinDir(path string) {
	if path == "" {
		path = os.Getenv("GD_TOOLS_BIN_DIR")
	}
	if path == "" {
		path = "/usr/local/bin"
	}
	BinDir = path
}

func GetBinDir(paths ...string) string {
	SetBinDir("")
	if len(paths) == 0 {
		return BinDir
	}
	return filepath.Join(append([]string{BinDir}, paths...)...)
}
