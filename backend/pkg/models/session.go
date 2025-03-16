package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MatchSession represents a tennis match session
type MatchSession struct {
	SessionID    uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"session_id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"-"`
	StartTime    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	OpponentName *string   `gorm:"type:varchar(100)" json:"opponent_name,omitempty"`
	Location     *string   `gorm:"type:varchar(100)" json:"location,omitempty"`
	Score        *string   `gorm:"type:varchar(50)" json:"score,omitempty"`
	Notes        *string   `gorm:"type:text" json:"notes,omitempty"`
	ErrorLogs    []ErrorLog `gorm:"foreignKey:SessionID" json:"error_logs,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (s *MatchSession) BeforeCreate(tx *gorm.DB) error {
	if s.SessionID == uuid.Nil {
		s.SessionID = uuid.New()
	}
	return nil
}

// IsActive checks if a session is currently active (not ended)
func (s *MatchSession) IsActive() bool {
	return s.EndTime == nil
}

// End marks a session as ended
func (s *MatchSession) End(tx *gorm.DB) error {
	now := time.Now()
	s.EndTime = &now
	return tx.Model(s).Update("end_time", now).Error
}
