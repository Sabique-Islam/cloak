package utils

import (
	"time"
)

// Calculate time between artifact creation and first detection
func CalculateLatency(createdAt, detectedAt time.Time) time.Duration {
	return detectedAt.Sub(createdAt)
}

// Calculate how long an artifact remains detectable
func CalculatePersistence(firstSeen, lastSeen time.Time) time.Duration {
	return lastSeen.Sub(firstSeen)
}
