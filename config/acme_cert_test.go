package config

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-acme/lego/v4/certificate"
)

func sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func makeTestCertificatePEM(t *testing.T, commonName string, dnsNames []string, notAfter time.Time) []byte {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Fatalf("serial failed: %v", err)
	}

	tpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             time.Now().Add(-1 * time.Hour).UTC(),
		NotAfter:              notAfter.UTC(),
		DNSNames:              dnsNames,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("CreateCertificate failed: %v", err)
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		t.Fatalf("pem.Encode failed: %v", err)
	}

	return buf.Bytes()
}

func TestReadValidUntil(t *testing.T) {
	want := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	pemBytes := makeTestCertificatePEM(t, "example.org", []string{"example.org", "www.example.org"}, want)

	got, err := ReadValidUntil(pemBytes)
	if err != nil {
		t.Fatalf("ReadValidUntil returned error: %v", err)
	}

	if !got.Equal(want.Round(0)) {
		t.Fatalf("unexpected NotAfter: got %v want %v", got, want.Round(0))
	}
}

func TestReadValidUntilNoCertificate(t *testing.T) {
	_, err := ReadValidUntil([]byte("not a certificate"))
	if err == nil {
		t.Fatal("expected error for invalid pem data")
	}
}

func TestExtractSANList(t *testing.T) {
	pemBytes := makeTestCertificatePEM(
		t,
		"example.org",
		[]string{"www.example.org", "example.org", "api.example.org"},
		time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC),
	)

	got, err := ExtractSANList("example.org", pemBytes)
	if err != nil {
		t.Fatalf("ExtractSANList returned error: %v", err)
	}

	want := []string{"example.org", "api.example.org", "www.example.org"}
	if !reflect.DeepEqual(got.Lines(), want) {
		t.Fatalf("unexpected SAN list:\n got: %#v\nwant: %#v", got.Lines(), want)
	}
}

func TestExtractSANListAddsDomainIfMissing(t *testing.T) {
	pemBytes := makeTestCertificatePEM(
		t,
		"ignored-cn.example.org",
		[]string{"www.example.org"},
		time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC),
	)

	got, err := ExtractSANList("example.org", pemBytes)
	if err != nil {
		t.Fatalf("ExtractSANList returned error: %v", err)
	}

	want := []string{"example.org", "www.example.org"}
	if !reflect.DeepEqual(got.Lines(), want) {
		t.Fatalf("unexpected SAN list:\n got: %#v\nwant: %#v", got.Lines(), want)
	}
}

func TestExtractSANListInvalidPEM(t *testing.T) {
	_, err := ExtractSANList("example.org", []byte("broken"))
	if err == nil {
		t.Fatal("expected error for invalid pem data")
	}
}

func TestEnsureCertificateReusesExistingValidCertificate(t *testing.T) {
	domain := "example.org"
	fullchain := makeTestCertificatePEM(
		t,
		domain,
		[]string{domain, "www.example.org"},
		time.Now().Add(90*24*time.Hour).UTC(),
	)

	cfg := &Config{}

	loadCalls := 0
	cfg.loadFile = func(name string) ([]byte, error) {
		loadCalls++
		switch filepath.Base(name) {
		case "fullchain.pem":
			return fullchain, nil
		case "privkey.pem":
			return []byte("key"), nil
		case "issuer.pem":
			return []byte("issuer"), nil
		default:
			t.Fatalf("unexpected LoadFile path: %s", name)
			return nil, nil
		}
	}

	var infoLines []string
	cfg.infof = func(format string, args ...any) {
		infoLines = append(infoLines, sprintf(format, args...))
	}

	calledProvider := false
	cfg.getPrivateKey = func(name string) (any, error) {
		t.Fatalf("GetPrivateKey must not be called when cert is reused")
		return nil, nil
	}
	cfg.getHetznerCertificate = func(domains []string, email string, key any) (*certificate.Resource, error) {
		calledProvider = true
		return nil, nil
	}

	err := cfg.EnsureCertificate(domain, "www.example.org")
	if err != nil {
		t.Fatalf("EnsureCertificate returned error: %v", err)
	}

	if calledProvider {
		t.Fatal("provider must not be called for reusable certificate")
	}

	if loadCalls != 3 {
		t.Fatalf("unexpected LoadFile count: got %d want 3", loadCalls)
	}

	if len(infoLines) == 0 {
		t.Fatal("expected info log for reused certificate")
	}
}

