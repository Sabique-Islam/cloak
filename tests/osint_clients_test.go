package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Sabique-Islam/cloak/internal/ingestion"
)

// TestGitHubClient tests the GitHub search client
func TestGitHubClient(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("GITHUB_TOKEN not set, skipping GitHub client test")
	}

	client := ingestion.NewGitHubClient(token, "")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("SearchRepositories", func(t *testing.T) {
		result, err := client.SearchRepositories(ctx, "language:go stars:>1000")
		if err != nil {
			t.Fatalf("SearchRepositories failed: %v", err)
		}
		if result == nil {
			t.Fatal("SearchRepositories returned nil result")
		}
		t.Logf("Found %d repositories", result.TotalCount)
		if len(result.Items) > 0 {
			t.Logf("First result: %s (%s)", result.Items[0].FullName, result.Items[0].HTMLURL)
		}
	})

	t.Run("SearchCode", func(t *testing.T) {
		result, err := client.SearchCode(ctx, "package main language:go")
		if err != nil {
			t.Fatalf("SearchCode failed: %v", err)
		}
		if result == nil {
			t.Fatal("SearchCode returned nil result")
		}
		t.Logf("Found %d code results", result.TotalCount)
	})

	t.Run("ToJSON", func(t *testing.T) {
		result, err := client.SearchRepositories(ctx, "language:go")
		if err != nil {
			t.Fatalf("SearchRepositories failed: %v", err)
		}
		jsonStr, err := result.ToJSON()
		if err != nil {
			t.Fatalf("ToJSON failed: %v", err)
		}
		if jsonStr == "" {
			t.Fatal("ToJSON returned empty string")
		}
		t.Logf("JSON output length: %d bytes", len(jsonStr))
	})
}

