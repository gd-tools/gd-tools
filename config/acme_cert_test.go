package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestUniqueSortedStrings(t *testing.T) {
	got := UniqueSortedStrings([]string{
		"mail.example.org",
		"",
		"example.org",
		"mail.example.org",
		"smtp.example.org",
	})

	want := []string{
		"example.org",
		"mail.example.org",
		"smtp.example.org",
	}

	if !EqualStrings(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func TestEqualStrings(t *testing.T) {
	if !EqualStrings([]string{"a", "b"}, []string{"a", "b"}) {
		t.Fatal("expected equal slices")
	}

	if EqualStrings([]string{"a"}, []string{"a", "b"}) {
		t.Fatal("expected different slices by length")
	}

	if EqualStrings([]string{"a", "b"}, []string{"b", "a"}) {
		t.Fatal("expected different slices by order")
	}
}

func TestExtractSANList(t *testing.T) {
	pemBytes, notAfter := mustTestCertificatePEM(t,
		"example.org",
		[]string{"imap.example.org", "smtp.example.org"},
	)

	got, err := ExtractSANList(pemBytes)
	if err != nil {
		t.Fatalf("ExtractSANList() returned error: %v", err)
	}

	want := []string{
		"example.org",
		"imap.example.org",
		"smtp.example.org",
	}

	if !EqualStrings(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}

	validUntil, err := ReadValidUntil(pemBytes)
	if err != nil {
		t.Fatalf("ReadValidUntil() returned error: %v", err)
	}

	if !validUntil.Equal(notAfter.UTC().Round(0)) {
		t.Fatalf("got %v want %v", validUntil, notAfter.UTC().Round(0))
	}
}

func TestExtractSANListInvalidPEM(t *testing.T) {
	_, err := ExtractSANList([]byte("not a pem"))
	if err == nil {
		t.Fatal("expected error for invalid pem")
	}
}

func TestReadValidUntilInvalidPEM(t *testing.T) {
	_, err := ReadValidUntil([]byte("not a pem"))
	if err == nil {
		t.Fatal("expected error for invalid pem")
	}
}

func TestExtractSANs(t *testing.T) {
	pemBytes, _ := mustTestCertificatePEM(t,
		"example.org",
		[]string{"smtp.example.org", "imap.example.org"},
	)

	got, err := ExtractSANs(pemBytes)
	if err != nil {
		t.Fatalf("ExtractSANs() returned error: %v", err)
	}

	wantParts := []string{
		"example.org",
		"imap.example.org",
		"smtp.example.org",
	}

	for _, part := range wantParts {
		if !strings.Contains(got, part) {
			t.Fatalf("missing %q in %q", part, got)
		}
	}
}

func mustTestCertificatePEM(t *testing.T, commonName string, dnsNames []string) ([]byte, time.Time) {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey() failed: %v", err)
	}

	notAfter := time.Now().UTC().Add(90 * 24 * time.Hour).Round(0)

	tpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             time.Now().UTC().Add(-1 * time.Hour),
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("x509.CreateCertificate() failed: %v", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})

	return pemBytes, notAfter
}
