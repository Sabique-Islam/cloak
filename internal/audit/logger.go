package audit

import (
	"log"
)

// For structured audit logging
type Logger struct {
	
}

// Creates new audit logger
func NewLogger() *Logger {
	return &Logger{}
}

// Records an auditable action
func (l *Logger) LogAction(actor, action, resource string, details map[string]interface{}) {
	
	log.Printf("AUDIT: actor=%s action=%s resource=%s", actor, action, resource) 
}
