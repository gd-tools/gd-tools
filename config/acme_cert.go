package config

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gd-tools/gd-tools/acme"
	"github.com/go-acme/lego/v4/certificate"
)

const (
	ACMECertDir    = "acme-certs"
	ACMEAccountKey = "../acme_account.key"
)

func (cfg *Config) EnsureCertificate(domain string, sans ...string) error {
	if domain == "" {
		return fmt.Errorf("missing domain in certificate request")
	}
	force := cfg.Force

	sans = UniqueSortedStrings(sans)
	wantDomains := UniqueSortedStrings(append([]string{domain}, sans...))

	baseDir := filepath.Join(ACMECertDir, domain)
	fullchainPath := filepath.Join(baseDir, "fullchain.pem")
	privkeyPath := filepath.Join(baseDir, "privkey.pem")
	issuerPath := filepath.Join(baseDir, "issuer.pem")

	if !force {
		fullchain, err1 := os.ReadFile(fullchainPath)
		_, err2 := os.ReadFile(privkeyPath)
		_, err3 := os.ReadFile(issuerPath)

		if err1 == nil && err2 == nil && err3 == nil {
			validUntil, err := ReadValidUntil(fullchain)
			if err != nil {
				return err
			}

			gotDomains, err := ExtractSANList(fullchain)
			if err != nil {
				return err
			}

			if !EqualStrings(gotDomains, wantDomains) {
				cfg.Sayf("certificate SANs changed, requesting new")
				force = true
			} else if time.Until(validUntil) <= 30*24*time.Hour {
				cfg.Sayf("certificate expires soon on %s, requesting new", validUntil.Format("2006-01-02"))
				force = true
			} else {
				cfg.Sayf("certificate valid until %s: %s",
					validUntil.Format("2006-01-02"),
					strings.Join(gotDomains, " "),
				)
				return nil
			}
		} else {
			force = true
		}
	}

	key, err := acme.GetPrivateKey(ACMEAccountKey)
	if err != nil {
		return err
	}

	var resource *certificate.Resource

	if cfg.HetznerToken != "" {
		os.Setenv("HETZNER_API_TOKEN", cfg.HetznerToken)
		defer os.Unsetenv("HETZNER_API_TOKEN")
		resource, err = acme.GetHetznerCertificate(wantDomains, cfg.SysAdmin, key)
	} else if cfg.IonosToken != "" {
		os.Setenv("IONOS_API_KEY", cfg.IonosToken)
		defer os.Unsetenv("IONOS_API_KEY")
		resource, err = acme.GetIonosCertificate(wantDomains, cfg.SysAdmin, key)
	}
	// add other providers here TODO Cloudflare

	if err != nil {
		return err
	}
	if resource == nil {
		return fmt.Errorf("missing provider for DNS-01 certificates")
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to mkdir %s: %w", baseDir, err)
	}

	write := func(name string, data []byte) error {
		path := filepath.Join(baseDir, name)
		return os.WriteFile(path, data, 0644)
	}

	if err := write("fullchain.pem", resource.Certificate); err != nil {
		return err
	}
	if err := write("privkey.pem", resource.PrivateKey); err != nil {
		return err
	}
	if err := write("issuer.pem", resource.IssuerCertificate); err != nil {
		return err
	}

	validUntil, err := ReadValidUntil(resource.Certificate)
	if err != nil {
		return err
	}

	gotDomains, err := ExtractSANList(resource.Certificate)
	if err != nil {
		return err
	}

	cfg.Sayf("certificate valid until %s: %s",
		validUntil.Format("2006-01-02"),
		strings.Join(gotDomains, " "),
	)

	return nil
}

func ExtractSANList(fullchain []byte) ([]string, error) {
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

	names := make([]string, 0, len(cert.DNSNames)+1)
	if cert.Subject.CommonName != "" {
		names = append(names, cert.Subject.CommonName)
	}
	names = append(names, cert.DNSNames...)

	return UniqueSortedStrings(names), nil
}

func ExtractSANs(fullchain []byte) (string, error) {
	names, err := ExtractSANList(fullchain)
	if err != nil {
		return "", err
	}
	return strings.Join(names, " "), nil
}

func UniqueSortedStrings(in []string) []string {
	m := make(map[string]struct{}, len(in))
	for _, s := range in {
		if s == "" {
			continue
		}
		m[s] = struct{}{}
	}

	out := make([]string, 0, len(m))
	for s := range m {
		out = append(out, s)
	}

	sort.Strings(out)
	return out
}

func EqualStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

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
