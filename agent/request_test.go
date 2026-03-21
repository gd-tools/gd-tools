package agent

import (
	"encoding/json"
	"testing"
)

func TestRequestJSON(t *testing.T) {
	req := Request{
		Version:  ProtocolVersion,
		Hello:    "test",
		FQDN:     "example.org",
		TimeZone: "Europe/Berlin",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Request
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.Hello != "test" {
		t.Fatalf("expected hello=test got %s", decoded.Hello)
	}

	if decoded.FQDN != "example.org" {
		t.Fatalf("expected fqdn=example.org got %s", decoded.FQDN)
	}
}

func TestResponseSay(t *testing.T) {
	resp := Response{}

	resp.Say("line1\nline2\nline3")

	if len(resp.Result) != 3 {
		t.Fatalf("expected 3 lines got %d", len(resp.Result))
	}

	if resp.Result[1] != "line2" {
		t.Fatalf("unexpected line: %s", resp.Result[1])
	}
}

func TestResponseSayf(t *testing.T) {
	resp := Response{}

	resp.Sayf("hello %s", "world")

	if len(resp.Result) != 1 {
		t.Fatalf("expected 1 line got %d", len(resp.Result))
	}

	if resp.Result[0] != "hello world" {
		t.Fatalf("unexpected result: %s", resp.Result[0])
	}
}

func TestResponseAddServiceDedup(t *testing.T) {
	resp := Response{}

	resp.AddService("nginx")
	resp.AddService("nginx")
	resp.AddService("redis")

	if len(resp.Services) != 2 {
		t.Fatalf("expected 2 services got %d", len(resp.Services))
	}
}

func TestRequestAddService(t *testing.T) {
	req := Request{}

	req.AddService("nginx")
	req.AddService("")

	if len(req.Services) != 1 {
		t.Fatalf("expected 1 service got %d", len(req.Services))
	}
}

func TestRequestAddFirewall(t *testing.T) {
	req := Request{}

	req.AddFirewall("80")
	req.AddFirewall("")

	if len(req.Firewall) != 1 {
		t.Fatalf("expected 1 firewall entry got %d", len(req.Firewall))
	}
}

func TestRequestString(t *testing.T) {
	req := Request{
		Version: ProtocolVersion,
		Hello:   "world",
	}

	s := req.String()

	if len(s) == 0 {
		t.Fatal("expected non-empty string")
	}
}

func TestResponseString(t *testing.T) {
	resp := Response{}
	resp.Say("hello")

	s := resp.String()

	if len(s) == 0 {
		t.Fatal("expected non-empty string")
	}
}
