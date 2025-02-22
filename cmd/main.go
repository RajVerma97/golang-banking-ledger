package main

import (
	"fmt"

	"github.com/RajVerma97/golang-banking-ledger/internal/api/routes"
	"github.com/RajVerma97/golang-banking-ledger/internal/db"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/mongodb"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/postgres"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	"github.com/RajVerma97/golang-banking-ledger/pkg/middleware"
	"github.com/gin-gonic/gin"
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

	logger.Info("Server starting",
		zap.String("port", "3000"),
		zap.String("mode", gin.Mode()),
	)

	postgresDB := db.InitPostgres()
	mongoDB, _, _ := db.InitMongo()

	accountRepo := postgres.NewAccountRepository(postgresDB)
	transactionRepo := mongodb.NewTransactionRepository(mongoDB)

	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)

	routes.Setup(router, accountService, transactionService)
	PORT := 3000
	fmt.Printf("Server Listening on Port %d", PORT)
	router.Run(fmt.Sprintf(":%d", PORT))
}
