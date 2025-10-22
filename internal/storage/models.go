package storage

import (
	"time"

	"gorm.io/gorm"
)

// GPT generated (Check later)

type Artifact struct {
	ID           uint           `gorm:"primarykey"`
	CreatedAt    time.Time      `gorm:"index"`
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Type         string         `gorm:"index;not null"`
	Identifier   string         `gorm:"uniqueIndex;not null"`
	Metadata     string         `gorm:"type:json"`
	TeardownInfo string         `gorm:"type:json"`
	Status       string         `gorm:"default:'active'"`
}

type Detection struct {
	ID           uint           `gorm:"primarykey"`
	CreatedAt    time.Time      `gorm:"index"`
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	ArtifactID   uint           `gorm:"index;not null"`
	Artifact     Artifact       `gorm:"foreignKey:ArtifactID"`
	Provider     string         `gorm:"index;not null"`
	DetectedAt   time.Time      `gorm:"index;not null"`
	RawResponse  string         `gorm:"type:text"`
	MatchScore   float64
	FirstSeen    bool           `gorm:"default:false"`
}

type Job struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Type      string         `gorm:"index;not null"`
	Status    string         `gorm:"index;not null"`
	StartedAt *time.Time
	EndedAt   *time.Time
	Error     string         `gorm:"type:text"`
	Metadata  string         `gorm:"type:json"`
}

type AuditLog struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	Actor     string    `gorm:"index"`
	Action    string    `gorm:"index;not null"`
	Resource  string    `gorm:"index"`
	Details   string    `gorm:"type:json"`
	IPAddress string
}
