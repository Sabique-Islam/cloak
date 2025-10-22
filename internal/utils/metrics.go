package utils

// Computed experiment metrics
type Metrics struct {
	TotalArtifacts     int
	DetectedArtifacts  int
	DetectionRate      float64
	AverageLatency     float64
	ProviderCoverage   map[string]int
}

// Calculate aggregate metrics from detections
func ComputeMetrics() (*Metrics, error) {
	return &Metrics{}, nil
}
