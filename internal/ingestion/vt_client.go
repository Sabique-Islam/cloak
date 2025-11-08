package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// VTClient handles VirusTotal API queries
type VTClient struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

// VTResponse represents a generic VirusTotal API response
type VTResponse struct {
	Data VTData                 `json:"data"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// VTData represents the data portion of VT response
type VTData struct {
	ID         string       `json:"id"`
	Type       string       `json:"type"`
	Attributes VTAttributes `json:"attributes"`
	Links      VTLinks      `json:"links"`
}

// VTAttributes contains various attributes depending on object type
type VTAttributes struct {
	// Common fields
	LastAnalysisStats    VTAnalysisStats     `json:"last_analysis_stats"`
	LastAnalysisResults  map[string]VTResult `json:"last_analysis_results"`
	LastAnalysisDate     int64               `json:"last_analysis_date"`
	LastModificationDate int64               `json:"last_modification_date"`

	// Domain/IP specific
	Reputation int               `json:"reputation"`
	Categories map[string]string `json:"categories"`
	TotalVotes VTVotes           `json:"total_votes"`

	// File specific
	SHA256          string   `json:"sha256"`
	SHA1            string   `json:"sha1"`
	MD5             string   `json:"md5"`
	Size            int64    `json:"size"`
	TypeDescription string   `json:"type_description"`
	Names           []string `json:"names"`

	// URL specific
	URL      string `json:"url"`
	Title    string `json:"title"`
	FinalURL string `json:"final_url"`

	// Additional metadata
	Tags    []string `json:"tags"`
	Whois   string   `json:"whois"`
	ASN     int      `json:"asn"`
	ASOwner string   `json:"as_owner"`
	Country string   `json:"country"`
}

// VTAnalysisStats contains detection statistics
type VTAnalysisStats struct {
	Malicious  int `json:"malicious"`
	Suspicious int `json:"suspicious"`
	Undetected int `json:"undetected"`
	Harmless   int `json:"harmless"`
	Timeout    int `json:"timeout"`
}

// VTResult represents a single vendor's analysis result
type VTResult struct {
	Category   string `json:"category"`
	Result     string `json:"result"`
	Method     string `json:"method"`
	EngineName string `json:"engine_name"`
}

// VTVotes represents community votes
type VTVotes struct {
	Harmless  int `json:"harmless"`
	Malicious int `json:"malicious"`
}

// VTLinks contains related API links
type VTLinks struct {
	Self string `json:"self"`
}

// VTSearchResponse represents search results
type VTSearchResponse struct {
	Data  []VTData          `json:"data"`
	Meta  VTSearchMeta      `json:"meta"`
	Links map[string]string `json:"links"`
}

// VTSearchMeta contains search metadata
type VTSearchMeta struct {
	Count int `json:"count"`
}

// NewVTClient creates a new VirusTotal API client
func NewVTClient(apiKey string) *VTClient {
	client := resty.New()
	client.SetBaseURL("https://www.virustotal.com/api/v3")
	client.SetHeader("x-apikey", apiKey)
	client.SetHeader("Accept", "application/json")
	client.SetTimeout(30 * time.Second)

	return &VTClient{
		client:  client,
		apiKey:  apiKey,
		baseURL: "https://www.virustotal.com/api/v3",
	}
}

// GetIPReport retrieves analysis report for an IP address
func (v *VTClient) GetIPReport(ctx context.Context, ip string) (*VTResponse, error) {
	var result VTResponse

	resp, err := v.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/ip_addresses/%s", ip))

	if err != nil {
		return nil, fmt.Errorf("virustotal ip report request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("virustotal ip report returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// GetDomainReport retrieves analysis report for a domain
func (v *VTClient) GetDomainReport(ctx context.Context, domain string) (*VTResponse, error) {
	var result VTResponse

	resp, err := v.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/domains/%s", domain))

	if err != nil {
		return nil, fmt.Errorf("virustotal domain report request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("virustotal domain report returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// GetURLReport retrieves analysis report for a URL
func (v *VTClient) GetURLReport(ctx context.Context, urlID string) (*VTResponse, error) {
	var result VTResponse

	resp, err := v.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/urls/%s", urlID))

	if err != nil {
		return nil, fmt.Errorf("virustotal url report request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("virustotal url report returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// GetFileReport retrieves analysis report for a file hash
func (v *VTClient) GetFileReport(ctx context.Context, hash string) (*VTResponse, error) {
	var result VTResponse

	resp, err := v.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/files/%s", hash))

	if err != nil {
		return nil, fmt.Errorf("virustotal file report request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("virustotal file report returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// Search performs a search query using VT Intelligence
func (v *VTClient) Search(ctx context.Context, query string) (*VTSearchResponse, error) {
	var result VTSearchResponse

	resp, err := v.client.R().
		SetContext(ctx).
		SetQueryParam("query", query).
		SetResult(&result).
		Get("/search")

	if err != nil {
		return nil, fmt.Errorf("virustotal search request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("virustotal search returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// ToJSON converts the VT response to JSON string
func (v *VTResponse) ToJSON() (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal virustotal response: %w", err)
	}
	return string(data), nil
}

// IsMalicious checks if the resource is flagged as malicious by any vendor
func (v *VTResponse) IsMalicious() bool {
	return v.Data.Attributes.LastAnalysisStats.Malicious > 0
}

// GetDetectionRate returns the detection rate as a string (e.g., "5/70")
func (v *VTResponse) GetDetectionRate() string {
	stats := v.Data.Attributes.LastAnalysisStats
	total := stats.Malicious + stats.Suspicious + stats.Undetected + stats.Harmless
	detected := stats.Malicious + stats.Suspicious
	return fmt.Sprintf("%d/%d", detected, total)
}
