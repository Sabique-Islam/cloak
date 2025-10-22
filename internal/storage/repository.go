package storage

// CRUD operations for domain models
type Repository struct {
	//Add DB connection
}

// Creates new repo instance
func NewRepository() *Repository {
	return &Repository{}
}

// Store new artifact
func (r *Repository) CreateArtifact(artifact *Artifact) error {
	return nil
}

// Retrieves artifacts with filters
func (r *Repository) GetArtifacts(filters map[string]interface{}) ([]Artifact, error) {
	return nil, nil
}

//Stores a new detection
func (r *Repository) CreateDetection(detection *Detection) error {
	return nil
}

// Retrieves detections with filters
func (r *Repository) GetDetections(filters map[string]interface{}) ([]Detection, error) {
	return nil, nil
}
