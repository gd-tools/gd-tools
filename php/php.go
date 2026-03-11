package php

import (
	"fmt"
	"path/filepath"
)

type TemplateData struct {
	PhpVersion string
}

var (
	PhpVersion string
	PhpEtcDir  string
)

func Set_PhpVersion(version string) {
	PhpVersion = version
}

func Get_PhpVersion() string {
	if PhpVersion == "" {
		PhpVersion = "8.3" // default for Ubuntu 24.4 noble
	}

	return PhpVersion
}

func Get_PhpFpmService() string {
	return "php" + Get_PhpVersion() + "-fpm"
}

func Set_PhpEtcDir(path string) {
	if path == "" {
		path = "/etc/php/" + Get_PhpVersion()
	}
	PhpEtcDir = path
}

func Get_PhpEtcDir(paths ...string) string {
	if PhpEtcDir == "" {
		Set_PhpEtcDir("")
	}
	if len(paths) == 0 {
		return PhpEtcDir
	}
	return filepath.Join(append([]string{PhpEtcDir}, paths...)...)
}

func Get_PhpFpmPoolPath(number int, name string) string {
	filename := fmt.Sprintf("%02d-%s.conf", number, name)
	return Get_PhpEtcDir("fpm", "pool.d", filename)
}

func Get_TemplateData() TemplateData {
	return TemplateData{
		PhpVersion: Get_PhpVersion(),
	}
}
