package assets

import (
	"fmt"
	"sort"
)

// GetVersions returns all version keys for a product.
func GetVersions(product string) ([]string, error) {
	cat, err := LoadCatalog()
	if err != nil {
		return nil, err
	}

	prod, _, err := cat.GetProduct(product, "")
	if err != nil {
		return nil, err
	}

	var list []string

	for _, v := range prod.Versions {
		if v.Number != "" {
			list = append(list, v.Number)
		}
	}

	sort.Strings(list)

	return list, nil
}

// GetDefaultVersion returns the default version for a product.
func GetDefaultVersion(product string) (string, error) {
	cat, err := LoadCatalog()
	if err != nil {
		return "", err
	}

	prod, _, err := cat.GetProduct(product, "")
	if err != nil {
		return "", err
	}

	if prod.Default == "" {
		return "", fmt.Errorf("no default version for %s", product)
	}

	return prod.Default, nil
}
