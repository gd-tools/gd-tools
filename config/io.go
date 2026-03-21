package config

import (
	"os"

	"github.com/gd-tools/gd-tools/utils"
)

// MkdirAll creates a directory path including missing parents.
// In tests, this can be overridden via cfg.mkdirAll.
func (cfg *Config) MkdirAll(path string, perm os.FileMode) error {
	if cfg != nil {
		if fn := cfg.mkdirAll; fn != nil {
			return fn(path, perm)
		}
	}
	return os.MkdirAll(path, perm)
}

// Chdir sets the current working directory.
// In tests, this can be overridden via cfg.chdir.
func (cfg *Config) Chdir(path string) error {
	if cfg != nil {
		if fn := cfg.chdir; fn != nil {
			return fn(path)
		}
	}
	return os.Chdir(path)
}

// Setenv sets an environment variable.
// In tests, this can be overridden via cfg.setenv.
func (cfg *Config) Setenv(key, value string) error {
	if cfg != nil {
		if fn := cfg.setenv; fn != nil {
			return fn(key, value)
		}
	}
	return os.Setenv(key, value)
}

// Unsetenv removes an environment variable.
// In tests, this can be overridden via cfg.unsetenv.
func (cfg *Config) Unsetenv(key string) error {
	if cfg != nil {
		if fn := cfg.unsetenv; fn != nil {
			return fn(key)
		}
	}
	return os.Unsetenv(key)
}

// LoadFile reads a file from disk.
// In tests, this can be overridden via cfg.loadFile.
func (cfg *Config) LoadFile(name string) ([]byte, error) {
	if cfg != nil {
		if fn := cfg.loadFile; fn != nil {
			return fn(name)
		}
	}
	return utils.LoadFile(name)
}

// LoadJSON reads a JSON file into the given structure.
// In tests, this can be overridden via cfg.loadJSON.
func (cfg *Config) LoadJSON(name string, v any) error {
	if cfg != nil {
		if fn := cfg.loadJSON; fn != nil {
			return fn(name, v)
		}
	}
	return utils.LoadJSON(name, v)
}

// SaveFile writes raw data to a file.
// In tests, this can be overridden via cfg.saveFile.
func (cfg *Config) SaveFile(name string, data []byte) error {
	if cfg != nil {
		if fn := cfg.saveFile; fn != nil {
			return fn(name, data)
		}
	}
	return utils.SaveFile(name, data)
}

// SaveJSON writes a structure as JSON to disk.
// In tests, this can be overridden via cfg.saveJSON.
func (cfg *Config) SaveJSON(name string, v any) error {
	if cfg != nil {
		if fn := cfg.saveJSON; fn != nil {
			return fn(name, v)
		}
	}
	return utils.SaveJSON(name, v)
}
