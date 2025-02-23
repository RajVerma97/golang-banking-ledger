package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo() (*mongo.Database, *mongo.Collection, error) {
	var mongoURI string

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}

	mongoURI = os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, nil, fmt.Errorf("MONGO_URI is not set in environment variables")
	}

	var client *mongo.Client
	var err error
	maxRetries := 5
	retryDelay := time.Second * 3

	for i := 0; i < maxRetries; i++ {
		clientOptions := options.Client().ApplyURI(mongoURI)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			err = client.Ping(ctx, nil)
			if err == nil {
				cancel()
				log.Printf("Connected to MongoDB successfully after %d attempts", i+1)
				break
			}
		}

		cancel()
		if i < maxRetries-1 {
			log.Printf("Failed to connect to MongoDB (attempt %d/%d): %v. Retrying in %v...",
				i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB after %d attempts: %v",
			maxRetries, err)
	}

	db := client.Database("banking_ledger")
	collection := db.Collection("transactions")

	return db, collection, nil
}
