package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gd-tools/gd-tools/assets"
)

func (file *File) IsApacheAvailable() (string, string, bool) {
	for _, kind := range []string{"conf", "mods", "sites"} {
		dir := assets.GetApacheEtcDir(kind + "-available/")
		if module, ok := strings.CutPrefix(file.Path, dir); ok {
			return kind, module, true
		}
	}

	return "", "", false
}

func (file *File) ApacheEnable(kind, module string) error {
	src := assets.GetApacheEtcDir(kind+"-available", module)
	dest := assets.GetApacheEtcDir(kind+"-enabled", module)

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("failed to create symlink for %s: %w", module, err)
	}
	if _, err := os.Lstat(dest); err == nil {
		return nil
	}

	relSrc := filepath.Join("..", kind+"-available", module)
	if err := os.Symlink(relSrc, dest); err != nil {
		return fmt.Errorf("failed to create symlink for %s: %w", module, err)
	}

	return nil
}

func (file *File) ApacheDisable(kind, module string) error {
	enabled := assets.GetApacheEtcDir(kind + "-enabled")

	for _, ext := range []string{".load", ".conf"} {
		dest := filepath.Join(enabled, module+ext)

		if _, err := os.Lstat(dest); err == nil {
			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("failed to remove symlink for %s: %w", module, err)
			}
		}
	}

	return nil
}
