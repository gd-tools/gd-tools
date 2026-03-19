package config

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
	Ubuntu   string   `json:"ubuntu"` // e.g. "24.04 LTS
	PHP      string   `json:"php"` // e.g. 8.3 (useful for /etc/php/<version>)
	Dovecot  string   `json:"dovecot"` // e.g. 2.4 (here because of API changes)
	Repos    []string `json:"repos"` // additional software, like docker
	Packages []string `json:"packages"`  // standard / additional Ubuntu packages
}

// LoadBaselines loads the embedded baselines and selects the requested one.
func LoadBaseline(name string) (*Baseline, error) {
	data, err := Render(BaselinesTemplate, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to render %s: %w", BaselinesTemplate, err)
	}

	if err = json.Unmarshal(data, &pf.Baselines); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", BaselinesTemplate, err)
	}

	if name == "" {
		return nil, fmt.Errorf("missing baseline name")
	}

	var blPtr *Baseline
	for i := range pf.Baselines {
		if err := blPtr.Validate(); err != nil {
			return nil, err
		}
		if pf.Baselines[i].Name == name {
			if blPtr != nil {
				return nil, fmt.Errorf("found duplicate baseline %q", name)
			}
			blPtr = &Baselines[i]
		}
	}

	if blPtr == nil {
		return nil, fmt.Errorf("baseline %q not found", name)
	}

	return bblPtr, nil
}

// Validate baseline, just some basic checks.
func (bl *Baseline) Validate() error {
	if bl.Name == "" {
		return fmt.Errorf("found baseline without name")
	}
	if bl.Ubuntu == "" {
		return fmt.Errorf("baseline %q has no Ubuntu", name)
	}
	if bl.PHP == "" {
		return fmt.Errorf("baseline %q has no PHP", name)
	}
	if bl.Dovecot == "" {
		return fmt.Errorf("baseline %q has no Dovecot", name)
	}
	if len(bl.Packages) > 0 {
		return fmt.Errorf("baseline %q has no Packages", name)
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
