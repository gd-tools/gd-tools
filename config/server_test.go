package config

import "testing"

func TestServerFQDN(t *testing.T) {
	s := &Server{
		HostName:   "host",
		DomainName: "example.org",
	}

	if got := s.FQDN(); got != "host.example.org" {
		t.Fatalf("got %q", got)
	}
}

func TestServerFQDNVariants(t *testing.T) {
	s := &Server{
		HostName:   "host",
		DomainName: "example.org",
	}

	if s.FQDNdot() != "host.example.org." {
		t.Fatal("FQDNdot failed")
	}
	if s.DotFQDN() != ".host.example.org" {
		t.Fatal("DotFQDN failed")
	}
}

func TestServerRootUser(t *testing.T) {
	s := &Server{
		HostName:   "host",
		DomainName: "example.org",
	}

	if got := s.RootUser(); got != "root@host.example.org" {
		t.Fatalf("got %q", got)
	}
}

func TestServerLocale(t *testing.T) {
	s := &Server{}
	if s.Locale() != "" {
		t.Fatal("expected empty locale")
	}

	s.Language = "de"
	if s.Locale() != "de" {
		t.Fatal("language only failed")
	}

	s.Region = "DE"
	if s.Locale() != "de_DE" {
		t.Fatal("full locale failed")
	}
}
