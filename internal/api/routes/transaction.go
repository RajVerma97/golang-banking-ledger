package routes

import (
	"github.com/RajVerma97/golang-banking-ledger/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func TransactionRoutes(r *gin.Engine, transactionHandler *handlers.TransactionHandler) {
	r.GET("/transaction/:id", transactionHandler.GetTransactionByID)
	r.POST("/transaction", transactionHandler.CreateTransaction)

}
