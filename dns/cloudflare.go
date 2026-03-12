package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	CloudflareBaseURL = "https://api.cloudflare.com/client/v4"
)

type CloudflareRecord struct {
	Name  string `json:"name"`
	ID    string `json:"id,omitempty"`
	TTL   int    `json:"ttl"`
	Type  string `json:"type"`
	Value string `json:"content"`
	Prio  int    `json:"priority"`
	Text  string `json:"-"`
}

type CloudflareUpdate struct {
	Value    string `json:"content"`
	TTL      int    `json:"ttl"`
	Prio     int    `json:"prio"`
	Disabled bool   `json:"disabled"`
}

type CloudflareZone struct {
	Name    string             `json:"name"`
	ID      string             `json:"id"`
	Records []CloudflareRecord `json:"records"`
	rrSet   RRSet              `json:"-"`
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

func (p *CloudflareProvider) getZoneData(ctx context.Context, zone string, rr RRSet) (*CloudflareZone, error) {
	zone = normalizeZone(zone)

	zoneID, err := p.getZoneID(ctx, zone)
	if err != nil {
		return nil, err
	}

	var theZone CloudflareZone
	url := CloudflareBaseURL + "/zones/" + zoneID + "/dns_records"
	data, err := p.sendToAPI(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &theZone); err != nil {
		return nil, err
	}

	wantType := strings.ToUpper(string(rr.Type))
	wantName := rr.Name
	var records []CloudflareRecord
	for _, zoneRec := range theZone.Records {
		if strings.ToUpper(zoneRec.Type) != wantType {
			continue
		}
		if normalizeZone(zoneRec.Name) != wantName {
			continue
		}
		if zoneRec.Value = strings.TrimSpace(zoneRec.Value); zoneRec.Value == "" {
			continue
		}
		switch rr.Type {
		case RR_CNAME:
			host := strings.TrimSuffix(zoneRec.Value, ".")
			zoneRec.Value = host
			zoneRec.Text = host + "."
		case RR_MX:
			host := strings.TrimSuffix(zoneRec.Value, ".")
			zoneRec.Value = host
			zoneRec.Text = fmt.Sprintf("%d %s.", zoneRec.Prio, host)
		case RR_TXT:
			zoneRec.Text = canonTXTWire(zoneRec.Value)
		default:
			zoneRec.Text = zoneRec.Value
		}
		records = append(records, zoneRec)
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].Prio != records[j].Prio {
			return records[i].Prio < records[j].Prio
		}
		return records[i].Value < records[j].Value
	})

	theZone.Records = records // filtered by typ and name, already sorted
	theZone.rrSet = rr

	return &theZone, nil
}

func (p *CloudflareProvider) getZoneID(ctx context.Context, name string) (string, error) {
	url := CloudflareBaseURL + "/zones"
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

	return "", fmt.Errorf("cloudflare: ID not found for zone %s", name)
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
		return nil, fmt.Errorf("cloudflare: HTTP %d %s: %s", resp.StatusCode, resp.Status, strings.TrimSpace(string(data)))
	}

	return data, nil
}
