
-- Table: artifacts
-- Stores GitHub commits, repos, and DNS records created for testing
CREATE TABLE IF NOT EXISTS artifacts (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL, -- 'github_commit', 'github_repo', 'dns'
    name TEXT NOT NULL, -- Human-readable name
    subject TEXT NOT NULL, -- GitHub repo URL or DNS record name
    value TEXT NOT NULL, -- Commit SHA, repo name, or DNS value
    published_at TEXT NOT NULL, -- ISO8601 timestamp when artifact was created
    metadata TEXT, -- JSON with provider-specific details (e.g., DNS record type, TTL, commit author)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table: detections
-- Logs when artifacts are found in OSINT feeds (GitHub Search, VirusTotal, etc.)
CREATE TABLE IF NOT EXISTS detections (
    id TEXT PRIMARY KEY,
    artifact_id TEXT NOT NULL,
    provider TEXT NOT NULL, -- 'github_search', 'virustotal', 'greynoise', 'otx'
    detected_at TEXT NOT NULL, -- ISO8601 timestamp when detection occurred
    raw_response TEXT, -- JSON raw response from provider for analysis
    metadata TEXT, -- JSON normalized/extracted fields
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (artifact_id) REFERENCES artifacts(id) ON DELETE CASCADE
);

-- Table: jobs
-- Manages background tasks for artifact generation and OSINT ingestion
CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL, -- 'generate_artifact', 'ingest_osint', 'teardown_artifact'
    status TEXT NOT NULL DEFAULT 'pending', -- 'pending', 'running', 'completed', 'failed'
    artifact_id TEXT, -- Optional: linked artifact for ingestion/teardown jobs
    config TEXT, -- JSON job configuration
    result TEXT, -- JSON result or error message
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (artifact_id) REFERENCES artifacts(id) ON DELETE SET NULL
);

-- Table: audit_logs
-- Tracks every action with user, timestamp, and reason for full accountability
CREATE TABLE IF NOT EXISTS audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL, -- 'create_artifact', 'query_provider', 'teardown', etc.
    entity_type TEXT, -- 'artifact', 'detection', 'job'
    entity_id TEXT, -- ID of the related entity
    user TEXT, -- User or system component that triggered the action
    details TEXT, -- JSON additional context
    reason TEXT, -- Why this action was taken
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_artifacts_type ON artifacts(type);
CREATE INDEX IF NOT EXISTS idx_artifacts_published_at ON artifacts(published_at);
CREATE INDEX IF NOT EXISTS idx_detections_artifact_id ON detections(artifact_id);
CREATE INDEX IF NOT EXISTS idx_detections_provider ON detections(provider);
CREATE INDEX IF NOT EXISTS idx_detections_detected_at ON detections(detected_at);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_type ON jobs(type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
