package routes

import (
	"github.com/RajVerma97/golang-banking-ledger/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func TransactionRoutes(r *gin.Engine) {
	r.GET("/transaction", handlers.GetTransactions)
	r.GET("/transaction/:id", handlers.GetTransactionByID)
	r.POST("/transaction", handlers.CreateTransaction)
	r.PATCH("/transaction/:id", handlers.UpdateTransaction)
	r.DELETE("/transaction/:id", handlers.DeleteTransaction)

}
