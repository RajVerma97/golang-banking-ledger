package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAccounts(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "get all accounts",
	})
}

func GetAccountByID(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "get account by ID",
	})
}
func CreateAccount(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "create account",
	})
}

func UpdateAccount(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": " update account",
	})
}

func DeleteAccount(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "delete account",
	})
}
