package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitPostgres() *gorm.DB {
	dsn := "host=localhost user=postgres password=secret dbname=bank sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(" Failed to connect to PostgreSQL: ", err)
	}
	fmt.Println(" Connected to PostgreSQL")
	return db
}
