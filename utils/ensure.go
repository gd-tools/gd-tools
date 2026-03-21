package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigFile = "config.json"
)

// EnsureServerDir ensures that the current working directory represents
// a gd-tools production server directory.
//
// A server directory is identified by the presence of a config.json file.
func EnsureServerDir() error {
	if _, err := os.Stat(ConfigFile); err == nil {
		return nil
	}
	return fmt.Errorf(
		"this command must be executed from a server directory (config.json not found)",
	)
}

// EnsureBaseDir ensures that the current working directory is located
// within the gd-tools base directory defined by GD_TOOLS_BASE.
//
// This is used as a safety mechanism to prevent accidental execution
// outside the managed infrastructure tree.
func EnsureBaseDir() error {
	base := os.Getenv("GD_TOOLS_BASE")
	if base == "" {
		return fmt.Errorf(
			"GD_TOOLS_BASE is not set; please set it to your gd-tools base directory",
		)
	}

	baseAbs, err := filepath.Abs(base)
	if err != nil {
		return fmt.Errorf("failed to resolve GD_TOOLS_BASE: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to determine working directory: %w", err)
	}
	cwdAbs, err := filepath.Abs(cwd)
	if err != nil {
		return fmt.Errorf("failed to resolve working directory: %w", err)
	}

	// Ensure cwd is equal to or below base
	if cwdAbs == baseAbs {
		return nil
	}
	if strings.HasPrefix(cwdAbs, baseAbs+string(os.PathSeparator)) {
		return nil
	}

	return fmt.Errorf(
		"current directory is outside GD_TOOLS_BASE (%s)", baseAbs,
	)
}

// EnsureBaseOrServerDir allows execution either from the gd-tools base
// directory (or any of its subdirectories) or from a server directory.
//
// This is useful for commands like 'list' that can operate globally
// or on a single server.
func EnsureBaseOrServerDir() error {
	if err := EnsureBaseDir(); err == nil {
		return nil
	}
	if err := EnsureServerDir(); err == nil {
		return nil
	}

	return fmt.Errorf(
		"command must be executed either from the gd-tools base directory " +
			"(GD_TOOLS_BASE) or from a server directory (config.json present)",
	)
}
