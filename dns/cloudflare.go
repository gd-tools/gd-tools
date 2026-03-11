package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	CloudflareBaseURL = "https://api.cloudflare.com/client/v4/zones"
)

type CloudflareZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CloudflareResponse struct {
	Success bool             `json:"success"`
	Result  []CloudflareZone `json:"result"`
}

type CloudflareProvider struct {
	AuthAPIToken string
	Client       *http.Client
}

func NewCloudflareProvider(apiToken string) (*CloudflareProvider, error) {
	if strings.TrimSpace(apiToken) == "" {
		return nil, fmt.Errorf("cloudflare: api token required")
	}

	return &CloudflareProvider{
		AuthAPIToken: apiToken,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (p *CloudflareProvider) UpsertRRSet(ctx context.Context, zone string, rrIn RRSet) (string, error) {
	return "Cloudflare TODO", nil
}

func (p *CloudflareProvider) getZoneID(ctx context.Context, name string) (string, error) {
	url := CloudflareBaseURL
	data, err := p.sendToAPI(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	var zones []CloudflareZone
	if err := json.Unmarshal(data, &zones); err != nil {
		return "", err
	}

	for _, zone := range zones {
		if name == zone.Name {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("ionos: ID not found for zone %s", name)
}

func (p *CloudflareProvider) sendToAPI(ctx context.Context, method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.AuthAPIToken)

	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Keep body because Cloudflare returns useful JSON/text errors.
		return nil, fmt.Errorf("ionos: HTTP %d %s: %s", resp.StatusCode, resp.Status, strings.TrimSpace(string(data)))
	}

	return data, nil
}
