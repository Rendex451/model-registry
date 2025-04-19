package models

import "time"

type Model struct {
	ID          int       `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"size:1024"`
	Type        string    `gorm:"size:255;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type ModelVersion struct {
	ID          int       `gorm:"primaryKey"`
	ModelID     int       `gorm:"not null"`
	VersionName string    `gorm:"size:255;not null"`
	FilePath    string    `gorm:"size:1024;not null"`
	DateAdded   time.Time `gorm:"autoCreateTime"`
}
