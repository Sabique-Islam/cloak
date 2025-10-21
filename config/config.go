package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment         string
	DBPath              string
	DBDialect           string
	LabDomain           string
	LabSubdomain        string
	GitHubToken         string
	GitHubOwner         string
	GitHubAPIBaseURL    string
	CloudflareAPIToken  string
	CloudflareZoneID    string
	VirusTotalAPIKey    string
	GreyNoiseAPIKey     string
	OTXAPIKey           string
	PassiveDNSAPIKey    string
	ShodanAPIKey        string
	WorkerConcurrency   int
	ProviderRPS         int
	PollIntervalSeconds int
	APIHost             string
	APIPort             string
	AlertsWebhookURL    string
	SentryDSN           string
	AutoTeardown        bool
	TeardownTTLSeconds  int
	DefaultDNSTTL       int
	DefaultCommitEmail  string
	DefaultCommitName   string
}

// Loads .env files from the repository root.
func Load() (*Config, error) {
	loadEnvFiles()

	return &Config{
		Environment:         getEnv("ENVIRONMENT", "development"),
		DBPath:              getEnv("DB_PATH", "./data/results.db"),
		DBDialect:           getEnv("DB_DIALECT", "sqlite"),
		LabDomain:           getEnv("LAB_DOMAIN", "smthg.xd"),
		LabSubdomain:        getEnv("LAB_SUBDOMAIN", "cloak"),
		GitHubToken:         os.Getenv("GITHUB_TOKEN"),
		GitHubOwner:         getEnv("GITHUB_OWNER", ""),
		GitHubAPIBaseURL:    getEnv("GITHUB_API_BASE_URL", "https://api.github.com"),
		CloudflareAPIToken:  os.Getenv("CLOUDFLARE_API_TOKEN"),
		CloudflareZoneID:    os.Getenv("CLOUDFLARE_ZONE_ID"),
		VirusTotalAPIKey:    os.Getenv("VIRUSTOTAL_API_KEY"),
		GreyNoiseAPIKey:     os.Getenv("GREYNOISE_API_KEY"),
		OTXAPIKey:           os.Getenv("OTX_API_KEY"),
		PassiveDNSAPIKey:    os.Getenv("PASSIVEDNS_API_KEY"),
		ShodanAPIKey:        os.Getenv("SHODAN_API_KEY"),
		WorkerConcurrency:   getEnvAsInt("WORKER_CONCURRENCY", 6),
		ProviderRPS:         getEnvAsInt("PROVIDER_RPS", 5),
		PollIntervalSeconds: getEnvAsInt("POLL_INTERVAL_SECONDS", 300),
		APIHost:             getEnv("API_HOST", "127.0.0.1"),
		APIPort:             getEnv("API_PORT", "8080"),
		AlertsWebhookURL:    os.Getenv("ALERTS_WEBHOOK_URL"),
		SentryDSN:           os.Getenv("SENTRY_DSN"),
		AutoTeardown:        getEnvAsBool("AUTO_TEARDOWN", false),
		TeardownTTLSeconds:  getEnvAsInt("TEARDOWN_TTL_SECONDS", 86400),
		DefaultDNSTTL:       getEnvAsInt("DEFAULT_DNS_TTL", 120),
		DefaultCommitEmail:  getEnv("DEFAULT_COMMIT_EMAIL", ""),
		DefaultCommitName:   getEnv("DEFAULT_COMMIT_NAME", "Cloak_Lab_Bot"),
	}, nil
}

// Locates and loads .env file from the repository root directory.
func loadEnvFiles() {
	root := findRepoRoot()
	_ = godotenv.Load(filepath.Join(root, ".env"))
}

// findRepoRoot walks up to the directory tree from the current working directory
// to find the repository root by looking for go.mod, returns "." if root cant be found.
func findRepoRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

// Retrieves value of env variable named by key.
// If the variable is not present or is empty, it returns the defaultValue.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Retrieves value of env variable named by key and converts it to an integer. If the variable is not present,empty, or cant be parsed as an integer, it returns the defaultValue.
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// Retrieves value of env variable named by key and converts it to a boolean. It accepts values like "1", "t", "true", "0", "f", "false" (case-insensitive). If the variable is not present, empty, or cant be parsed as a boolean, it returns the defaultValue.
func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
