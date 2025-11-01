
---

# Cloak Setup Guide

---

## Prerequisites


| Tool | Minimum Version | Installation |
|------|----------------|--------------|
| **Go** | 1.25.3 | [Download](https://go.dev/dl/) |
| **Git** | 2.0+ | [Download](https://git-scm.com/downloads) |
| **golangci-lint** | 2.5.0+ | See [Installation](#golangci-lint-installation) |
| **Make** | Build automation | Pre-installed on macOS/Linux |
| **SQLite Browser** | Database inspection | [Download](https://sqlitebrowser.org/) |



### API Access Requirements

<details>
<summary>Click to expand</summary>

| Provider | Required For | Registration |
|----------|-------------|--------------|
| **GitHub** | Artifact generation & search | [Create Token](https://github.com/settings/tokens) |
| **Cloudflare** | DNS artifact generation | [Get API Token](https://dash.cloudflare.com/profile/api-tokens) |
| **VirusTotal** | OSINT ingestion | [API Key](https://www.virustotal.com/gui/join-us) |
| **GreyNoise** | OSINT ingestion | [API Key](https://www.greynoise.io/) |
| **AlienVault OTX** | OSINT ingestion | [API Key](https://otx.alienvault.com/) |
| **Shodan** | OSINT ingestion | [API Key](https://account.shodan.io/) |

</details>

---

## Initial Setup

### 1. Fork the Repo & Clone ur Fork

### 2. Verify Go Installation

```bash
go version
# Expected output: go version go1.25.3 or higher
```

### 3. Install Dependencies

```bash
go mod download
go mod verify
```

Expected output:
```
all modules verified
```

### 4. golangci-lint Installation

<details>
<summary>macOS</summary>

```bash
brew install golangci-lint
```

</details>

<details>
<summary>Linux</summary>

todo

</details>

<details>
<summary>Windows</summary>

todo

</details>

Verify installation:
```bash
golangci-lint --version
# Expected: golangci-lint has version 2.5.0 or higher
```

---

## Configuration

### Environment Variables Setup

1. **Copy example configuration:**

```bash
cp .env.example .env
```

2. **Edit `.env` with your credentials:**

### Validate Configuration

```bash
go test -v ./tests -run TestLoadConfig
```

Expected output should display all loaded configuration values with masked secrets.

---

## Database Setup

### Automatic Migration

The database is automatically initialized on first run. The SQLite database will be created at the path specified in `DB_PATH`.


### Database Schema

<details>
<summary>View Database Tables</summary>

**artifacts** - Stores generated GitHub commits, repos, and DNS records
- `id` - Primary key
- `artifact_type` - Type: github_commit, github_repo, dns_record
- `identifier` - Unique identifier (commit SHA, repo name, DNS record ID)
- `metadata` - JSON metadata
- `created_at` - Timestamp
- `status` - active, deleted, expired

**detections** - Logs when artifacts are found in OSINT feeds
- `id` - Primary key
- `artifact_id` - Foreign key to artifacts
- `provider` - OSINT provider name
- `detected_at` - Detection timestamp
- `raw_response` - JSON response from provider
- `confidence` - Detection confidence score

**jobs** - Manages async tasks
- `id` - Primary key
- `job_type` - Type: generation, ingestion, teardown
- `status` - pending, running, completed, failed
- `created_at`, `started_at`, `completed_at` - Timestamps
- `error_message` - Error details if failed

**audit_logs** - Tracks all actions
- `id` - Primary key
- `action` - Action performed
- `user` - User/system identifier
- `timestamp` - When action occurred
- `details` - JSON details
- `reason` - Why action was performed

</details>

---

## Running the Application

### Build the Application

```bash
# Build main binary
go build -o bin/cloakd ./cmd/cloakd

# Build utility scripts
go build -o bin/seed ./scripts/seed
go build -o bin/monitor ./scripts/monitor
go build -o bin/export ./scripts/export
go build -o bin/teardown ./scripts/teardown
```

### Access the Application

Once running, access the web UI:

```
http://127.0.0.1:8080
```

API endpoints:
- `GET /api/artifacts` - List artifacts
- `GET /api/detections` - List detections
- `GET /api/metrics` - System metrics
- `GET /api/audit` - Audit logs

---

## Workflow

### Project Structure

```
cloak/
├── cmd/
│   └── cloakd/          # Main application entry point
├── config/              # Configuration management
├── internal/
│   ├── api/            # HTTP handlers and routes
│   ├── artifacts/      # Artifact generation (GitHub, DNS)
│   ├── audit/          # Audit logging
│   ├── ingestion/      # OSINT provider clients
│   ├── storage/        # Database layer
│   └── utils/          # Shared utilities
├── scripts/            # Utility scripts
│   ├── seed/          # Generate sample data
│   ├── monitor/       # Monitor artifact propagation
│   ├── export/        # Export results
│   └── teardown/      # Clean up artifacts
├── tests/             # Test files
└── data/              # Runtime data
    ├── exports/       # Exported data
    └── logs/          # Application logs
```

### Code Organization Guidelines

- Place business logic in `internal/` packages
- Add tests for all new features

---

## Code Quality Checks

### Linting

#### Run All Linters

```bash
golangci-lint run
```

#### Run with Auto-fix

```bash
golangci-lint run --fix
```

#### Run Specific Linters

```bash
# Check for common mistakes
golangci-lint run --disable-all -E errcheck -E gosimple -E staticcheck

# Check code formatting
golangci-lint run --disable-all -E gofmt -E goimports

# Security checks
golangci-lint run --disable-all -E gosec
```

#### Verbose Output

```bash
golangci-lint run -v
```

### Common Linting Issues and Fixes

<details>
<summary>Unused Variables</summary>

**Error:**
```
variable 'foo' is unused
```

**Fix:**
```go
// Remove unused variable or use blank identifier
_ = foo
```

</details>

<details>
<summary>Error Handling</summary>

**Error:**
```
Error return value is not checked
```

**Fix:**
```go
// Before
file.Close()

// After
if err := file.Close(); err != nil {
    log.Printf("error closing file: %v", err)
}
```

</details>

<details>
<summary>Formatting Issues</summary>

**Fix all formatting:**
```bash
gofmt -w .
goimports -w .
```

</details>

### Pre-commit Checks

Run before committing:

```bash
# Format code
gofmt -w .

# Run linters
golangci-lint run

# Run tests
go test ./...

# Check for security issues
golangci-lint run --disable-all -E gosec
```

---

## Testing

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run specific test file
go test -v ./tests -run TestConfig

# Run DNS tests only
go test -v ./tests -run TestDNS

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---
