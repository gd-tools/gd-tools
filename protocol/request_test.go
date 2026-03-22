package protocol

import (
	"strings"
	"testing"
)

func TestRequestStringNil(t *testing.T) {
	var req *Request
	if got := req.String(); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestRequestStringContainsFields(t *testing.T) {
	req := &Request{
		Version: 1,
		Hello:   Hello{Greeting: "hello"},
		Bootstrap: Bootstrap{
			FQDN:     "host.example.org",
			TimeZone: "Europe/Berlin",
		},
	}

	got := req.String()

	if !strings.Contains(got, `"version": 1`) {
		t.Fatalf("expected version in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"greeting": "hello"`) {
		t.Fatalf("expected greeting in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"fqdn": "host.example.org"`) {
		t.Fatalf("expected fqdn in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"time_zone": "Europe/Berlin"`) {
		t.Fatalf("expected time_zone in output, got:\n%s", got)
	}
}

func TestMakeFQDN(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		domain string
		want   string
	}{
		{
			name:   "host and domain",
			host:   "cloud",
			domain: "example.org",
			want:   "cloud.example.org",
		},
		{
			name:   "host only",
			host:   "cloud",
			domain: "",
			want:   "cloud",
		},
		{
			name:   "domain only",
			host:   "",
			domain: "example.org",
			want:   "example.org",
		},
		{
			name:   "both empty",
			host:   "",
			domain: "",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeFQDN(tt.host, tt.domain)
			if got != tt.want {
				t.Fatalf("MakeFQDN(%q, %q) = %q, want %q", tt.host, tt.domain, got, tt.want)
			}
		})
	}
}
