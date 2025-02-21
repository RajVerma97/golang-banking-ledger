package handlers

import (
	"net/http"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}
func (accountHandler *AccountHandler) GetAccounts(c *gin.Context) {

	accounts, err := accountHandler.accountService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch accounts"})

	}
	c.JSON(http.StatusOK, accounts)
}

func (accountHandler *AccountHandler) GetAccountByID(c *gin.Context) {

	accountIDStr := c.Param["id"]
	id, err := uuid.Parse("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}
	account, err := accountHandler.accountService.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}
	c.JSON(http.StatusOK, account)
}
func (accountHandler *AccountHandler) CreateAccount(c *gin.Context) {

	var newAccount models.Account

	if err := c.ShouldBindJSON(&newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	newAccount.ID = uuid.New()
	if err := accountHandler.accountService.Create(c, &newAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}
	c.JSON(http.StatusCreated, newAccount)

}
func (accountHandler *AccountHandler) UpdateAccount(c *gin.Context) {
	accountIDStr := c.Param("id")
	id, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var updateData models.AccountUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := accountHandler.accountService.Update(c, id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account updated successfully"})
}

func (accountHandler *AccountHandler) DeleteAccount(c *gin.Context) {
	accountIDStr := c.Param("id")
	id, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	if err := accountHandler.accountService.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}
