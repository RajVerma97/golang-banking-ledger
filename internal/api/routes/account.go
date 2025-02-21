package routes

import (
	"github.com/RajVerma97/golang-banking-ledger/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func AccountRoutes(r *gin.Engine) {
	r.GET("/account", handlers.GetAccounts)
	r.GET("/account/:id", handlers.GetAccountByID)
	r.POST("/account", handlers.CreateAccount)
	r.PATCH("/account/:id", handlers.UpdateAccount)
	r.DELETE("/account/:id", handlers.DeleteAccount)
}
