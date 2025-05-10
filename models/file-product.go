package models

import (
	"time"

	"gorm.io/gorm"
)

// File represents the files table in the database
type FileProduct struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename     string         `gorm:"type:varchar(255);not null" json:"filename"`
	OriginalName string         `gorm:"type:varchar(255);not null" json:"originalname"`
	MimeType     string         `gorm:"type:varchar(150);not null" json:"mimetype"`
	Path         string         `gorm:"type:varchar(500);not null" json:"path"`
	Size         int64          `gorm:"not null" json:"size"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
