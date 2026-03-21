package config

import (
	"crypto"

	"github.com/gd-tools/gd-tools/acme"
	"github.com/go-acme/lego/v4/certificate"
)

// GetPrivateKey loads the ACME account key.
// In tests, this can be overridden via cfg.getPrivateKey.
func (cfg *Config) GetPrivateKey(path string) (crypto.PrivateKey, error) {
	if cfg != nil {
		if fn := cfg.getPrivateKey; fn != nil {
			return fn(path)
		}
	}
	return acme.GetPrivateKey(path)
}

// GetCloudflareCertificate requests a certificate via the Cloudflare DNS provider.
// In tests, this can be overridden via cfg.getCloudflareCertificate.
func (cfg *Config) GetCloudflareCertificate(domains []string, email string, key crypto.PrivateKey) (*certificate.Resource, error) {
	if cfg != nil {
		if fn := cfg.getCloudflareCertificate; fn != nil {
			return fn(domains, email, key)
		}
	}
	return acme.GetCloudflareCertificate(domains, email, key)
}

// GetHetznerCertificate requests a certificate via the Hetzner DNS provider.
// In tests, this can be overridden via cfg.getHetznerCertificate.
func (cfg *Config) GetHetznerCertificate(domains []string, email string, key crypto.PrivateKey) (*certificate.Resource, error) {
	if cfg != nil {
		if fn := cfg.getHetznerCertificate; fn != nil {
			return fn(domains, email, key)
		}
	}
	return acme.GetHetznerCertificate(domains, email, key)
}

// GetIonosCertificate requests a certificate via the IONOS DNS provider.
// In tests, this can be overridden via cfg.getIonosCertificate.
func (cfg *Config) GetIonosCertificate(domains []string, email string, key crypto.PrivateKey) (*certificate.Resource, error) {
	if cfg != nil {
		if fn := cfg.getIonosCertificate; fn != nil {
			return fn(domains, email, key)
		}
	}
	return acme.GetIonosCertificate(domains, email, key)
}
