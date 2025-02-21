package main

import (
	"fmt"

	"github.com/RajVerma97/golang-banking-ledger/internal/api/routes"
	"github.com/RajVerma97/golang-banking-ledger/internal/db"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/mongodb"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/postgres"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	postgresDB := db.InitPostgres()
	mongoDB, _ := db.InitMongo()

	accountRepo := postgres.NewAccountRepository(postgresDB)
	transactionRepo := mongodb.NewTransactionRepository(mongoDB)

	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(transactionRepo)

	routes.Setup(r, accountService, transactionService)
	PORT := 3000
	fmt.Printf("Server Listening on Port %d", PORT)
	r.Run(fmt.Sprintf(":%d", PORT))
}
