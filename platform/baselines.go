package platform

import (
	"encoding/json"
	"fmt"

	"github.com/gd-tools/gd-tools/utils"
)

const BaselinesPath = "system/baselines.json"

// Baseline describes the platform runtime environment.
type Baseline struct {
	Name     string   `json:"name"` // e.g. "noble-8.3-2.4"
	Ubuntu   string   `json:"ubuntu"`
	PHP      string   `json:"php"`
	Dovecot  string   `json:"dovecot"`
	Repos    []string `json:"repos"`
	Packages []string `json:"packages"`
}

// DefaultBaselines returns the baselines embedded in the gdt binary.
func DefaultBaselines() []Baseline {
	data, err := Render(BaselinesPath, nil)
	if err != nil {
		return nil // error will be reported later
	}

	var baselines []Baseline
	err = json.Unmarshal(data, &baselines)
	if err != nil {
		return nil // error will be reported later
	}

	return baselines
}

// GetBaseline returns one baseline by name.
func (pf *Platform) GetBaseline(name string) (*Baseline, error) {
	for i := range pf.Baselines {
		if pf.Baselines[i].Name == name {
			return &pf.Baselines[i], nil
		}
	}
	return nil, fmt.Errorf("baseline %q not found", name)
}

// Info returns formatted information for the baseline.
func (bl *Baseline) Info() []string {
	var lb utils.LineBuffer

	lb.Addf("Baseline: %s", bl.Name)

	if bl.Ubuntu != "" {
		lb.Addf("Ubuntu:   %s", bl.Ubuntu)
	}
	if bl.PHP != "" {
		lb.Addf("PHP:      %s", bl.PHP)
	}
	if bl.Dovecot != "" {
		lb.Addf("Dovecot:  %s", bl.Dovecot)
	}

	if len(bl.Repos) > 0 {
		lb.Add("Repositories:")
		for _, repo := range bl.Repos {
			lb.Addf("  - %s", repo)
		}
	}

	if len(bl.Packages) > 0 {
		lb.Add("Packages:")
		for _, pkg := range bl.Packages {
			lb.Addf("  - %s", pkg)
		}
	}

	return lb.Lines()
}
