package tests

import (
	"strings"
	"testing"

	"github.com/Sabique-Islam/cloak/config"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	separator := strings.Repeat("=", 60)
	t.Log("\n" + separator)
	t.Log("CONFIGURATION LOADED FROM ENVIRONMENT")
	t.Log(separator)

	t.Logf("%-30s: %s", "Environment", cfg.Environment)
	t.Logf("%-30s: %s", "DB Path", cfg.DBPath)
	t.Logf("%-30s: %s", "DB Dialect", cfg.DBDialect)

	t.Log("\n--- Lab / Domain ---")
	t.Logf("%-30s: %s", "Lab Domain", cfg.LabDomain)
	t.Logf("%-30s: %s", "Lab Subdomain", cfg.LabSubdomain)
	t.Logf("%-30s: %s.%s", "Full Lab Domain", cfg.LabSubdomain, cfg.LabDomain)

	t.Log("\n--- GitHub ---")
	t.Logf("%-30s: %s", "GitHub Token", maskSecret(cfg.GitHubToken))
	t.Logf("%-30s: %s", "GitHub Owner", cfg.GitHubOwner)
	t.Logf("%-30s: %s", "GitHub API Base URL", cfg.GitHubAPIBaseURL)

	t.Log("\n--- Cloudflare ---")
	t.Logf("%-30s: %s", "Cloudflare API Token", maskSecret(cfg.CloudflareAPIToken))
	t.Logf("%-30s: %s", "Cloudflare Zone ID", maskSecret(cfg.CloudflareZoneID))

	t.Log("\n--- OSINT Provider API Keys ---")
	t.Logf("%-30s: %s", "VirusTotal API Key", maskSecret(cfg.VirusTotalAPIKey))
	t.Logf("%-30s: %s", "GreyNoise API Key", maskSecret(cfg.GreyNoiseAPIKey))
	t.Logf("%-30s: %s", "OTX API Key", maskSecret(cfg.OTXAPIKey))
	t.Logf("%-30s: %s", "PassiveDNS API Key", maskSecret(cfg.PassiveDNSAPIKey))
	t.Logf("%-30s: %s", "Shodan API Key", maskSecret(cfg.ShodanAPIKey))

	t.Log("\n--- Worker / Polling ---")
	t.Logf("%-30s: %d", "Worker Concurrency", cfg.WorkerConcurrency)
	t.Logf("%-30s: %d", "Provider RPS", cfg.ProviderRPS)
	t.Logf("%-30s: %d", "Poll Interval Seconds", cfg.PollIntervalSeconds)

	t.Log("\n--- API Server ---")
	t.Logf("%-30s: %s", "API Host", cfg.APIHost)
	t.Logf("%-30s: %s", "API Port", cfg.APIPort)
	t.Logf("%-30s: http://%s:%s", "Full API URL", cfg.APIHost, cfg.APIPort)

	t.Log("\n--- Security / Ops ---")
	t.Logf("%-30s: %s", "Alerts Webhook URL", maskSecret(cfg.AlertsWebhookURL))
	t.Logf("%-30s: %s", "Sentry DSN", maskSecret(cfg.SentryDSN))

	t.Log("\n--- Teardown & Safety ---")
	t.Logf("%-30s: %v", "Auto Teardown", cfg.AutoTeardown)
	t.Logf("%-30s: %d seconds (%d hours)", "Teardown TTL", cfg.TeardownTTLSeconds, cfg.TeardownTTLSeconds/3600)

	t.Log("\n--- Experiment Defaults ---")
	t.Logf("%-30s: %d seconds", "Default DNS TTL", cfg.DefaultDNSTTL)
	t.Logf("%-30s: %s", "Default Commit Email", cfg.DefaultCommitEmail)
	t.Logf("%-30s: %s", "Default Commit Name", cfg.DefaultCommitName)

	t.Log("\n" + separator)
}

func maskSecret(secret string) string {
	if secret == "" {
		return "<not set>"
	}
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "..." + secret[len(secret)-4:]
}
