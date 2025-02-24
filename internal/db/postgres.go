package db

import (
	"fmt"
	"log"
	"os"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		return nil, fmt.Errorf("POSTGRES_URI is not set in environment variables")
	}

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return nil, fmt.Errorf("failed to enable uuid-ossp extension: %w", err)
	}

	log.Println("Running Migrations...")
	if err := db.AutoMigrate(&models.Account{}); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}
