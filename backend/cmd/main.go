package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/jimsyyap/error_app/backend/config"
	"github.com/jimsyyap/error_app/backend/internal/handlers"
	"github.com/jimsyyap/error_app/backend/internal/middleware"
	"github.com/jimsyyap/error_app/backend/pkg/models"
)

// main is the entry point of the Tennis Error Tracker backend.
// It initializes the configuration, database, router, and starts the server.
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v. Proceeding with system environment variables.", err)
	}

	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Establish database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to the database")

	// Auto-migrate database schema
	err = db.AutoMigrate(&models.User{}, &models.MatchSession{}, &models.ErrorType{}, &models.ErrorLog{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	// Seed Error_Types table if empty
	seedErrorTypes(db)

	// Initialize Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORS()) // Enable CORS for frontend integration

	// Define public routes
	router.POST("/register", handlers.Register(db))
	router.POST("/login", handlers.Login(db))

	// Define protected routes group with JWT authentication middleware
	protected := router.Group("/")
	protected.Use(middleware.Authenticate(db, cfg.JWTSecret))
	{
		protected.POST("/sessions", handlers.StartSession(db))
		protected.PUT("/sessions/:session_id", handlers.EndSession(db))
		protected.GET("/sessions", handlers.ListSessions(db))
		protected.POST("/errors", handlers.LogError(db))
		protected.DELETE("/errors/last", handlers.UndoLastError(db))
		protected.GET("/sessions/:session_id/summary", handlers.GetSummary(db))
		protected.GET("/error-types", handlers.GetErrorTypes(db))
	}

	// Determine server port
	port := cfg.Port
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Start the server
	log.Printf("Starting Tennis Error Tracker server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// seedErrorTypes populates the Error_Types table with predefined error types if it's empty.
func seedErrorTypes(db *gorm.DB) {
	var count int64
	if err := db.Model(&models.ErrorType{}).Count(&count).Error; err != nil {
		log.Fatalf("Failed to check Error_Types table: %v", err)
	}

	if count == 0 {
		errorTypes := []models.ErrorType{
			{Name: "Forehand"},
			{Name: "Backhand"},
			{Name: "Serve"},
			{Name: "Volley"},
		}
		if err := db.Create(&errorTypes).Error; err != nil {
			log.Fatalf("Failed to seed error types: %v", err)
		}
		log.Println("Successfully seeded Error_Types table")
	}
}
