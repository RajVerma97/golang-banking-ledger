package routes

import (
	"github.com/RajVerma97/golang-banking-ledger/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func AccountRoutes(r *gin.Engine, accountHandler *handlers.AccountHandler, transactionHandler *handlers.TransactionHandler) {
	r.GET("/account", accountHandler.GetAccounts)
	r.GET("/account/:id", accountHandler.GetAccountByID)
	r.POST("/account", accountHandler.CreateAccount)
	r.PATCH("/account/:id", accountHandler.UpdateAccount)
	r.DELETE("/account/:id", accountHandler.DeleteAccount)
	r.GET("/accounts/:accountID/transactions", transactionHandler.GetTransactionHistory)
}
