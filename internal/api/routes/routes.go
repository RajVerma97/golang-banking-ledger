package routes

import (
	"github.com/RajVerma97/golang-banking-ledger/internal/api/handlers"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, accountService *service.AccountService, transactionService *service.TransactionService) {
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService, accountService)
	AccountRoutes(r, accountHandler, transactionHandler)
	TransactionRoutes(r, transactionHandler)
}
