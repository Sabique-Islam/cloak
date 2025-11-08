// AI generated

package tests

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Sabique-Islam/cloak/config"
	"github.com/Sabique-Islam/cloak/internal/artifacts"
	"github.com/Sabique-Islam/cloak/internal/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Auto-migrate test schema
	err = db.AutoMigrate(&storage.Artifact{}, &storage.AuditLog{})
	if err != nil {
		t.Fatalf("Failed to migrate test schema: %v", err)
	}

	return db
}

// TestDNSGeneratorInitialization tests creating a new DNS generator
func TestDNSGeneratorInitialization(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db := setupTestDB(t)

	generator, err := artifacts.NewDNSGenerator(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create DNS generator: %v", err)
	}

	if generator == nil {
		t.Fatal("DNS generator should not be nil")
	}

	t.Log("✓ DNS generator initialized successfully")
}

// TestDNSGeneratorInitializationWithoutAPIToken tests that initialization fails without API token
func TestDNSGeneratorInitializationWithoutAPIToken(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Save original token and clear it
	originalToken := cfg.CloudflareAPIToken
	cfg.CloudflareAPIToken = ""

	db := setupTestDB(t)

	generator, err := artifacts.NewDNSGenerator(cfg, db)

	// Restore original token
	cfg.CloudflareAPIToken = originalToken

	if err == nil {
		t.Fatal("Expected error when initializing without API token, got nil")
	}

	if generator != nil {
		t.Fatal("Generator should be nil when initialization fails")
	}

	t.Logf("✓ Correctly failed with error: %v", err)
}

