package platform

import (
	"fmt"
)

const (
	DefaultBaseline = "noble-8.3-2.4"
)

// Platform describes the runtime environment of gd-tools.
type Platform struct {
	// Baselines describe the available platforms.
	// This contains Ubuntu, PHP and Dovecot Versions.
	Baselines []Baseline

	// Baseline points to the version currently in use.
	Baseline *Baseline `json:"baseline"`

	// options describe paths and functions that can be faked for test.
	options *Options

	// Products describe the downloadable archives and binaries.
	Products []Product `json:"products"`
}

// LoadPlatform loads the platform definition.
func LoadPlatform(name string, opts *Options) (*Platform, error) {
	var pf Platform

	if opts != nil {
		pf.options = opts
	} else {
		pf.options = defaultOptions()
	}

	if err := pf.LoadBaselines(name); err != nil {
		return nil, err
	}

	if err := pf.LoadProducts(); err != nil {
		return nil, err
	}

	if err := pf.Validate(); err != nil {
		return nil, err
	}

	return &pf, nil
}

func (pf *Platform) Validate() error {
	if len(pf.Baselines) == 0 {
		return fmt.Errorf("platform has no baselines")
	}
	if pf.Baseline == nil {
		return fmt.Errorf("platform has no baseline pointer")
	}

	if len(pf.Products) == 0 {
		return fmt.Errorf("platform has no products")
	}

	if pf.options == nil {
		return fmt.Errorf("platform has no options")
	}
	if pf.options.rootDir == "" {
		return fmt.Errorf("platform has no rootDir")
	}
	if pf.options.varDir == "" {
		return fmt.Errorf("platform has no varDir")
	}
	if pf.options.etcDir == "" {
		return fmt.Errorf("platform has no etcDir")
	}
	if pf.options.binDir == "" {
		return fmt.Errorf("platform has no binDir")
	}
	if pf.options.runDir == "" {
		return fmt.Errorf("platform has no runDir")
	}

	return nil
}
