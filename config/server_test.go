package config

import "testing"

func TestServerFQDN(t *testing.T) {
	srv := &Server{
		HostName:   "host00",
		DomainName: "example.org",
	}

	if got := srv.FQDN(); got != "host00.example.org" {
		t.Fatalf("unexpected fqdn: %q", got)
	}
}

func TestServerFQDNHostOnly(t *testing.T) {
	srv := &Server{
		HostName: "host00",
	}

	if got := srv.FQDN(); got != "host00" {
		t.Fatalf("unexpected fqdn: %q", got)
	}
}

func TestServerFQDNDomainOnly(t *testing.T) {
	srv := &Server{
		DomainName: "example.org",
	}

	if got := srv.FQDN(); got != "example.org" {
		t.Fatalf("unexpected fqdn: %q", got)
	}
}

func TestServerFQDNNil(t *testing.T) {
	var srv *Server

	if got := srv.FQDN(); got != "" {
		t.Fatalf("unexpected fqdn: %q", got)
	}
}
