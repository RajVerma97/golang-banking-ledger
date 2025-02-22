package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
	accountService     *service.AccountService
}

func NewTransactionHandler(transactionService *service.TransactionService, accountService *service.AccountService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService, accountService: accountService}
}
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	transactionIDStr := c.Param("id")

	transaction, err := h.transactionService.GetByID(c, transactionIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
func (h *TransactionHandler) GetTransactionByAccountID(c *gin.Context) {
	accountIDStr := c.Param("id")

	transaction, err := h.transactionService.GetByAccountID(c, accountIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var newTransaction models.Transaction

	if err := c.ShouldBindJSON(&newTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if newTransaction.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the amount should be greater than 0 "})
		return
	}
	accountUUID, err := uuid.Parse(newTransaction.AccountID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID format"})
		return
	}

	var account models.Account
	if account, err = h.accountService.GetByID(context.Background(), accountUUID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	if newTransaction.Type == models.WITHDRAWL && account.Balance < newTransaction.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}
	h.initializeTransaction(&newTransaction, accountUUID)

	if err := h.transactionService.Create(c, &newTransaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	if err := h.transactionService.PublishTransactionEvent(c, &newTransaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish transaction event"})
		return
	}
	c.JSON(http.StatusCreated, newTransaction)
}

func (h *TransactionHandler) initializeTransaction(transaction *models.Transaction, accountUUID uuid.UUID) {
	transactionUUID := uuid.New()
	transaction.ID = transactionUUID.String()
	transaction.AccountID = accountUUID.String()
	transaction.Status = models.PENDING
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

}
