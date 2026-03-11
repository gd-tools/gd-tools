package releases

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/gd-tools/gd-tools/templates"
)

const (
	ReleasesFile = "releases.json"
)

// Load reads the release catalog from assets/releases.json.
func Load() (*Catalog, error) {
	releasesPath := filepath.Join("assets", ReleasesFile)
	data, err := templates.Load(releasesPath, false)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", ReleasesFile, err)
	}

	var catalog Catalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", ReleasesFile, err)
	}

	if err := catalog.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate %s: %w", ReleasesFile, err)
	}

	return &catalog, nil
}
