package protocol

type FirewallList struct {
	FirewallRules []string `json:"firewall_rules,omitempty"`
}

func (req *Request) AddFirewall(rule string) {
	if req == nil || rule == "" {
		return
	}
	for _, check := range req.FirewallRules {
		if check == rule {
			return
		}
	}
	req.FirewallRules = append(req.FirewallRules, rule)
}

func (req *Request) HasFirewallList() bool {
	if req == nil {
		return false
	}
	return len(req.FirewallRules) > 0
}
