package api

import (
	"net/http"
)

// GetArtifacts handles GET /api/artifacts
func GetArtifacts(w http.ResponseWriter, r *http.Request) {

}

// CreateArtifact handles POST /api/artifacts
func CreateArtifact(w http.ResponseWriter, r *http.Request) {

}

// GetDetections handles GET /api/detections
func GetDetections(w http.ResponseWriter, r *http.Request) {

}

// GetMetrics handles GET /api/metrics
func GetMetrics(w http.ResponseWriter, r *http.Request) {

}

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
