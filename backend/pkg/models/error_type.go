package models

// ErrorType represents a type of tennis error (e.g., Forehand, Backhand)
type ErrorType struct {
	ErrorTypeID int    `gorm:"primaryKey;autoIncrement" json:"error_type_id"`
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
}
