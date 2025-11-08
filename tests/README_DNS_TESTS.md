`dns_artifacts_test.go` - Tests for DNS record creation via Cloudflare API

```bash
go test -v ./tests
```

### Only DNS Tests
```bash
go test -v ./tests -run TestDNS
```

### Run Specific Test
```bash
go test -v ./tests -run TestDNSGeneratorInitialization
```

### Unit Tests (No API Keys Required)

- `TestDNSGeneratorInitialization` - Tests generator initialization
- `TestDNSGeneratorInitializationWithoutAPIToken` - Tests error handling
- `TestGetDNSRecordInfo` - Tests metadata parsing
- `TestDNSRecordMetadataJSON` - Tests JSON serialization
- `TestDNSTeardownInfoJSON` - Tests teardown info serialization

**Run unit tests:**
```bash
go test -v ./tests -run "TestDNSGenerator|TestGetDNS|TestDNSRecord|TestDNSTeardown"
```

### Integration Tests (Require API Keys)

- `TestGenerateARecord` - Creates and validates A records
- `TestGenerateCNAMERecord` - Creates and validates CNAME records
- `TestGenerateTXTRecord` - Creates and validates TXT records
- `TestGenerateRandomArtifact` - Tests random artifact generation
- `TestListActiveDNSArtifacts` - Tests querying active artifacts
- `TestDeleteDNSRecord` - Tests artifact deletion

**Run integration tests (requires Cloudflare credentials):**
```bash
# Set .env first, then

go test -v ./tests -run "TestGenerate|TestList|TestDelete"
```

## Required

- `CLOUDFLARE_API_TOKEN`
- `CLOUDFLARE_ZONE_ID`
- `LAB_DOMAIN`
- `LAB_SUBDOMAIN`
- `DEFAULT_DNS_TTL` - TTL for DNS records (default: 120)

## Test Database

Tests use an in-memory SQLite database that is created and destroyed for each test run. No persistent data is stored.

## Safety

- All integration tests automatically clean up created DNS records
- Tests use RFC 5737 reserved IP addresses (203.0.113.0/24) for A records
- Unique timestamped names prevent conflicts
- Audit logs are created for all operations 

## Coverage

Run tests with coverage:
```bash
go test -v ./tests -cover -run TestDNS
```

Generate coverage report:
```bash
go test -v ./tests -coverprofile=coverage.out -run TestDNS
go tool cover -html=coverage.out
```
