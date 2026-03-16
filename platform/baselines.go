package platform

import (
	"encoding/json"
	"fmt"

	"github.com/gd-tools/gd-tools/utils"
)

const (
	BaselinesTemplate = "system/baselines.json"
)

// Baseline describes the platform runtime environment.
type Baseline struct {
	Name     string   `json:"name"` // e.g. "noble-8.3-2.4"
	Ubuntu   string   `json:"ubuntu"`
	PHP      string   `json:"php"`
	Dovecot  string   `json:"dovecot"`
	Repos    []string `json:"repos"`
	Packages []string `json:"packages"`
}

// LoadBaselines loads the embedded baselines and selects the requested one.
func (pf *Platform) LoadBaselines(name string) error {
	data, err := Render(BaselinesTemplate, nil)
	if err != nil {
		return fmt.Errorf("failed to render %s: %w", BaselinesTemplate, err)
	}

	if err = json.Unmarshal(data, &pf.baselines); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", BaselinesTemplate, err)
	}

	if name == "" {
		return fmt.Errorf("missing baseline name")
	}
	for i := range pf.baselines {
		if pf.baselines[i].Name == name {
			pf.Baseline = &pf.baselines[i]
			break
		}
	}
	if pf.Baseline == nil {
		return fmt.Errorf("baseline %q not found", name)
	}

	return nil
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
