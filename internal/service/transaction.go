package service

import (
	"context"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo TransactionRepository
	accountRepo     AccountRepository
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Transaction, error)
}

func NewTransactionService(transactionRepo TransactionRepository, accountRepo AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}
