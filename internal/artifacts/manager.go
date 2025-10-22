package artifacts

// Artifact creation jobs
type Manager struct {
	//Add dependencies (DB, GitHub client, Cloudflare client)
}

// Creates new artifact manager
func NewManager() *Manager {
	return &Manager{}
}

// Creates GitHub-based artifacts
func (m *Manager) GenerateGitHubArtifacts() error {
	return nil
}

//Creates DNS-based artifacts
func (m *Manager) GenerateDNSArtifacts() error {
	return nil
}
