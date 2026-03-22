package protocol

import (
	"strings"
)

// FirewallList contains firewall rules to be applied on prod.
type FirewallList struct {
	Rules []string `json:"firewall_rules,omitempty"`
}

// AddFirewall adds a firewall rule if it is non-empty and not already present (deduplicated).
func (req *Request) AddFirewall(rule string) {
	if req == nil {
		return
	}

	rule = strings.TrimSpace(rule)
	if rule == "" {
		return
	}

	for _, check := range req.Rules {
		if check == rule {
			return
		}
	}

	req.Rules = append(req.Rules, rule)
}

// HasFirewallList reports whether the request contains firewall rules.
func (req *Request) HasFirewallList() bool {
	if req == nil {
		return false
	}
	return len(req.Rules) > 0
}
