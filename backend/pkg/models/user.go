package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user of the Tennis Error Tracker application
type User struct {
	UserID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"user_id"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	LastLogin    *time.Time `json:"last_login"`
	Sessions     []MatchSession `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}
	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (u *User) UpdateLastLogin(tx *gorm.DB) error {
	now := time.Now()
	u.LastLogin = &now
	return tx.Model(u).Update("last_login", now).Error
}
