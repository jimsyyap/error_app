package database

import (
	"log"

	"github.com/jimsyyap/error_app/backend/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// BuildDSN builds a PostgreSQL DSN string from the config
func (c *Config) BuildDSN() string {
	return "host=" + c.Host + 
		" user=" + c.User + 
		" password=" + c.Password + 
		" dbname=" + c.DBName + 
		" port=" + c.Port + 
		" sslmode=" + c.SSLMode + 
		" TimeZone=UTC"
}

// DB is the database instance
type DB struct {
	*gorm.DB
}

// New creates a new database connection
func New(config *Config) (*DB, error) {
	dsn := config.BuildDSN()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	
	return &DB{db}, nil
}

// Initialize sets up the database schema and seeds initial data
func (db *DB) Initialize() error {
	// Enable UUID extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")

	// AutoMigrate creates/updates tables based on the model structs
	err := db.AutoMigrate(
		&models.User{},
		&models.MatchSession{},
		&models.ErrorType{},
		&models.ErrorLog{},
	)
	if err != nil {
		return err
	}

	// Add check constraints that GORM doesn't handle automatically
	db.Exec("ALTER TABLE match_sessions DROP CONSTRAINT IF EXISTS chk_end_after_start")
	db.Exec("ALTER TABLE match_sessions ADD CONSTRAINT chk_end_after_start CHECK (end_time IS NULL OR end_time >= start_time)")
	
	db.Exec("ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_email_format")
	db.Exec("ALTER TABLE users ADD CONSTRAINT chk_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$')")

	// Seed error types if table is empty
	var count int64
	db.Model(&models.ErrorType{}).Count(&count)
	if count == 0 {
		log.Println("Seeding error types...")
		errorTypes := []models.ErrorType{
			{Name: "Forehand"},
			{Name: "Backhand"},
			{Name: "Serve"},
			{Name: "Volley"},
			{Name: "Return"},
			{Name: "Drop Shot"},
			{Name: "Overhead"},
			{Name: "Footwork"},
		}
		result := db.Create(&errorTypes)
		if result.Error != nil {
			return result.Error
		}
		log.Println("Seeded", result.RowsAffected, "error types")
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
