package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ErrorType represents an entry in the Error_Types table
type ErrorType struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve database connection details from environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// AutoMigrate ensures the table exists
	db.AutoMigrate(&ErrorType{})

	// Define predefined error types
	errorTypes := []ErrorType{
		{Name: "Forehand"},
		{Name: "Backhand"},
		{Name: "Serve"},
		{Name: "Volley"},
	}

	// Seed the database (if the entries don't already exist)
	for _, et := range errorTypes {
		db.FirstOrCreate(&et, ErrorType{Name: et.Name})
	}

	fmt.Println("✅ Error types seeded successfully!")
}
