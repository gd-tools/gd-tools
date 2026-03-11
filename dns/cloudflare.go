package dns

import (
	"fmt"
)

const (
	CloudflareBaseURL = " https://api.cloudflare.com/client/v4/zones"
)

type CloudflareProvider struct {
	AuthAPIToken string
}

func Reference(token string) error {
	fmt.Println()
	return nil
}
