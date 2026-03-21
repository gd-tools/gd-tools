package config

import (
	"net"

	"github.com/gd-tools/gd-tools/utils"
)

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

// LookupIP resolves a hostname to IP addresses.
// In tests, this can be overridden via cfg.lookupIP.
func (cfg *Config) LookupIP(host string) ([]net.IP, error) {
	if cfg != nil {
		if fn := cfg.lookupIP; fn != nil {
			return fn(host)
		}
	}
	return net.LookupIP(host)
}

// DHParams generates Diffie-Hellman parameters.
// In tests, this can be overridden via cfg.dhParams.
func (cfg *Config) DHParams(bits int) ([]byte, error) {
	if cfg != nil {
		if fn := cfg.dhParams; fn != nil {
			return fn(bits)
		}
	}
	return utils.DHParams(bits)
}

// RSAKeyPair generates an RSA key pair for a given FQDN.
// In tests, this can be overridden via cfg.rsaKeyPair.
func (cfg *Config) RSAKeyPair(fqdn string) ([]byte, []byte, error) {
	if cfg != nil {
		if fn := cfg.rsaKeyPair; fn != nil {
			return fn(fqdn)
		}
	}
	return utils.RSAKeyPair(fqdn)
}

// RunShell executes a list of shell commands.
// In tests, this can be overridden via cfg.runShell.
func (cfg *Config) RunShell(commands []string) ([]byte, error) {
	if cfg != nil {
		if fn := cfg.runShell; fn != nil {
			return fn(commands)
		}
	}
	return utils.RunShell(commands)
}

// RunCommand executes a system command.
// In tests, this can be overridden via cfg.runCommand.
func (cfg *Config) RunCommand(name string, args ...string) ([]byte, error) {
	if cfg != nil {
		if fn := cfg.runCommand; fn != nil {
			return fn(name, args...)
		}
	}
	return utils.RunCommand(name, args...)
}
