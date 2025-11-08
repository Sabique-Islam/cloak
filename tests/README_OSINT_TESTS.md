# OSINT Clients Tests

This directory contains integration tests for all OSINT API clients.

## Overview

The test suite validates the following OSINT clients:
- **GitHub Search** - Repository, code, and commit searches
- **GreyNoise** - IP reputation and scanning data
- **AlienVault OTX** - Threat intelligence pulses and indicators
- **Shodan** - Internet-wide device and service scanning
- **VirusTotal** - File, URL, domain, and IP analysis

## Running Tests

### 1. Basic Initialization Tests (No API Keys Required)

Test that all clients can be initialized:

```bash
go test -v ./tests -run TestAllClientsInitialization
```

### 2. Full Integration Tests (API Keys Required)

Set your API keys as environment variables in your `.env` file or export them:

```bash
# Copy .env.example to .env and fill in your API keys
cp .env.example .env

# Then run specific tests
go test -v ./tests -run TestGitHubClient
go test -v ./tests -run TestGreyNoiseClient
go test -v ./tests -run TestOTXClient
go test -v ./tests -run TestShodanClient
go test -v ./tests -run TestVirusTotalClient
```

### 3. Run All Tests

```bash
# Run all tests (will skip tests for missing API keys)
go test -v ./tests

# Run with coverage
go test -v -cover ./tests
```

## Required Environment Variables

| Variable | Description | Get API Key From |
|----------|-------------|------------------|
| `GITHUB_TOKEN` | GitHub Personal Access Token | https://github.com/settings/tokens |
| `GREYNOISE_API_KEY` | GreyNoise API Key | https://www.greynoise.io/account/api-key |
| `OTX_API_KEY` | AlienVault OTX API Key | https://otx.alienvault.com/api |
| `SHODAN_API_KEY` | Shodan API Key | https://account.shodan.io/ |
| `VIRUSTOTAL_API_KEY` | VirusTotal API Key | https://www.virustotal.com/gui/my-apikey |

## Test Behavior

- **Automatic Skipping**: Tests automatically skip if their required API key is not set
- **Safe Test Data**: Tests use well-known, safe IPs and domains (e.g., 8.8.8.8, google.com)
- **Timeouts**: All tests have 30-second context timeouts to prevent hanging
- **Rate Limiting**: Be mindful of API rate limits when running tests repeatedly

## Example: Running a Single Test

```bash
# Test GitHub client only
export GITHUB_TOKEN="your_github_token_here"
go test -v ./tests -run TestGitHubClient

# Test VirusTotal with verbose output
export VIRUSTOTAL_API_KEY="your_vt_key_here"
go test -v ./tests -run TestVirusTotalClient
```

## CI/CD Integration

For CI/CD pipelines, set API keys as secrets and run:

```bash
go test ./tests -short  # Run only fast tests
```

## Notes

- Tests use real API endpoints, so they require internet connectivity
- Some APIs have rate limits - avoid running tests too frequently
- Free tier API keys may have limited functionality
- Test results will vary based on the current state of external services
