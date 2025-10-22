package ingestion

import (
	"context"
)

// Worker handles concurrent OSINT ingestion with rate limiting
type Worker struct {
}

// Creates new ingestion worker
func NewWorker() *Worker {
	return &Worker{}
}

// start ingestion worker
func (w *Worker) Start(ctx context.Context) error {
	// Implement concurrent queries with errgroup
	// Rate limits
	// Store detections in DB
	return nil
}
