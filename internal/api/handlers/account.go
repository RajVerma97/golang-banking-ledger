package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "hellooo",
	})
}
