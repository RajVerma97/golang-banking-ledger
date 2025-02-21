package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTransactions(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "get all transactions ",
	})
}

func GetTransactionByID(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "get transaction by ID",
	})
}

func CreateTransaction(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "create transaction",
	})
}

func UpdateTransaction(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "update transaction",
	})
}
func DeleteTransaction(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "delete transaction",
	})
}
