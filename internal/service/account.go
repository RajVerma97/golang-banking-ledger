package service

import (
	"context"
	"sync"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/google/uuid"
)

type AccountService struct {
	accountRepo AccountRepository
	mutex       sync.RWMutex
}

type AccountRepository interface {
	GetAll(ctx context.Context) (models.Accounts, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Account, error)
	Create(ctx context.Context, account *models.Account) error
	Update(ctx context.Context, id uuid.UUID, updates models.AccountUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewAccountService(accountRepo AccountRepository) *AccountService {

	return &AccountService{
		accountRepo: accountRepo,
	}
}
