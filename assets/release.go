package assets

import (
	"encoding/json"
	"fmt"

	"github.com/gd-tools/gd-tools/utils"
)

const ReleasesFile = "releases.json"

// Download describes a specific loadable asset, e.g. a zip archive or binary.
type Download struct {
	DownloadURL string `json:"download_url"`
	Filename    string `json:"filename"`
	Directory   string `json:"directory"`
	Binary      string `json:"binary"`
	MD5         string `json:"md5"`
	SHA256      string `json:"sha256"`
	SHA512      string `json:"sha512"`
}

// Baseline describes the platform runtime environment (for recovery).
type Baseline struct {
	Name     string   `json:"name"` // e.g. "noble-8.3-2.4"
	Ubuntu   string   `json:"ubuntu"`
	PHP      string   `json:"php"`
	Dovecot  string   `json:"dovecot"`
	Repos    []string `json:"repos"`
	Packages []string `json:"packages"`
}

// Release describes one specific release entry in the catalog.
type Release struct {
	Number   string   `json:"number"`
	Series   string   `json:"series,omitempty"`
	Download Download `json:"download"`
}

// Product describes a given product (like Nextcloud, WordPress, etc.).
type Product struct {
	Name      string    `json:"name"`
	SourceURL string    `json:"source_url,omitempty"`
	Default   string    `json:"default"`
	Directory string    `json:"directory,omitempty"`
	Binary    string    `json:"binary,omitempty"`
	Versions  []Release `json:"versions"`
}

// Catalog is the root structure of releases.json.
type Catalog struct {
	Baselines []Baseline `json:"baselines"`
	Products  []Product  `json:"products"`
}

// LoadCatalog reads the release catalog from assets/templates/system/releases.json.
func LoadCatalog() (*Catalog, error) {
	data, err := Render("system/"+ReleasesFile, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", ReleasesFile, err)
	}

	var catalog Catalog
	err = json.Unmarshal(data, &catalog)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", ReleasesFile, err)
	}

	return &catalog, nil
}

// GetBaseline provides the entry point for the system generation.
// It is used for e.g. loading packages or finding the php-fpm pool.
func (c *Catalog) GetBaseline(name string) (*Baseline, error) {
	for i := range c.Baselines {
		if c.Baselines[i].Name == name {
			return &c.Baselines[i], nil
		}
	}
	return nil, fmt.Errorf("baseline %q not found", name)
}

// GetProduct returns a specific release for one product.
// Return the default release if num is empty.
func (c *Catalog) GetProduct(name, num string) (*Product, *Release, error) {
	for i := range c.Products {
		pr := &c.Products[i]
		if pr.Name != name {
			continue
		}

		rel, err := pr.GetRelease(num)
		if err != nil {
			return nil, nil, err
		}

		return pr, rel, nil
	}

	return nil, nil, fmt.Errorf("product %q not found", name)
}

// GetRelease returns a specific product release.
// Return the default release if num is empty.
func (pr *Product) GetRelease(num string) (*Release, error) {
	if num == "" {
		num = pr.Default
	}

	for i := range pr.Versions {
		rel := pr.Versions[i]
		if rel.Number != num {
			continue
		}

		// copy release so we don't mutate the catalog
		r := rel

		if r.Download.Directory == "" {
			r.Download.Directory = pr.Directory
		}

		if r.Download.Binary == "" {
			r.Download.Binary = pr.Binary
		}

		return &r, nil
	}

	return nil, fmt.Errorf("release %q not found for product %q", num, pr.Name)
}

// Info returns formatted information for the given product.
func (pr *Product) Info() []string {
	var lb utils.LineBuffer

	lb.Addf("Product:    %s", pr.Name)
	if pr.SourceURL != "" {
		lb.Addf("Source URL: %s", pr.SourceURL)
	}

	lb.Add("Known releases:")

	for _, rel := range pr.Versions {
		line := rel.Number
		if rel.Number == pr.Default {
			line += " (default)"
		}

		lb.Addf("  - %s", line)
		lb.Addf("    Version:    %s", rel.Number)

		if rel.Series != "" {
			lb.Addf("    Series:     %s", rel.Series)
		}

		lb.Addf("    File:       %s", rel.Download.Filename)

		dir := rel.Download.Directory
		if dir == "" {
			dir = pr.Directory
		}
		if dir != "" {
			lb.Addf("    Directory:  %s", dir)
		}

		bin := rel.Download.Binary
		if bin == "" {
			bin = pr.Binary
		}
		if bin != "" {
			lb.Addf("    Binary:     %s", bin)
		}

		if rel.Download.MD5 != "" {
			lb.Addf("    MD5:        %s", rel.Download.MD5)
		}
		if rel.Download.SHA256 != "" {
			lb.Addf("    SHA256:     %s", rel.Download.SHA256)
		}
		if rel.Download.SHA512 != "" {
			lb.Addf("    SHA512:     %s", rel.Download.SHA512)
		}

		lb.Addf("    URL:        %s", rel.Download.DownloadURL)
	}

	return lb.Lines()
}
