package config

import (
	"fmt"

	"github.com/gd-tools/gd-tools/utils"
)

const (
	BaselinesTemplate = "system/baselines.json"
)

// Baseline describes the desired system baseline.
type Baseline struct {
	Name     string   `json:"name"`     // e.g. "noble-8.3-2.4"
	Ubuntu   string   `json:"ubuntu"`   // e.g. "24.04 LTS"
	PHP      string   `json:"php"`      // e.g. "8.3"
	Dovecot  string   `json:"dovecot"`  // e.g. "2.4"
	Repos    []string `json:"repos"`    // additional software repositories
	Packages []string `json:"packages"` // standard and additional Ubuntu packages
}

// loadBaselines loads and validates the embedded baselines.
func loadBaselines() ([]Baseline, error) {
	var baselines []Baseline

	if err := RenderJSON(BaselinesTemplate, nil, &baselines); err != nil {
		return nil, fmt.Errorf("load baselines: %w", err)
	}

	seen := make(map[string]struct{}, len(baselines))

	for i := range baselines {
		if err := baselines[i].Validate(); err != nil {
			return nil, err
		}

		if _, ok := seen[baselines[i].Name]; ok {
			return nil, fmt.Errorf("found duplicate baseline %q", baselines[i].Name)
		}
		seen[baselines[i].Name] = struct{}{}
	}

	return baselines, nil
}

// LoadBaseline loads the embedded baselines and selects the requested one.
func LoadBaseline(name string) (*Baseline, error) {
	if name == "" {
		return nil, fmt.Errorf("missing baseline name")
	}

	baselines, err := loadBaselines()
	if err != nil {
		return nil, err
	}

	for i := range baselines {
		if baselines[i].Name == name {
			return &baselines[i], nil
		}
	}

	return nil, fmt.Errorf("baseline %q not found", name)
}

// Validate baseline, just some basic checks.
func (bl *Baseline) Validate() error {
	if bl == nil {
		return fmt.Errorf("baseline is nil")
	}
	if bl.Name == "" {
		return fmt.Errorf("found baseline without name")
	}
	if bl.Ubuntu == "" {
		return fmt.Errorf("baseline %q has no Ubuntu", bl.Name)
	}
	if bl.PHP == "" {
		return fmt.Errorf("baseline %q has no PHP", bl.Name)
	}
	if bl.Dovecot == "" {
		return fmt.Errorf("baseline %q has no Dovecot", bl.Name)
	}
	if len(bl.Packages) == 0 {
		return fmt.Errorf("baseline %q has no packages", bl.Name)
	}

	return nil
}

// Info returns formatted information for the baseline.
func (bl *Baseline) Info() []string {
	if bl == nil {
		return nil
	}

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
