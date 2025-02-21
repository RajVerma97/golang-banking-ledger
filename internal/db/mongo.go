package db

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func InitMongo() (*mongo.Client, *mongo.Database, *mongo.Collection, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(" Warning: No .env file found, using system environment variables.")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal(" MONGODB_URI is not set in environment variables")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(" Failed to connect to MongoDB: ", err)
		return nil, nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(" MongoDB connection test failed: ", err)
		return nil, nil, nil, err
	}

	fmt.Println(" Connected to MongoDB")

	db := client.Database("banking")
	collection := db.Collection("transactions")

	return client, db, collection, nil
}
