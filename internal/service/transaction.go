package service

import (
	"context"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
)

type TransactionService struct {
	transactionRepo TransactionRepository
	accountRepo     AccountRepository
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	GetByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error)
}

func NewTransactionService(transactionRepo TransactionRepository, accountRepo AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (ts *TransactionService) Create(ctx context.Context, tx *models.Transaction) error {
	return ts.transactionRepo.Create(ctx, tx)
}

func (ts *TransactionService) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	return ts.transactionRepo.GetByID(ctx, id)
}

func (ts *TransactionService) GetByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error) {
	return ts.transactionRepo.GetByAccountID(ctx, accountID)
}
