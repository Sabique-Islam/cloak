package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// ShodanClient handles Shodan API queries
type ShodanClient struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

// ShodanSearchResponse represents search results from Shodan
type ShodanSearchResponse struct {
	Matches []ShodanHost `json:"matches"`
	Total   int          `json:"total"`
}

// ShodanHost represents a host discovered by Shodan
type ShodanHost struct {
	IP              string         `json:"ip_str"`
	Port            int            `json:"port"`
	Organization    string         `json:"org"`
	Hostnames       []string       `json:"hostnames"`
	Domains         []string       `json:"domains"`
	Location        ShodanLocation `json:"location"`
	Timestamp       string         `json:"timestamp"`
	ISP             string         `json:"isp"`
	ASN             string         `json:"asn"`
	Transport       string         `json:"transport"`
	Product         string         `json:"product"`
	Version         string         `json:"version"`
	CPE             []string       `json:"cpe"`
	OS              string         `json:"os"`
	Data            string         `json:"data"`
	Banner          string         `json:"banner"`
	SSL             *ShodanSSL     `json:"ssl,omitempty"`
	HTTP            *ShodanHTTP    `json:"http,omitempty"`
	Vulnerabilities []string       `json:"vulns,omitempty"`
}

// ShodanLocation represents geographical location
type ShodanLocation struct {
	City        string  `json:"city"`
	RegionCode  string  `json:"region_code"`
	AreaCode    int     `json:"area_code"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Country     string  `json:"country_name"`
	CountryCode string  `json:"country_code"`
}

// ShodanSSL represents SSL/TLS information
type ShodanSSL struct {
	Cert   ShodanCertificate `json:"cert"`
	Cipher ShodanCipher      `json:"cipher"`
	Chain  []string          `json:"chain"`
	ALPN   []string          `json:"alpn"`
	JA3S   string            `json:"ja3s"`
}

// ShodanCertificate represents SSL certificate details
type ShodanCertificate struct {
	Issued      string `json:"issued"`
	Expires     string `json:"expires"`
	Subject     string `json:"subject"`
	Issuer      string `json:"issuer"`
	Fingerprint string `json:"fingerprint"`
}

// ShodanCipher represents cipher information
type ShodanCipher struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Bits    int    `json:"bits"`
}

// ShodanHTTP represents HTTP-specific information
type ShodanHTTP struct {
	Status  int               `json:"status"`
	Title   string            `json:"title"`
	Server  string            `json:"server"`
	Headers map[string]string `json:"headers"`
	HTML    string            `json:"html"`
	Favicon *ShodanFavicon    `json:"favicon,omitempty"`
}

// ShodanFavicon represents website favicon information
type ShodanFavicon struct {
	Hash     string `json:"hash"`
	Data     string `json:"data"`
	Location string `json:"location"`
}

// ShodanHostInfo represents detailed host information
type ShodanHostInfo struct {
	IP              string       `json:"ip_str"`
	OS              string       `json:"os"`
	Organization    string       `json:"org"`
	Data            []ShodanHost `json:"data"`
	Ports           []int        `json:"ports"`
	Hostnames       []string     `json:"hostnames"`
	Vulnerabilities []string     `json:"vulns"`
	Tags            []string     `json:"tags"`
}

// NewShodanClient creates a new Shodan API client
func NewShodanClient(apiKey string) *ShodanClient {
	client := resty.New()
	client.SetBaseURL("https://api.shodan.io")
	client.SetQueryParam("key", apiKey)
	client.SetTimeout(30 * time.Second)

	return &ShodanClient{
		client:  client,
		apiKey:  apiKey,
		baseURL: "https://api.shodan.io",
	}
}

// Search performs a Shodan search query
func (s *ShodanClient) Search(ctx context.Context, query string) (*ShodanSearchResponse, error) {
	var result ShodanSearchResponse

	resp, err := s.client.R().
		SetContext(ctx).
		SetQueryParam("query", query).
		SetResult(&result).
		Get("/shodan/host/search")

	if err != nil {
		return nil, fmt.Errorf("shodan search request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("shodan search returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// GetHostInfo retrieves detailed information about a specific IP
func (s *ShodanClient) GetHostInfo(ctx context.Context, ip string) (*ShodanHostInfo, error) {
	var result ShodanHostInfo

	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/shodan/host/%s", ip))

	if err != nil {
		return nil, fmt.Errorf("shodan host info request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("shodan host info returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// SearchFacets performs a search with faceting for statistical analysis
func (s *ShodanClient) SearchFacets(ctx context.Context, query string, facets []string) (*ShodanSearchResponse, error) {
	var result ShodanSearchResponse

	req := s.client.R().
		SetContext(ctx).
		SetQueryParam("query", query)

	if len(facets) > 0 {
		for _, facet := range facets {
			req.SetQueryParam("facets", facet)
		}
	}

	resp, err := req.
		SetResult(&result).
		Get("/shodan/host/search")

	if err != nil {
		return nil, fmt.Errorf("shodan faceted search request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("shodan faceted search returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// ToJSON converts the search response to JSON string
func (s *ShodanSearchResponse) ToJSON() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal shodan response: %w", err)
	}
	return string(data), nil
}
