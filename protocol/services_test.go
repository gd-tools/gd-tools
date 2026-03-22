package protocol

import "testing"

func TestRequestAddService(t *testing.T) {
	req := &Request{}

	req.AddService("nginx")

	if len(req.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(req.Services))
	}
	if req.Services[0] != "nginx" {
		t.Fatalf("unexpected service value: %q", req.Services[0])
	}
}

func TestRequestAddServiceDedup(t *testing.T) {
	req := &Request{}

	req.AddService("nginx")
	req.AddService("nginx")

	if len(req.Services) != 1 {
		t.Fatalf("expected deduplicated service list, got %d", len(req.Services))
	}
}

func TestRequestAddServiceTrim(t *testing.T) {
	req := &Request{}

	req.AddService("  nginx  ")

	if len(req.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(req.Services))
	}
	if req.Services[0] != "nginx" {
		t.Fatalf("expected trimmed service, got %q", req.Services[0])
	}
}

func TestRequestAddServiceIgnoreEmpty(t *testing.T) {
	req := &Request{}

	req.AddService("")

	if len(req.Services) != 0 {
		t.Fatalf("expected 0 services, got %d", len(req.Services))
	}
}

func TestRequestAddServiceNilReceiver(t *testing.T) {
	var req *Request
	req.AddService("nginx")
}

func TestRequestHasServiceList(t *testing.T) {
	if (&Request{}).HasServiceList() {
		t.Fatalf("expected false for empty request")
	}

	req := &Request{}
	req.AddService("nginx")

	if !req.HasServiceList() {
		t.Fatalf("expected true after adding service")
	}
}
