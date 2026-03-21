package config

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gd-tools/gd-tools/utils"
	"github.com/go-acme/lego/v4/certificate"
)

const (
	ACMECertDir    = "acme-certs"
	ACMEAccountKey = "../acme_account.key"
)

var (
	RenewalTime = 30 * 24 * time.Hour
)

// EnsureCertificate ensures that a certificate exists for the given domain and SANs.
// It reuses an existing certificate if it is still valid, not close to expiry,
// and its SAN set matches the requested domains.
func (cfg *Config) EnsureCertificate(domain string, sans ...string) error {
	if domain == "" {
		return fmt.Errorf("missing domain in certificate request")
	}

	var wantDomains utils.LineBuffer
	wantDomains.Append(sans...)
	wantDomains.NormalizeWithFirst(domain)

	baseDir := filepath.Join(ACMECertDir, domain)
	fullchainPath := filepath.Join(baseDir, "fullchain.pem")
	privkeyPath := filepath.Join(baseDir, "privkey.pem")
	issuerPath := filepath.Join(baseDir, "issuer.pem")

	// Issue a new certificate if expiry is near or SANs have changed.
	// Renew also if 'gdt cert' was called with --force.
	fullchain, err1 := cfg.LoadFile(fullchainPath)
	_, err2 := cfg.LoadFile(privkeyPath)
	_, err3 := cfg.LoadFile(issuerPath)

	if err1 == nil && err2 == nil && err3 == nil {
		validUntil, err := ReadValidUntil(fullchain)
		if err != nil {
			return err
		}

		gotDomains, err := ExtractSANList(domain, fullchain)
		if err != nil {
			return err
		}

		if !wantDomains.IsEqual(gotDomains) {
			cfg.Infof("certificate SANs changed, requesting new")
		} else if time.Until(validUntil) <= RenewalTime {
			cfg.Infof("certificate expires soon on %s, requesting new", validUntil.Format("2006-01-02"))
		} else if cfg.Force {
			// Fall through: 'gdt cert' was called with --force
		} else {
			cfg.Infof(
				"certificate valid until %s: %s",
				validUntil.Format("2006-01-02"),
				strings.Join(gotDomains.Lines(), " "),
			)
			return nil
		}
	}

	key, err := cfg.GetPrivateKey(ACMEAccountKey)
	if err != nil {
		return err
	}

	var resource *certificate.Resource

	switch {
	case cfg.CloudflareToken != "":
		if err := cfg.Setenv("CF_DNS_API_TOKEN", cfg.CloudflareToken); err != nil {
			return err
		}
		defer func() {
			_ = cfg.Unsetenv("CF_DNS_API_TOKEN")
		}()

		resource, err = cfg.GetCloudflareCertificate(wantDomains.Lines(), cfg.SysAdmin, key)

	case cfg.HetznerToken != "":
		if err := cfg.Setenv("HETZNER_API_TOKEN", cfg.HetznerToken); err != nil {
			return err
		}
		defer func() {
			_ = cfg.Unsetenv("HETZNER_API_TOKEN")
		}()

		resource, err = cfg.GetHetznerCertificate(wantDomains.Lines(), cfg.SysAdmin, key)

	case cfg.IonosToken != "":
		if err := cfg.Setenv("IONOS_API_KEY", cfg.IonosToken); err != nil {
			return err
		}
		defer func() {
			_ = cfg.Unsetenv("IONOS_API_KEY")
		}()

		resource, err = cfg.GetIonosCertificate(wantDomains.Lines(), cfg.SysAdmin, key)

	default:
		resource = nil
	}

	if err != nil {
		return err
	}
	if resource == nil {
		return fmt.Errorf("missing provider for DNS-01 certificates")
	}

	if err := cfg.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to mkdir %s: %w", baseDir, err)
	}

	if err := cfg.SaveFile(filepath.Join(baseDir, "fullchain.pem"), resource.Certificate); err != nil {
		return err
	}
	if err := cfg.SaveFile(filepath.Join(baseDir, "privkey.pem"), resource.PrivateKey); err != nil {
		return err
	}
	if err := cfg.SaveFile(filepath.Join(baseDir, "issuer.pem"), resource.IssuerCertificate); err != nil {
		return err
	}

	validUntil, err := ReadValidUntil(resource.Certificate)
	if err != nil {
		return err
	}

	gotDomains, err := ExtractSANList(domain, resource.Certificate)
	if err != nil {
		return err
	}

	cfg.Infof(
		"certificate valid until %s: %s",
		validUntil.Format("2006-01-02"),
		strings.Join(gotDomains.Lines(), " "),
	)

	return nil
}

// ExtractSANList extracts the DNS SAN list from the first certificate in a PEM chain.
// The given domain is normalized as first entry and added if missing.
func ExtractSANList(domain string, fullchain []byte) (*utils.LineBuffer, error) {
	var certBlock *pem.Block
	rest := fullchain

	for {
		certBlock, rest = pem.Decode(rest)
		if certBlock == nil {
			return nil, fmt.Errorf("no certificate found in fullchain")
		}
		if certBlock.Type == "CERTIFICATE" {
			break
		}
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %w", err)
	}

	var gotDomains utils.LineBuffer
	gotDomains.Append(cert.DNSNames...)
	gotDomains.NormalizeWithFirst(domain)

	return &gotDomains, nil
}

// ReadValidUntil returns the NotAfter timestamp of the first certificate found in PEM data.
func ReadValidUntil(pemBytes []byte) (time.Time, error) {
	var block *pem.Block
	rest := pemBytes

	for {
		if block, rest = pem.Decode(rest); block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				return cert.NotAfter.UTC().Round(0), nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("no certificate in response")
}
