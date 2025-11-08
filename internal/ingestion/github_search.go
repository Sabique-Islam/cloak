package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// GitHubClient handles GitHub Search API queries
type GitHubClient struct {
	client  *resty.Client
	token   string
	baseURL string
}

// GitHubSearchResult represents a GitHub search response
type GitHubSearchResult struct {
	TotalCount        int                `json:"total_count"`
	IncompleteResults bool               `json:"incomplete_results"`
	Items             []GitHubSearchItem `json:"items"`
}

// GitHubSearchItem represents a single search result item
type GitHubSearchItem struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	FullName    string            `json:"full_name"`
	HTMLURL     string            `json:"html_url"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Owner       GitHubOwner       `json:"owner"`
	Repository  *GitHubRepository `json:"repository,omitempty"`
}

// GitHubOwner represents repository owner information
type GitHubOwner struct {
	Login   string `json:"login"`
	ID      int64  `json:"id"`
	HTMLURL string `json:"html_url"`
	Type    string `json:"type"`
}

// GitHubRepository represents repository information in commits
type GitHubRepository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(token, baseURL string) *GitHubClient {
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}

	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetHeader("Accept", "application/vnd.github.v3+json")
	if token != "" {
		client.SetHeader("Authorization", fmt.Sprintf("token %s", token))
	}
	client.SetTimeout(30 * time.Second)

	return &GitHubClient{
		client:  client,
		token:   token,
		baseURL: baseURL,
	}
}

// SearchRepositories searches for repositories containing the query string
func (g *GitHubClient) SearchRepositories(ctx context.Context, query string) (*GitHubSearchResult, error) {
	var result GitHubSearchResult

	resp, err := g.client.R().
		SetContext(ctx).
		SetQueryParam("q", query).
		SetQueryParam("per_page", "100").
		SetResult(&result).
		Get("/search/repositories")

	if err != nil {
		return nil, fmt.Errorf("github search repositories request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("github search repositories returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// SearchCommits searches for commits containing the query string
func (g *GitHubClient) SearchCommits(ctx context.Context, query string) (*GitHubSearchResult, error) {
	var result GitHubSearchResult

	resp, err := g.client.R().
		SetContext(ctx).
		SetQueryParam("q", query).
		SetQueryParam("per_page", "100").
		SetHeader("Accept", "application/vnd.github.cloak-preview+json"). // Required for commit search
		SetResult(&result).
		Get("/search/commits")

	if err != nil {
		return nil, fmt.Errorf("github search commits request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("github search commits returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// SearchCode searches for code containing the query string
func (g *GitHubClient) SearchCode(ctx context.Context, query string) (*GitHubSearchResult, error) {
	var result GitHubSearchResult

	resp, err := g.client.R().
		SetContext(ctx).
		SetQueryParam("q", query).
		SetQueryParam("per_page", "100").
		SetResult(&result).
		Get("/search/code")

	if err != nil {
		return nil, fmt.Errorf("github search code request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("github search code returned error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return &result, nil
}

// ToJSON converts the search result to JSON string
func (g *GitHubSearchResult) ToJSON() (string, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return "", fmt.Errorf("failed to marshal github search result: %w", err)
	}
	return string(data), nil
}
