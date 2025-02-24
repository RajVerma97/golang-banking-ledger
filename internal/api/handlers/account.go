package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	service service.AccountServiceInterface
}

func NewAccountHandler(service service.AccountServiceInterface) *AccountHandler {
	return &AccountHandler{
		service: service,
	}
}
func (accountHandler *AccountHandler) GetAccounts(c *gin.Context) {

	accounts, err := accountHandler.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch accounts"})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func (accountHandler *AccountHandler) GetAccountByID(c *gin.Context) {

	accountIDStr := c.Param("id")
	id, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}
	account, err := accountHandler.service.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}
	c.JSON(http.StatusOK, account)
}
func (accountHandler *AccountHandler) CreateAccount(c *gin.Context) {
	var newAccountRequest models.AccountCreate

	if err := c.ShouldBindJSON(&newAccountRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	newAccount := models.Account{
		ID:            uuid.New(),
		AccountNumber: generateAccountNumber(),
		FirstName:     newAccountRequest.FirstName,
		LastName:      newAccountRequest.LastName,
		Email:         newAccountRequest.Email,
		Phone:         newAccountRequest.Phone,
		Balance:       newAccountRequest.Balance,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := accountHandler.service.Create(c.Request.Context(), &newAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, newAccount)
}

func generateAccountNumber() int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Intn(900000) + 100000
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

	if updateData.FirstName == nil && updateData.LastName == nil && updateData.Phone == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one valid field must be provided"})
		return
	}

	if err := accountHandler.service.Update(c.Request.Context(), id, updateData); err != nil {
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

	if err := accountHandler.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}