// TestGenerateARecord tests A record creation (requires valid Cloudflare credentials)
func TestGenerateARecord(t *testing.T) {
	if os.Getenv("CLOUDFLARE_API_TOKEN") == "" {
		t.Skip("Skipping live DNS test: CLOUDFLARE_API_TOKEN not set")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db := setupTestDB(t)

	generator, err := artifacts.NewDNSGenerator(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create DNS generator: %v", err)
	}

	ctx := context.Background()
	testName := "test-a-record"
	testIP := "203.0.113.1" // RFC 5737 test IP

	t.Logf("Creating A record: %s.%s.%s -> %s", testName, cfg.LabSubdomain, cfg.LabDomain, testIP)

	artifact, err := generator.GenerateARecord(ctx, testName, testIP)
	if err != nil {
		t.Fatalf("Failed to generate A record: %v", err)
	}

	// Verify artifact was created
	if artifact.ID == 0 {
		t.Fatal("Artifact ID should not be zero")
	}

	if artifact.Type != "dns" {
		t.Errorf("Expected artifact type 'dns', got '%s'", artifact.Type)
	}

	if artifact.Status != "active" {
		t.Errorf("Expected artifact status 'active', got '%s'", artifact.Status)
	}

	// Parse and verify metadata
	var metadata artifacts.DNSArtifactMetadata
	err = json.Unmarshal([]byte(artifact.Metadata), &metadata)
	if err != nil {
		t.Fatalf("Failed to parse metadata: %v", err)
	}

	if metadata.RecordType != "A" {
		t.Errorf("Expected record type 'A', got '%s'", metadata.RecordType)
	}

	if metadata.RecordValue != testIP {
		t.Errorf("Expected record value '%s', got '%s'", testIP, metadata.RecordValue)
	}

	t.Logf("✓ A record created successfully: %s", metadata.RecordName)
	t.Logf("  - Artifact ID: %d", artifact.ID)
	t.Logf("  - Record Name: %s", metadata.RecordName)
	t.Logf("  - Record Value: %s", metadata.RecordValue)
	t.Logf("  - TTL: %d", metadata.TTL)

	// Verify audit log was created
	var auditLogs []storage.AuditLog
	db.Where("action = ?", "create_dns_artifact").Find(&auditLogs)
	if len(auditLogs) == 0 {
		t.Error("Expected audit log entry, found none")
	} else {
		t.Logf("✓ Audit log created: %s", auditLogs[0].Action)
	}

	// Clean up
	t.Log("Cleaning up test artifact...")
	err = generator.DeleteDNSRecord(ctx, artifact)
	if err != nil {
		t.Errorf("Failed to delete test artifact: %v", err)
	} else {
		t.Log("✓ Test artifact deleted successfully")
	}
}

// TestGenerateCNAMERecord tests CNAME record creation
func TestGenerateCNAMERecord(t *testing.T) {
	if os.Getenv("CLOUDFLARE_API_TOKEN") == "" {
		t.Skip("Skipping live DNS test: CLOUDFLARE_API_TOKEN not set")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db := setupTestDB(t)

	generator, err := artifacts.NewDNSGenerator(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create DNS generator: %v", err)
	}

	ctx := context.Background()
	testName := "test-cname"
	testTarget := "target.example.com"

	t.Logf("Creating CNAME record: %s.%s.%s -> %s", testName, cfg.LabSubdomain, cfg.LabDomain, testTarget)

	artifact, err := generator.GenerateCNAMERecord(ctx, testName, testTarget)
	if err != nil {
		t.Fatalf("Failed to generate CNAME record: %v", err)
	}

	// Parse metadata
	var metadata artifacts.DNSArtifactMetadata
	err = json.Unmarshal([]byte(artifact.Metadata), &metadata)
	if err != nil {
		t.Fatalf("Failed to parse metadata: %v", err)
	}

	if metadata.RecordType != "CNAME" {
		t.Errorf("Expected record type 'CNAME', got '%s'", metadata.RecordType)
	}

	t.Logf("✓ CNAME record created successfully: %s -> %s", metadata.RecordName, metadata.RecordValue)

	// Clean up
	err = generator.DeleteDNSRecord(ctx, artifact)
	if err != nil {
		t.Errorf("Failed to delete test artifact: %v", err)
	}
}

// TestGenerateTXTRecord tests TXT record creation
func TestGenerateTXTRecord(t *testing.T) {
	if os.Getenv("CLOUDFLARE_API_TOKEN") == "" {
		t.Skip("Skipping live DNS test: CLOUDFLARE_API_TOKEN not set")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db := setupTestDB(t)

	generator, err := artifacts.NewDNSGenerator(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create DNS generator: %v", err)
	}

	ctx := context.Background()
	testName := "test-txt"
	testContent := "cloak-test-record-v1"

	t.Logf("Creating TXT record: %s.%s.%s -> \"%s\"", testName, cfg.LabSubdomain, cfg.LabDomain, testContent)

	artifact, err := generator.GenerateTXTRecord(ctx, testName, testContent)
	if err != nil {
		t.Fatalf("Failed to generate TXT record: %v", err)
	}

	// Parse metadata
	var metadata artifacts.DNSArtifactMetadata
	err = json.Unmarshal([]byte(artifact.Metadata), &metadata)
	if err != nil {
		t.Fatalf("Failed to parse metadata: %v", err)
	}

	if metadata.RecordType != "TXT" {
		t.Errorf("Expected record type 'TXT', got '%s'", metadata.RecordType)
	}

	if metadata.RecordValue != testContent {
		t.Errorf("Expected record value '%s', got '%s'", testContent, metadata.RecordValue)
	}

	t.Logf("✓ TXT record created successfully: %s", metadata.RecordName)

	// Clean up
	err = generator.DeleteDNSRecord(ctx, artifact)
	if err != nil {
		t.Errorf("Failed to delete test artifact: %v", err)
	}
}

// TestGenerateRandomArtifact tests random artifact generation
func TestGenerateRandomArtifact(t *testing.T) {
	if os.Getenv("CLOUDFLARE_API_TOKEN") == "" {
		t.Skip("Skipping live DNS test: CLOUDFLARE_API_TOKEN not set")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db := setupTestDB(t)

	generator, err := artifacts.NewDNSGenerator(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create DNS generator: %v", err)
	}

	ctx := context.Background()

	// Test all record types
	recordTypes := []artifacts.DNSRecordType{
		artifacts.RecordTypeA,
		artifacts.RecordTypeCNAME,
		artifacts.RecordTypeTXT,
	}

	for _, recordType := range recordTypes {
		t.Run(string(recordType), func(t *testing.T) {
			artifact, err := generator.GenerateRandomArtifact(ctx, recordType)
			if err != nil {
				t.Fatalf("Failed to generate random %s artifact: %v", recordType, err)
			}

			var metadata artifacts.DNSArtifactMetadata
			err = json.Unmarshal([]byte(artifact.Metadata), &metadata)
			if err != nil {
				t.Fatalf("Failed to parse metadata: %v", err)
			}

			if metadata.RecordType != string(recordType) {
				t.Errorf("Expected record type '%s', got '%s'", recordType, metadata.RecordType)
			}

			t.Logf("✓ Random %s record created: %s -> %s", recordType, metadata.RecordName, metadata.RecordValue)

			// Clean up
			err = generator.DeleteDNSRecord(ctx, artifact)
			if err != nil {
				t.Errorf("Failed to delete test artifact: %v", err)
			}
		})
	}
}

// TestListActiveDNSArtifacts tests listing active artifacts
func TestListActiveDNSArtifacts(t *testing.T) {
	if os.Getenv("CLOUDFLARE_API_TOKEN") == "" {
		t.Skip("Skipping live DNS test: CLOUDFLARE_API_TOKEN not set")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db := setupTestDB(t)

	generator, err := artifacts.NewDNSGenerator(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create DNS generator: %v", err)
	}

	ctx := context.Background()

	// Create multiple artifacts
	artifact1, err := generator.GenerateARecord(ctx, "test-list-1", "203.0.113.1")
	if err != nil {
		t.Fatalf("Failed to create artifact 1: %v", err)
	}

	artifact2, err := generator.GenerateTXTRecord(ctx, "test-list-2", "test-content")
	if err != nil {
		t.Fatalf("Failed to create artifact 2: %v", err)
	}

	// List active artifacts
	artifacts, err := generator.ListActiveDNSArtifacts()
	if err != nil {
		t.Fatalf("Failed to list artifacts: %v", err)
	}

	if len(artifacts) < 2 {
		t.Errorf("Expected at least 2 artifacts, got %d", len(artifacts))
	}

	t.Logf("✓ Listed %d active DNS artifacts", len(artifacts))
	for _, artifact := range artifacts {
		metadata, err := generator.GetDNSRecordInfo(&artifact)
		if err == nil {
			t.Logf("  - %s: %s (%s)", metadata.RecordType, metadata.RecordName, artifact.Status)
		}
	}

	// Clean up
	generator.DeleteDNSRecord(ctx, artifact1)
	generator.DeleteDNSRecord(ctx, artifact2)
}

// TestGetDNSRecordInfo tests parsing artifact metadata
func TestGetDNSRecordInfo(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db := setupTestDB(t)

	generator, err := artifacts.NewDNSGenerator(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create DNS generator: %v", err)
	}

	// Create a mock artifact with metadata
	metadata := artifacts.DNSArtifactMetadata{
		RecordType:  "A",
		RecordName:  "test.cloak.example.com",
		RecordValue: "203.0.113.1",
		TTL:         120,
		Proxied:     false,
		ZoneID:      "test-zone-id",
		ZoneName:    "example.com",
	}

	metadataJSON, _ := json.Marshal(metadata)

	artifact := &storage.Artifact{
		Type:     "dns",
		Metadata: string(metadataJSON),
	}

	// Test parsing
	parsedMetadata, err := generator.GetDNSRecordInfo(artifact)
	if err != nil {
		t.Fatalf("Failed to parse metadata: %v", err)
	}

	if parsedMetadata.RecordType != metadata.RecordType {
		t.Errorf("Expected record type '%s', got '%s'", metadata.RecordType, parsedMetadata.RecordType)
	}

	if parsedMetadata.RecordName != metadata.RecordName {
		t.Errorf("Expected record name '%s', got '%s'", metadata.RecordName, parsedMetadata.RecordName)
	}

	if parsedMetadata.RecordValue != metadata.RecordValue {
		t.Errorf("Expected record value '%s', got '%s'", metadata.RecordValue, parsedMetadata.RecordValue)
	}

	t.Log("✓ Metadata parsed correctly")
	t.Logf("  - Type: %s", parsedMetadata.RecordType)
	t.Logf("  - Name: %s", parsedMetadata.RecordName)
	t.Logf("  - Value: %s", parsedMetadata.RecordValue)
	t.Logf("  - TTL: %d", parsedMetadata.TTL)
}

// TestDeleteDNSRecord tests artifact deletion
func TestDeleteDNSRecord(t *testing.T) {
	if os.Getenv("CLOUDFLARE_API_TOKEN") == "" {
		t.Skip("Skipping live DNS test: CLOUDFLARE_API_TOKEN not set")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db := setupTestDB(t)

	generator, err := artifacts.NewDNSGenerator(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create DNS generator: %v", err)
	}

	ctx := context.Background()

	// Create artifact
	artifact, err := generator.GenerateARecord(ctx, "test-delete", "203.0.113.1")
	if err != nil {
		t.Fatalf("Failed to create artifact: %v", err)
	}

	originalStatus := artifact.Status
	t.Logf("Created artifact with status: %s", originalStatus)

	// Delete it
	err = generator.DeleteDNSRecord(ctx, artifact)
	if err != nil {
		t.Fatalf("Failed to delete artifact: %v", err)
	}

	// Verify status changed
	var updatedArtifact storage.Artifact
	db.First(&updatedArtifact, artifact.ID)

	if updatedArtifact.Status != "deleted" {
		t.Errorf("Expected status 'deleted', got '%s'", updatedArtifact.Status)
	}

	t.Log("✓ Artifact deleted successfully")
	t.Logf("  - Status changed from '%s' to '%s'", originalStatus, updatedArtifact.Status)

	// Verify delete audit log
	var auditLogs []storage.AuditLog
	db.Where("action = ?", "delete_dns_artifact").Find(&auditLogs)
	if len(auditLogs) == 0 {
		t.Error("Expected delete audit log entry, found none")
	} else {
		t.Logf("✓ Delete audit log created")
	}
}

// TestDNSRecordMetadataJSON tests JSON serialization/deserialization
func TestDNSRecordMetadataJSON(t *testing.T) {
	metadata := artifacts.DNSArtifactMetadata{
		RecordType:  "A",
		RecordName:  "test.example.com",
		RecordValue: "192.0.2.1",
		TTL:         300,
		Proxied:     true,
		ZoneID:      "zone123",
		ZoneName:    "example.com",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("Failed to marshal metadata: %v", err)
	}

	t.Logf("JSON: %s", string(jsonData))

	// Unmarshal back
	var decoded artifacts.DNSArtifactMetadata
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}

	// Verify all fields
	if decoded.RecordType != metadata.RecordType {
		t.Errorf("RecordType mismatch: expected '%s', got '%s'", metadata.RecordType, decoded.RecordType)
	}
	if decoded.RecordName != metadata.RecordName {
		t.Errorf("RecordName mismatch: expected '%s', got '%s'", metadata.RecordName, decoded.RecordName)
	}
	if decoded.RecordValue != metadata.RecordValue {
		t.Errorf("RecordValue mismatch: expected '%s', got '%s'", metadata.RecordValue, decoded.RecordValue)
	}
	if decoded.TTL != metadata.TTL {
		t.Errorf("TTL mismatch: expected %d, got %d", metadata.TTL, decoded.TTL)
	}
	if decoded.Proxied != metadata.Proxied {
		t.Errorf("Proxied mismatch: expected %v, got %v", metadata.Proxied, decoded.Proxied)
	}

	t.Log("✓ JSON serialization/deserialization works correctly")
}

// TestDNSTeardownInfoJSON tests teardown info serialization
func TestDNSTeardownInfoJSON(t *testing.T) {
	teardownInfo := artifacts.DNSTeardownInfo{
		CloudflareRecordID: "cf-record-123",
		ZoneID:             "zone-456",
		RecordName:         "test.example.com",
		RecordType:         "A",
		CreatedAt:          time.Now(),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(teardownInfo)
	if err != nil {
		t.Fatalf("Failed to marshal teardown info: %v", err)
	}

	t.Logf("Teardown JSON: %s", string(jsonData))

	// Unmarshal back
	var decoded artifacts.DNSTeardownInfo
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal teardown info: %v", err)
	}

	if decoded.CloudflareRecordID != teardownInfo.CloudflareRecordID {
		t.Errorf("CloudflareRecordID mismatch")
	}
	if decoded.ZoneID != teardownInfo.ZoneID {
		t.Errorf("ZoneID mismatch")
	}
	if decoded.RecordName != teardownInfo.RecordName {
		t.Errorf("RecordName mismatch")
	}

	t.Log("✓ Teardown info JSON serialization works correctly")
}
