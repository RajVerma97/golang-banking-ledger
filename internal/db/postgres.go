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

func InitPostgres() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		log.Fatal("POSTGRES_URI is not set in environment variables")
	}

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatal("Failed to enable uuid-ossp extension:", err)
	}
	fmt.Println("Running Migrations...")
	if err := db.AutoMigrate(&models.Account{}); err != nil {
		log.Fatal("Migration failed: ", err)
	}
	fmt.Println(" Connected to PostgreSQL")
	return db
}
