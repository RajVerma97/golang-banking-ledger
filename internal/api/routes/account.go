package routes

import (
	"github.com/RajVerma97/golang-banking-ledger/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.GET("/", handlers.Handler)
}
