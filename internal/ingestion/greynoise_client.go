package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// GreyNoiseClient handles GreyNoise API queries
type GreyNoiseClient struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

// GreyNoiseIPResponse represents the IP lookup response
type GreyNoiseIPResponse struct {
	IP             string            `json:"ip"`
	Seen           bool              `json:"seen"`
	Classification string            `json:"classification"` // unknown, benign, malicious
	FirstSeen      string            `json:"first_seen"`
	LastSeen       string            `json:"last_seen"`
	Actor          string            `json:"actor"`
	Tags           []string          `json:"tags"`
	Metadata       GreyNoiseMetadata `json:"metadata"`
	RawData        GreyNoiseRawData  `json:"raw_data"`
}

// GreyNoiseMetadata contains metadata about the IP
type GreyNoiseMetadata struct {
	ASN          string `json:"asn"`
	City         string `json:"city"`
	Country      string `json:"country"`
	CountryCode  string `json:"country_code"`
	Organization string `json:"organization"`
	Category     string `json:"category"`
	TOR          bool   `json:"tor"`
	RDNS         string `json:"rdns"`
	OS           string `json:"os"`
}

// GreyNoiseRawData contains raw scanning data
type GreyNoiseRawData struct {
	Scan      []GreyNoiseScan `json:"scan"`
	JA3       []GreyNoiseJA3  `json:"ja3"`
	Hassh     []string        `json:"hassh"`
	UserAgent []string        `json:"user_agent"`
}

// GreyNoiseScan represents scanning activity
type GreyNoiseScan struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

// GreyNoiseJA3 represents JA3 fingerprints
type GreyNoiseJA3 struct {
	Fingerprint string `json:"fingerprint"`
	Port        int    `json:"port"`
}

// GreyNoiseQueryResponse represents the GNQL query response
type GreyNoiseQueryResponse struct {
	Complete bool                  `json:"complete"`
	Count    int                   `json:"count"`
	Data     []GreyNoiseIPResponse `json:"data"`
	Message  string                `json:"message"`
	Query    string                `json:"query"`
}

// NewGreyNoiseClient creates a new GreyNoise API client
func NewGreyNoiseClient(apiKey string) *GreyNoiseClient {
	client := resty.New()
	client.SetBaseURL("https://api.greynoise.io/v3")
	client.SetHeader("key", apiKey)
	client.SetHeader("Accept", "application/json")
	client.SetTimeout(30 * time.Second)

	return &GreyNoiseClient{
		client:  client,
		apiKey:  apiKey,
		baseURL: "https://api.greynoise.io/v3",
	}
}

// LookupIP performs an IP lookup to get reputation data
func (g *GreyNoiseClient) LookupIP(ctx context.Context, ip string) (*GreyNoiseIPResponse, error) {
	var result GreyNoiseIPResponse

	resp, err := g.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/community/%s", ip))

	if err != nil {
		return nil, fmt.Errorf("greynoise lookup request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("greynoise lookup returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// Query performs a GNQL query for advanced searching
func (g *GreyNoiseClient) Query(ctx context.Context, query string) (*GreyNoiseQueryResponse, error) {
	var result GreyNoiseQueryResponse

	resp, err := g.client.R().
		SetContext(ctx).
		SetQueryParam("query", query).
		SetQueryParam("size", "100").
		SetResult(&result).
		Get("/query")

	if err != nil {
		return nil, fmt.Errorf("greynoise query request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("greynoise query returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// QuickCheck performs a quick IP check (community API)
func (g *GreyNoiseClient) QuickCheck(ctx context.Context, ip string) (bool, error) {
	result, err := g.LookupIP(ctx, ip)
	if err != nil {
		return false, err
	}

	return result.Seen, nil
}

// ToJSON converts the IP response to JSON string
func (g *GreyNoiseIPResponse) ToJSON() (string, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return "", fmt.Errorf("failed to marshal greynoise response: %w", err)
	}
	return string(data), nil
}
