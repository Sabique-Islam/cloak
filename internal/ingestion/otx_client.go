package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// OTXClient handles AlienVault OTX API queries
type OTXClient struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

// OTXPulseResponse represents a pulse from OTX
type OTXPulseResponse struct {
	Results  []OTXPulse `json:"results"`
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
}

// OTXPulse represents a threat intelligence pulse
type OTXPulse struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	AuthorName        string         `json:"author_name"`
	Modified          string         `json:"modified"`
	Created           string         `json:"created"`
	Tags              []string       `json:"tags"`
	References        []string       `json:"references"`
	TLP               string         `json:"TLP"`
	Indicators        []OTXIndicator `json:"indicators"`
	Industries        []string       `json:"industries"`
	TargetedCountries []string       `json:"targeted_countries"`
	MalwareFamilies   []string       `json:"malware_families"`
	AttackIDs         []string       `json:"attack_ids"`
}

// OTXIndicator represents an indicator of compromise
type OTXIndicator struct {
	ID          int64  `json:"id"`
	Indicator   string `json:"indicator"`
	Type        string `json:"type"` // IPv4, domain, URL, FileHash-SHA256, etc.
	Created     string `json:"created"`
	Content     string `json:"content"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// OTXIPResponse represents IP reputation data
type OTXIPResponse struct {
	Reputation int        `json:"reputation"`
	Country    string     `json:"country_code"`
	City       string     `json:"city"`
	ASN        string     `json:"asn"`
	Latitude   float64    `json:"latitude"`
	Longitude  float64    `json:"longitude"`
	PulseCount int        `json:"pulse_info.count"`
	Pulses     []OTXPulse `json:"pulse_info.pulses"`
}

// OTXDomainResponse represents domain reputation data
type OTXDomainResponse struct {
	Indicator  string     `json:"indicator"`
	Type       string     `json:"type"`
	TypeTitle  string     `json:"type_title"`
	Alexa      string     `json:"alexa"`
	PulseCount int        `json:"pulse_info.count"`
	Pulses     []OTXPulse `json:"pulse_info.pulses"`
}

// OTXSearchResponse represents search results
type OTXSearchResponse struct {
	Results []OTXSearchResult `json:"results"`
	Count   int               `json:"count"`
}

// OTXSearchResult represents a single search result
type OTXSearchResult struct {
	Indicator string `json:"indicator"`
	Type      string `json:"type"`
	Pulse     string `json:"pulse"`
}

// NewOTXClient creates a new AlienVault OTX API client
func NewOTXClient(apiKey string) *OTXClient {
	client := resty.New()
	client.SetBaseURL("https://otx.alienvault.com/api/v1")
	client.SetHeader("X-OTX-API-KEY", apiKey)
	client.SetHeader("Accept", "application/json")
	client.SetTimeout(30 * time.Second)

	return &OTXClient{
		client:  client,
		apiKey:  apiKey,
		baseURL: "https://otx.alienvault.com/api/v1",
	}
}

// GetPulses retrieves recent threat intelligence pulses
func (o *OTXClient) GetPulses(ctx context.Context, modified_since string) (*OTXPulseResponse, error) {
	var result OTXPulseResponse

	req := o.client.R().
		SetContext(ctx).
		SetResult(&result)

	if modified_since != "" {
		req.SetQueryParam("modified_since", modified_since)
	}

	resp, err := req.Get("/pulses/subscribed")

	if err != nil {
		return nil, fmt.Errorf("otx get pulses request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("otx get pulses returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// GetIPReputation retrieves reputation data for an IP address
func (o *OTXClient) GetIPReputation(ctx context.Context, ip string) (*OTXIPResponse, error) {
	var result OTXIPResponse

	resp, err := o.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/indicators/IPv4/%s/general", ip))

	if err != nil {
		return nil, fmt.Errorf("otx ip reputation request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("otx ip reputation returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// GetDomainReputation retrieves reputation data for a domain
func (o *OTXClient) GetDomainReputation(ctx context.Context, domain string) (*OTXDomainResponse, error) {
	var result OTXDomainResponse

	resp, err := o.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/indicators/domain/%s/general", domain))

	if err != nil {
		return nil, fmt.Errorf("otx domain reputation request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("otx domain reputation returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// SearchIndicators searches for indicators across OTX
func (o *OTXClient) SearchIndicators(ctx context.Context, query string) (*OTXSearchResponse, error) {
	var result OTXSearchResponse

	resp, err := o.client.R().
		SetContext(ctx).
		SetQueryParam("q", query).
		SetResult(&result).
		Get("/search/pulses")

	if err != nil {
		return nil, fmt.Errorf("otx search request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("otx search returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// ToJSON converts the pulse response to JSON string
func (o *OTXPulseResponse) ToJSON() (string, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return "", fmt.Errorf("failed to marshal otx pulse response: %w", err)
	}
	return string(data), nil
}
