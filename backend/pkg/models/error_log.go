package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrorLog represents a logged error during a tennis match
type ErrorLog struct {
	ErrorID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"error_id"`
	SessionID   uuid.UUID `gorm:"type:uuid;not null;index" json:"session_id"`
	Session     MatchSession `gorm:"foreignKey:SessionID" json:"-"`
	ErrorTypeID int       `gorm:"not null" json:"error_type_id"`
	ErrorType   ErrorType `gorm:"foreignKey:ErrorTypeID" json:"error_type,omitempty"`
	Timestamp   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"timestamp"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (e *ErrorLog) BeforeCreate(tx *gorm.DB) error {
	if e.ErrorID == uuid.Nil {
		e.ErrorID = uuid.New()
	}
	return nil
}
