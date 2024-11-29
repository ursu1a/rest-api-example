package db

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                     uint   `gorm:"primaryKey"`
	Email                  string `gorm:"unique;not null"`
	Name                   string `gorm:"not null"`
	GoogleID               *string
	PasswordHash           string
	EmailVerified          bool
	EmailVerificationToken string
	Picture                string
	CreatedAt              time.Time `gorm:"default:current_timestamp"`
	UpdatedAt              time.Time `gorm:"default:current_timestamp"`
}

func AddGoogleIDUniqueConstraint(db *gorm.DB) error {
	return db.Exec(`
		 CREATE UNIQUE INDEX IF NOT EXISTS idx_users_google_id_not_null 
		 ON users (google_id)
		 WHERE google_id IS NOT NULL
	`).Error
}
