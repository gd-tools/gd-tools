package platform

import "fmt"

type Platform struct {
	// Baselines are embedded in the gdt binary and describe
	// supported runtime combinations such as Ubuntu, PHP, and Dovecot.
	Baselines []Baseline `json:"baselines"`

	// Products are embedded in the gdt binary and describe
	// supported applications managed by gdt.
	Products []Product `json:"products"`

	// Paths describe the concrete runtime directory tree.
	// They may be overridden in tests.
	Paths []Path `json:"paths"`
}

// LoadPlatform returns the default runtime platform.
func LoadPlatform() (*Platform, error) {
	pf := &Platform{
		Baselines: DefaultBaselines(),
		Products:  DefaultProducts(),
		Paths:     DefaultPaths(),
	}

	if err := pf.Validate(); err != nil {
		return nil, err
	}

	return pf, nil
}

// LoadPlatformWithPaths returns the default runtime platform
// with overridden runtime paths. This is primarily useful for tests.
func LoadPlatformWithPaths(paths []Path) (*Platform, error) {
	pf, err := LoadPlatform()
	if err != nil {
		return nil, err
	}

	pf.Paths = ClonePaths(paths)

	if err := pf.Validate(); err != nil {
		return nil, err
	}

	return pf, nil
}

// Validate checks the platform for required runtime information.
func (pf *Platform) Validate() error {
	if pf == nil {
		return fmt.Errorf("platform is nil")
	}
	if len(pf.Baselines) == 0 {
		return fmt.Errorf("missing baselines")
	}
	if len(pf.Products) == 0 {
		return fmt.Errorf("missing products")
	}
	if len(pf.Paths) == 0 {
		return fmt.Errorf("missing paths")
	}
	return nil
}
