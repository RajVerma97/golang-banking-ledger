package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RajVerma97/golang-banking-ledger/internal/api/routes"
	"github.com/RajVerma97/golang-banking-ledger/internal/db"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/mongodb"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/postgres"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	"github.com/RajVerma97/golang-banking-ledger/pkg/middleware"
	"github.com/RajVerma97/golang-banking-ledger/pkg/queue"
	"github.com/RajVerma97/golang-banking-ledger/pkg/worker"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger, err := middleware.InitLogger()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	logger.Info("Server starting",
		zap.String("port", PORT),
		zap.String("mode", gin.Mode()),
	)

	postgresDB, _ := db.InitPostgres()
	mongoDB, _, _ := db.InitMongo()

	accountRepo := postgres.NewAccountRepository(postgresDB)
	transactionRepo := mongodb.NewTransactionRepository(mongoDB)

	rabbitMQConn, rabbitMQChannel, err := queue.InitRabbitMQ()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitMQConn.Close()
	defer rabbitMQChannel.Close()

	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo, rabbitMQChannel)

	go func() {
		worker := worker.NewTransactionWorker(rabbitMQChannel, accountRepo, transactionRepo)
		worker.ProcessTransactions()
	}()

	routes.Setup(router, accountService, transactionService)

	fmt.Printf("Server Listening on Port %s\n", PORT)
	router.Run(":" + PORT)
}
