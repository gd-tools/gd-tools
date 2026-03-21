package config

import (
	"net"

	"github.com/gd-tools/gd-tools/utils"
)

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
