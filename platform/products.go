package platform

import (
	"encoding/json"
	"fmt"

	"github.com/gd-tools/gd-tools/utils"
)

const ProductsTemplate = "system/products.json"

// Release describes one specific release entry for a product.
type Release struct {
	Number   string   `json:"number"`
	Series   string   `json:"series,omitempty"`
	Download Download `json:"download"`
}

// Product describes a managed application such as Nextcloud or WordPress.
type Product struct {
	Name      string    `json:"name"`
	SourceURL string    `json:"source_url,omitempty"`
	Default   string    `json:"default"`
	Directory string    `json:"directory,omitempty"`
	Binary    string    `json:"binary,omitempty"`
	Marker    string    `json:"marker,omitempty"`
	Versions  []Release `json:"versions"`
}

// LoadProducts loads the products embedded in the gdt binary.
func (pf *Platform) LoadProducts() error {
	data, err := Render(ProductsTemplate, nil)
	if err != nil {
		return fmt.Errorf("failed to render %s: %w", ProductsTemplate, err)
	}

	if err := json.Unmarshal(data, &pf.Products); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", ProductsTemplate, err)
	}

	return nil
}

// FindProduct returns one product by name.
func (pf *Platform) FindProduct(name string) (*Product, error) {
	for i := range pf.Products {
		if pf.Products[i].Name == name {
			return &pf.Products[i], nil
		}
	}
	return nil, fmt.Errorf("product %q not found", name)
}

// GetProduct returns one product and one selected release.
// If num is empty, the default release is used.
func (pf *Platform) GetProduct(name, num string) (*Product, *Release, error) {
	pr, err := pf.FindProduct(name)
	if err != nil {
		return nil, nil, err
	}

	rel, err := pr.GetRelease(num)
	if err != nil {
		return nil, nil, err
	}

	return pr, rel, nil
}

// GetRelease returns one specific product release.
// If num is empty, the product default is used.
func (pr *Product) GetRelease(num string) (*Release, error) {
	if num == "" {
		num = pr.Default
	}

	for i := range pr.Versions {
		rel := pr.Versions[i]
		if rel.Number != num {
			continue
		}

		r := rel

		if r.Download.Directory == "" {
			r.Download.Directory = pr.Directory
		}

		if r.Download.Binary == "" {
			r.Download.Binary = pr.Binary
		}

		if r.Download.Marker == "" {
			r.Download.Marker = pr.Marker
		}

		return &r, nil
	}

	return nil, fmt.Errorf("release %q not found for product %q", num, pr.Name)
}

// DefaultRelease returns the default release of the product.
func (pr *Product) DefaultRelease() (*Release, error) {
	return pr.GetRelease(pr.Default)
}

// ReleaseCandidate describes one completion candidate for shell completion.
type ReleaseCandidate struct {
	Value       string
	Description string
	Default     bool
}

// ReleaseCandidates returns all selectable releases for shell completion or UI.
func (pr *Product) ReleaseCandidates() []ReleaseCandidate {
	out := make([]ReleaseCandidate, 0, len(pr.Versions))

	for _, rel := range pr.Versions {
		c := ReleaseCandidate{
			Value:       rel.Number,
			Description: rel.Number,
			Default:     rel.Number == pr.Default,
		}

		if c.Default {
			c.Description += " (default)"
		}

		if rel.Series != "" {
			c.Description += " [" + rel.Series + "]"
		}

		out = append(out, c)
	}

	return out
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

		if rel.Series != "" {
			line += " [" + rel.Series + "]"
		}

		lb.Addf("  - %s", line)

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