// TestGreyNoiseClient tests the GreyNoise client
func TestGreyNoiseClient(t *testing.T) {
	apiKey := os.Getenv("GREYNOISE_API_KEY")
	if apiKey == "" {
		t.Skip("GREYNOISE_API_KEY not set, skipping GreyNoise client test")
	}

	client := ingestion.NewGreyNoiseClient(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Using a known public IP for testing
	testIP := "8.8.8.8"

	t.Run("LookupIP", func(t *testing.T) {
		result, err := client.LookupIP(ctx, testIP)
		if err != nil {
			t.Fatalf("LookupIP failed: %v", err)
		}
		if result == nil {
			t.Fatal("LookupIP returned nil result")
		}
		t.Logf("IP: %s, Seen: %v, Classification: %s", result.IP, result.Seen, result.Classification)
	})

	t.Run("QuickCheck", func(t *testing.T) {
		seen, err := client.QuickCheck(ctx, testIP)
		if err != nil {
			t.Fatalf("QuickCheck failed: %v", err)
		}
		t.Logf("IP %s seen by GreyNoise: %v", testIP, seen)
	})

	t.Run("ToJSON", func(t *testing.T) {
		result, err := client.LookupIP(ctx, testIP)
		if err != nil {
			t.Fatalf("LookupIP failed: %v", err)
		}
		jsonStr, err := result.ToJSON()
		if err != nil {
			t.Fatalf("ToJSON failed: %v", err)
		}
		if jsonStr == "" {
			t.Fatal("ToJSON returned empty string")
		}
		t.Logf("JSON output length: %d bytes", len(jsonStr))
	})
}

// TestOTXClient tests the AlienVault OTX client
func TestOTXClient(t *testing.T) {
	apiKey := os.Getenv("OTX_API_KEY")
	if apiKey == "" {
		t.Skip("OTX_API_KEY not set, skipping OTX client test")
	}

	client := ingestion.NewOTXClient(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testIP := "8.8.8.8"
	testDomain := "google.com"

	t.Run("GetIPReputation", func(t *testing.T) {
		result, err := client.GetIPReputation(ctx, testIP)
		if err != nil {
			t.Fatalf("GetIPReputation failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetIPReputation returned nil result")
		}
		t.Logf("IP: %s, Reputation: %d, Country: %s, Pulse Count: %d",
			testIP, result.Reputation, result.Country, result.PulseCount)
	})

	t.Run("GetDomainReputation", func(t *testing.T) {
		result, err := client.GetDomainReputation(ctx, testDomain)
		if err != nil {
			t.Fatalf("GetDomainReputation failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetDomainReputation returned nil result")
		}
		t.Logf("Domain: %s, Type: %s, Pulse Count: %d",
			result.Indicator, result.Type, result.PulseCount)
	})

	t.Run("GetPulses", func(t *testing.T) {
		result, err := client.GetPulses(ctx, "")
		if err != nil {
			t.Fatalf("GetPulses failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetPulses returned nil result")
		}
		t.Logf("Found %d pulses", result.Count)
		if len(result.Results) > 0 {
			t.Logf("First pulse: %s by %s", result.Results[0].Name, result.Results[0].AuthorName)
		}
	})

	t.Run("SearchIndicators", func(t *testing.T) {
		result, err := client.SearchIndicators(ctx, "malware")
		if err != nil {
			t.Fatalf("SearchIndicators failed: %v", err)
		}
		if result == nil {
			t.Fatal("SearchIndicators returned nil result")
		}
		t.Logf("Found %d search results", result.Count)
	})

	t.Run("ToJSON", func(t *testing.T) {
		result, err := client.GetPulses(ctx, "")
		if err != nil {
			t.Fatalf("GetPulses failed: %v", err)
		}
		jsonStr, err := result.ToJSON()
		if err != nil {
			t.Fatalf("ToJSON failed: %v", err)
		}
		if jsonStr == "" {
			t.Fatal("ToJSON returned empty string")
		}
		t.Logf("JSON output length: %d bytes", len(jsonStr))
	})
}

// TestShodanClient tests the Shodan client
func TestShodanClient(t *testing.T) {
	apiKey := os.Getenv("SHODAN_API_KEY")
	if apiKey == "" {
		t.Skip("SHODAN_API_KEY not set, skipping Shodan client test")
	}

	client := ingestion.NewShodanClient(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testIP := "8.8.8.8"

	t.Run("Search", func(t *testing.T) {
		result, err := client.Search(ctx, "apache")
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}
		if result == nil {
			t.Fatal("Search returned nil result")
		}
		t.Logf("Found %d total results", result.Total)
		if len(result.Matches) > 0 {
			t.Logf("First result: %s:%d (%s)",
				result.Matches[0].IP,
				result.Matches[0].Port,
				result.Matches[0].Organization)
		}
	})

	t.Run("GetHostInfo", func(t *testing.T) {
		result, err := client.GetHostInfo(ctx, testIP)
		if err != nil {
			t.Fatalf("GetHostInfo failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetHostInfo returned nil result")
		}
		t.Logf("Host: %s, OS: %s, Org: %s, Ports: %v",
			result.IP, result.OS, result.Organization, result.Ports)
	})

	t.Run("ToJSON", func(t *testing.T) {
		result, err := client.Search(ctx, "nginx")
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}
		jsonStr, err := result.ToJSON()
		if err != nil {
			t.Fatalf("ToJSON failed: %v", err)
		}
		if jsonStr == "" {
			t.Fatal("ToJSON returned empty string")
		}
		t.Logf("JSON output length: %d bytes", len(jsonStr))
	})
}

// TestVirusTotalClient tests the VirusTotal client
func TestVirusTotalClient(t *testing.T) {
	apiKey := os.Getenv("VIRUSTOTAL_API_KEY")
	if apiKey == "" {
		t.Skip("VIRUSTOTAL_API_KEY not set, skipping VirusTotal client test")
	}

	client := ingestion.NewVTClient(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testIP := "8.8.8.8"
	testDomain := "google.com"
	testHash := "44d88612fea8a8f36de82e1278abb02f" // EICAR test file MD5

	t.Run("GetIPReport", func(t *testing.T) {
		result, err := client.GetIPReport(ctx, testIP)
		if err != nil {
			t.Fatalf("GetIPReport failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetIPReport returned nil result")
		}
		t.Logf("IP: %s, Reputation: %d, Malicious: %d, Harmless: %d",
			result.Data.ID,
			result.Data.Attributes.Reputation,
			result.Data.Attributes.LastAnalysisStats.Malicious,
			result.Data.Attributes.LastAnalysisStats.Harmless)
	})

	t.Run("GetDomainReport", func(t *testing.T) {
		result, err := client.GetDomainReport(ctx, testDomain)
		if err != nil {
			t.Fatalf("GetDomainReport failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetDomainReport returned nil result")
		}
		t.Logf("Domain: %s, Reputation: %d, Detection Rate: %s",
			result.Data.ID,
			result.Data.Attributes.Reputation,
			result.GetDetectionRate())
	})

	t.Run("GetFileReport", func(t *testing.T) {
		result, err := client.GetFileReport(ctx, testHash)
		if err != nil {
			t.Fatalf("GetFileReport failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetFileReport returned nil result")
		}
		t.Logf("File: %s, Type: %s, Is Malicious: %v, Detection Rate: %s",
			result.Data.ID,
			result.Data.Attributes.TypeDescription,
			result.IsMalicious(),
			result.GetDetectionRate())
	})

	t.Run("ToJSON", func(t *testing.T) {
		result, err := client.GetDomainReport(ctx, testDomain)
		if err != nil {
			t.Fatalf("GetDomainReport failed: %v", err)
		}
		jsonStr, err := result.ToJSON()
		if err != nil {
			t.Fatalf("ToJSON failed: %v", err)
		}
		if jsonStr == "" {
			t.Fatal("ToJSON returned empty string")
		}
		t.Logf("JSON output length: %d bytes", len(jsonStr))
	})

	t.Run("HelperMethods", func(t *testing.T) {
		result, err := client.GetFileReport(ctx, testHash)
		if err != nil {
			t.Fatalf("GetFileReport failed: %v", err)
		}

		isMalicious := result.IsMalicious()
		detectionRate := result.GetDetectionRate()

		t.Logf("IsMalicious: %v", isMalicious)
		t.Logf("Detection Rate: %s", detectionRate)
	})
}

// TestAllClientsInitialization tests that all clients can be initialized without panicking
func TestAllClientsInitialization(t *testing.T) {
	t.Run("GitHubClient", func(t *testing.T) {
		client := ingestion.NewGitHubClient("test_token", "https://api.github.com")
		if client == nil {
			t.Fatal("NewGitHubClient returned nil")
		}
	})

	t.Run("GreyNoiseClient", func(t *testing.T) {
		client := ingestion.NewGreyNoiseClient("test_key")
		if client == nil {
			t.Fatal("NewGreyNoiseClient returned nil")
		}
	})

	t.Run("OTXClient", func(t *testing.T) {
		client := ingestion.NewOTXClient("test_key")
		if client == nil {
			t.Fatal("NewOTXClient returned nil")
		}
	})

	t.Run("ShodanClient", func(t *testing.T) {
		client := ingestion.NewShodanClient("test_key")
		if client == nil {
			t.Fatal("NewShodanClient returned nil")
		}
	})

	t.Run("VTClient", func(t *testing.T) {
		client := ingestion.NewVTClient("test_key")
		if client == nil {
			t.Fatal("NewVTClient returned nil")
		}
	})
}
