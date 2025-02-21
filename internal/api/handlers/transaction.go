package handlers

import (
	"net/http"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	transactions, err := h.transactionService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	transactionIDStr := c.Param("id")
	id, err := uuid.Parse(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	transaction, err := h.transactionService.GetByID(c, id)
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

	newTransaction.ID = uuid.New()
	if err := h.transactionService.Create(c, &newTransaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}
	c.JSON(http.StatusCreated, newTransaction)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	transactionIDStr := c.Param("id")
	id, err := uuid.Parse(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	var updatedTransaction models.Transaction
	if err := c.ShouldBindJSON(&updatedTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updatedTransaction.ID = id
	if err := h.transactionService.Update(c, &updatedTransaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update transaction"})
		return
	}
	c.JSON(http.StatusOK, updatedTransaction)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	transactionIDStr := c.Param("id")
	id, err := uuid.Parse(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	if err := h.transactionService.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted"})
}
