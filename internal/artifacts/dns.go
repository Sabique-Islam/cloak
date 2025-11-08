package artifacts

/*
Implement DNS artifact generation
- Create A, CNAME, TXT records via Cloudflare API
- Validate records using miekg/dns
- Store records in database
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/Sabique-Islam/cloak/config"
	"github.com/Sabique-Islam/cloak/internal/storage"
	"github.com/cloudflare/cloudflare-go"
	"github.com/miekg/dns"
	"gorm.io/gorm"
)

// Type of DNS record
type DNSRecordType string

const (
	RecordTypeA     DNSRecordType = "A"
	RecordTypeCNAME DNSRecordType = "CNAME"
	RecordTypeTXT   DNSRecordType = "TXT"
)

// DNS-specific metadata as JSON
type DNSArtifactMetadata struct {
	RecordType  string `json:"record_type"`
	RecordName  string `json:"record_name"`
	RecordValue string `json:"record_value"`
	TTL         int    `json:"ttl"`
	Proxied     bool   `json:"proxied"`
	ZoneID      string `json:"zone_id"`
	ZoneName    string `json:"zone_name"`
}

// Info needed to delete the DNS record
type DNSTeardownInfo struct {
	CloudflareRecordID string    `json:"cloudflare_record_id"`
	ZoneID             string    `json:"zone_id"`
	RecordName         string    `json:"record_name"`
	RecordType         string    `json:"record_type"`
	CreatedAt          time.Time `json:"created_at"`
}

// DNS artifact creation via Cloudflare
type DNSGenerator struct {
	cfg       *config.Config
	db        *gorm.DB
	cfClient  *cloudflare.API
	dnsClient *dns.Client
}

// Create new DNS artifact generator
func NewDNSGenerator(cfg *config.Config, db *gorm.DB) (*DNSGenerator, error) {
	// Initialize Cloudflare API client
	cfAPI, err := cloudflare.NewWithAPIToken(cfg.CloudflareAPIToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudflare client: %w", err)
	}

	// Init DNS client for validation
	dnsClient := new(dns.Client)
	dnsClient.Timeout = 10 * time.Second

	return &DNSGenerator{
		cfg:       cfg,
		db:        db,
		cfClient:  cfAPI,
		dnsClient: dnsClient,
	}, nil
}

// exponentialBackoffRetry retries an operation with exponential backoff
// maxAttempts: maximum number of retry attempts (default 5)
// initialDelay: initial delay between retries (default 1 second)
// maxDelay: maximum delay between retries (default 30 seconds)
func exponentialBackoffRetry(ctx context.Context, maxAttempts int, initialDelay time.Duration, maxDelay time.Duration, operation func() error) error {
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	if initialDelay <= 0 {
		initialDelay = 1 * time.Second
	}
	if maxDelay <= 0 {
		maxDelay = 30 * time.Second
	}

	var lastErr error
	delay := initialDelay

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := operation(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < maxAttempts-1 {
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			waitTime := time.Duration(math.Min(float64(delay+jitter), float64(maxDelay)))

			select {
			case <-time.After(waitTime):
				// Continue to next attempt
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled: %w", ctx.Err())
			}

			delay = time.Duration(math.Min(float64(delay*2), float64(maxDelay)))
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", maxAttempts, lastErr)
}

// Create A record artifact
func (g *DNSGenerator) GenerateARecord(ctx context.Context, name string, ipAddress string) (*storage.Artifact, error) {
	return g.createDNSRecord(ctx, RecordTypeA, name, ipAddress, false)
}

// Create CNAME record artifact
func (g *DNSGenerator) GenerateCNAMERecord(ctx context.Context, name string, target string) (*storage.Artifact, error) {
	return g.createDNSRecord(ctx, RecordTypeCNAME, name, target, false)
}

// Create TXT record artifact
func (g *DNSGenerator) GenerateTXTRecord(ctx context.Context, name string, content string) (*storage.Artifact, error) {
	return g.createDNSRecord(ctx, RecordTypeTXT, name, content, false)
}

// GenerateRandomArtifact creates a random DNS artifact with unique identifier
func (g *DNSGenerator) GenerateRandomArtifact(ctx context.Context, recordType DNSRecordType) (*storage.Artifact, error) {
	timestamp := time.Now().Unix()
	uniqueID := fmt.Sprintf("cloak-%d", timestamp)

	switch recordType {
	case RecordTypeA:
		// Use a reserved test IP (RFC 5737)
		return g.GenerateARecord(ctx, uniqueID, "203.0.113.1")
	case RecordTypeCNAME:
		return g.GenerateCNAMERecord(ctx, uniqueID, fmt.Sprintf("target-%d.%s", timestamp, g.cfg.LabDomain))
	case RecordTypeTXT:
		return g.GenerateTXTRecord(ctx, uniqueID, fmt.Sprintf("cloak-experiment-%d", timestamp))
	default:
		return nil, fmt.Errorf("unsupported record type: %s", recordType)
	}
}

// createDNSRecord is the core function that creates DNS records via Cloudflare
func (g *DNSGenerator) createDNSRecord(ctx context.Context, recordType DNSRecordType, name string, value string, proxied bool) (*storage.Artifact, error) {
	// Construct full record name
	fullRecordName := fmt.Sprintf("%s.%s.%s", name, g.cfg.LabSubdomain, g.cfg.LabDomain)

	// Prepare Cloudflare DNS record
	cfRecord := cloudflare.DNSRecord{
		Type:    string(recordType),
		Name:    fullRecordName,
		Content: value,
		TTL:     g.cfg.DefaultDNSTTL,
		Proxied: &proxied,
	}

	// Create record in Cloudflare
	zoneID := cloudflare.ZoneIdentifier(g.cfg.CloudflareZoneID)
	createParams := cloudflare.CreateDNSRecordParams{
		Type:    cfRecord.Type,
		Name:    cfRecord.Name,
		Content: cfRecord.Content,
		TTL:     cfRecord.TTL,
		Proxied: cfRecord.Proxied,
	}

	response, err := g.cfClient.CreateDNSRecord(ctx, zoneID, createParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS record in Cloudflare: %w", err)
	}

	// Validate DNS record with exponential backoff retry (up to 30 seconds)
	validationErr := exponentialBackoffRetry(ctx, 6, 1*time.Second, 30*time.Second, func() error {
		valid, err := g.validateDNSRecord(fullRecordName, string(recordType), value)
		if err != nil {
			return err
		}
		if !valid {
			return fmt.Errorf("DNS record not yet propagated")
		}
		return nil
	})

	// Log validation error but don't fail artifact creation - DNS propagation timing varies
	if validationErr != nil {
		fmt.Printf("Warning: DNS validation timed out for %s after retries: %v\n", fullRecordName, validationErr)
	}

	// Prepare metadata
	metadata := DNSArtifactMetadata{
		RecordType:  string(recordType),
		RecordName:  fullRecordName,
		RecordValue: value,
		TTL:         g.cfg.DefaultDNSTTL,
		Proxied:     proxied,
		ZoneID:      g.cfg.CloudflareZoneID,
		ZoneName:    g.cfg.LabDomain,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Prepare teardown info
	teardownInfo := DNSTeardownInfo{
		CloudflareRecordID: response.ID,
		ZoneID:             g.cfg.CloudflareZoneID,
		RecordName:         fullRecordName,
		RecordType:         string(recordType),
		CreatedAt:          time.Now(),
	}
	teardownJSON, err := json.Marshal(teardownInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal teardown info: %w", err)
	}

	// Create artifact record in database
	artifact := &storage.Artifact{
		Type:         "dns",
		Identifier:   fullRecordName,
		Metadata:     string(metadataJSON),
		TeardownInfo: string(teardownJSON),
		Status:       "active",
	}

	if err := g.db.Create(artifact).Error; err != nil {
		// Attempt to clean up Cloudflare record if DB insert fails
		_ = g.cfClient.DeleteDNSRecord(ctx, zoneID, response.ID)
		return nil, fmt.Errorf("failed to save artifact to database: %w", err)
	}

	// Create audit log entry - log errors but don't fail operation
	auditLog := &storage.AuditLog{
		Actor:    "dns_generator",
		Action:   "create_dns_artifact",
		Resource: fmt.Sprintf("artifact:%d", artifact.ID),
		Details:  string(metadataJSON),
	}
	if auditErr := g.db.Create(auditLog).Error; auditErr != nil {
		// Log the error without failing - audit trail is important but shouldn't block operations
		fmt.Printf("Error: Failed to create audit log for DNS artifact creation: %v\n", auditErr)
	}

	return artifact, nil
}

// validateDNSRecord performs actual DNS lookup to verify record exists
func (g *DNSGenerator) validateDNSRecord(name string, recordType string, expectedValue string) (bool, error) {
	var qType uint16
	switch recordType {
	case "A":
		qType = dns.TypeA
	case "CNAME":
		qType = dns.TypeCNAME
	case "TXT":
		qType = dns.TypeTXT
	default:
		return false, fmt.Errorf("unsupported record type: %s", recordType)
	}

	// Create DNS query
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), qType)
	msg.RecursionDesired = true

	// Query public DNS server (Cloudflare's 1.1.1.1)
	response, _, err := g.dnsClient.Exchange(msg, "1.1.1.1:53")
	if err != nil {
		return false, fmt.Errorf("DNS query failed: %w", err)
	}

	if response == nil || len(response.Answer) == 0 {
		return false, nil
	}

	// Check if any answer matches our expected value
	for _, ans := range response.Answer {
		switch recordType {
		case "A":
			if a, ok := ans.(*dns.A); ok {
				if a.A.String() == expectedValue {
					return true, nil
				}
			}
		case "CNAME":
			if cname, ok := ans.(*dns.CNAME); ok {
				if cname.Target == dns.Fqdn(expectedValue) {
					return true, nil
				}
			}
		case "TXT":
			if txt, ok := ans.(*dns.TXT); ok {
				for _, t := range txt.Txt {
					if t == expectedValue {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}

// DeleteDNSRecord removes a DNS record from Cloudflare and marks artifact as deleted
func (g *DNSGenerator) DeleteDNSRecord(ctx context.Context, artifact *storage.Artifact) error {
	// Parse teardown info
	var teardownInfo DNSTeardownInfo
	if err := json.Unmarshal([]byte(artifact.TeardownInfo), &teardownInfo); err != nil {
		return fmt.Errorf("failed to parse teardown info: %w", err)
	}

	// Delete from Cloudflare
	zoneID := cloudflare.ZoneIdentifier(teardownInfo.ZoneID)

	if err := g.cfClient.DeleteDNSRecord(ctx, zoneID, teardownInfo.CloudflareRecordID); err != nil {
		return fmt.Errorf("failed to delete DNS record from Cloudflare: %w", err)
	}

	// Update artifact status in database
	artifact.Status = "deleted"
	if err := g.db.Save(artifact).Error; err != nil {
		return fmt.Errorf("failed to update artifact status: %w", err)
	}

	// Create audit log - log errors but don't fail operation
	auditLog := &storage.AuditLog{
		Actor:    "dns_generator",
		Action:   "delete_dns_artifact",
		Resource: fmt.Sprintf("artifact:%d", artifact.ID),
		Details:  fmt.Sprintf(`{"record_name":"%s","record_type":"%s"}`, teardownInfo.RecordName, teardownInfo.RecordType),
	}
	if auditErr := g.db.Create(auditLog).Error; auditErr != nil {
		fmt.Printf("Error: Failed to create audit log for DNS artifact deletion: %v\n", auditErr)
	}

	return nil
}

// Return all active DNS artifacts
func (g *DNSGenerator) ListActiveDNSArtifacts() ([]storage.Artifact, error) {
	var artifacts []storage.Artifact
	err := g.db.Where("type = ? AND status = ?", "dns", "active").Find(&artifacts).Error
	return artifacts, err
}

// Return parsed metadata for a DNS artifact
func (g *DNSGenerator) GetDNSRecordInfo(artifact *storage.Artifact) (*DNSArtifactMetadata, error) {
	var metadata DNSArtifactMetadata
	if err := json.Unmarshal([]byte(artifact.Metadata), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	return &metadata, nil
}
