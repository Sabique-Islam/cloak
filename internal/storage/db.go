package storage

import (
	"gorm.io/gorm"
)

// holds the DB connection
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(dbPath string) error {
	// Init SQLite with GORM
	// Run migrations
	return nil
}

// close DB connection
func Close() error {
	return nil
}
