package protocol

import "testing"

func TestRequestAddFirewall(t *testing.T) {
	req := &Request{}

	req.AddFirewall("allow 80/tcp")

	if len(req.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(req.Rules))
	}
	if req.Rules[0] != "allow 80/tcp" {
		t.Fatalf("unexpected rule value: %q", req.Rules[0])
	}
}

func TestRequestAddFirewallDedup(t *testing.T) {
	req := &Request{}

	req.AddFirewall("allow 80/tcp")
	req.AddFirewall("allow 80/tcp")

	if len(req.Rules) != 1 {
		t.Fatalf("expected deduplicated rule list, got %d", len(req.Rules))
	}
}

func TestRequestAddFirewallTrim(t *testing.T) {
	req := &Request{}

	req.AddFirewall("  allow 80/tcp  ")

	if len(req.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(req.Rules))
	}
	if req.Rules[0] != "allow 80/tcp" {
		t.Fatalf("expected trimmed rule, got %q", req.Rules[0])
	}
}

func TestRequestAddFirewallIgnoreEmpty(t *testing.T) {
	req := &Request{}

	req.AddFirewall("")

	if len(req.Rules) != 0 {
		t.Fatalf("expected no rules, got %d", len(req.Rules))
	}
}

func TestRequestAddFirewallNilReceiver(t *testing.T) {
	var req *Request

	req.AddFirewall("allow 22/tcp")
}

func TestRequestHasFirewallList(t *testing.T) {
	if (&Request{}).HasFirewallList() {
		t.Fatalf("expected false for empty request")
	}

	req := &Request{}
	req.AddFirewall("allow 22/tcp")

	if !req.HasFirewallList() {
		t.Fatalf("expected true after adding rule")
	}
}
