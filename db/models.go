package db

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"unique;not null"`
	Name         string `gorm:"not null"`
	GoogleID     string `gorm:"uniqueIndex:idx_users_google_id,where:google_id IS NOT NULL"`
	PasswordHash string
	RefreshToken string
	Picture      string
	CreatedAt    time.Time `gorm:"default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"default:current_timestamp"`
}