func TestEnsureCertificateRequestsNewWhenSANChanged(t *testing.T) {
	domain := "example.org"
	fullchain := makeTestCertificatePEM(
		t,
		domain,
		[]string{domain, "old.example.org"},
		time.Now().Add(90*24*time.Hour).UTC(),
	)

	cfg := &Config{
		HetznerToken: "token",
		SysAdmin:     "admin@example.org",
	}

	cfg.loadFile = func(name string) ([]byte, error) {
		switch filepath.Base(name) {
		case "fullchain.pem":
			return fullchain, nil
		case "privkey.pem":
			return []byte("key"), nil
		case "issuer.pem":
			return []byte("issuer"), nil
		default:
			t.Fatalf("unexpected LoadFile path: %s", name)
			return nil, nil
		}
	}

	cfg.getPrivateKey = func(name string) (any, error) {
		return "private-key", nil
	}

	cfg.setenv = func(key, value string) error {
		if key != "HETZNER_API_TOKEN" || value != "token" {
			t.Fatalf("unexpected Setenv: %s=%s", key, value)
		}
		return nil
	}
	cfg.unsetenv = func(key string) error {
		if key != "HETZNER_API_TOKEN" {
			t.Fatalf("unexpected Unsetenv key: %s", key)
		}
		return nil
	}

	calledDomains := []string(nil)
	cfg.getHetznerCertificate = func(domains []string, email string, key any) (*certificate.Resource, error) {
		calledDomains = append([]string{}, domains...)

		return &certificate.Resource{
			Certificate:       makeTestCertificatePEM(t, domain, []string{domain, "www.example.org"}, time.Now().Add(90*24*time.Hour)),
			PrivateKey:        []byte("new-private-key"),
			IssuerCertificate: []byte("new-issuer"),
		}, nil
	}

	var mkdirPath string
	cfg.mkdirAll = func(path string, perm os.FileMode) error {
		mkdirPath = path
		return nil
	}

	saved := map[string][]byte{}
	cfg.saveFile = func(name string, data []byte) error {
		saved[name] = append([]byte{}, data...)
		return nil
	}

	err := cfg.EnsureCertificate(domain, "www.example.org")
	if err != nil {
		t.Fatalf("EnsureCertificate returned error: %v", err)
	}

	wantDomains := []string{"example.org", "www.example.org"}
	if !reflect.DeepEqual(calledDomains, wantDomains) {
		t.Fatalf("unexpected requested domains:\n got: %#v\nwant: %#v", calledDomains, wantDomains)
	}

	if mkdirPath != filepath.Join(ACMECertDir, domain) {
		t.Fatalf("unexpected mkdir path: %q", mkdirPath)
	}

	if len(saved) != 3 {
		t.Fatalf("expected 3 saved files, got %d", len(saved))
	}
}

func TestEnsureCertificateFailsWithoutProvider(t *testing.T) {
	domain := "example.org"
	fullchain := makeTestCertificatePEM(
		t,
		domain,
		[]string{domain},
		time.Now().Add(5*24*time.Hour).UTC(),
	)

	cfg := &Config{}

	cfg.loadFile = func(name string) ([]byte, error) {
		switch filepath.Base(name) {
		case "fullchain.pem":
			return fullchain, nil
		case "privkey.pem":
			return []byte("key"), nil
		case "issuer.pem":
			return []byte("issuer"), nil
		default:
			t.Fatalf("unexpected LoadFile path: %s", name)
			return nil, nil
		}
	}

	cfg.getPrivateKey = func(name string) (any, error) {
		return "private-key", nil
	}

	err := cfg.EnsureCertificate(domain)
	if err == nil {
		t.Fatal("expected error when no provider is configured")
	}

	if !strings.Contains(err.Error(), "missing provider") {
		t.Fatalf("unexpected error: %v", err)
	}
}
