package db

import (
	"fmt"
	"log"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() *gorm.DB {
	dsn := "host=localhost user=admin password=secret dbname=bank-ledger sslmode=disable port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(" Failed to connect to PostgreSQL: ", err)
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
