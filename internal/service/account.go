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
	Deposit(ctx context.Context, id uuid.UUID, amount float64) error
	Withdraw(ctx context.Context, id uuid.UUID, amount float64) error
}

func NewAccountService(accountRepo AccountRepository) *AccountService {

	return &AccountService{
		accountRepo: accountRepo,
	}
}

func (s *AccountService) GetAll(ctx context.Context) (models.Accounts, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.accountRepo.GetAll(ctx)
}

func (s *AccountService) GetByID(ctx context.Context, id uuid.UUID) (models.Account, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.accountRepo.GetByID(ctx, id)
}

func (s *AccountService) Create(ctx context.Context, account *models.Account) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.accountRepo.Create(ctx, account)

}

func (s *AccountService) Update(ctx context.Context, id uuid.UUID, updates models.AccountUpdate) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.accountRepo.Update(ctx, id, updates)

}

func (s *AccountService) Delete(ctx context.Context, id uuid.UUID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.accountRepo.Delete(ctx, id)

}

func (s *AccountService) Deposit(ctx context.Context, id uuid.UUID, amount float64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.accountRepo.Deposit(ctx, id, amount)
}

func (s *AccountService) Withdraw(ctx context.Context, id uuid.UUID, amount float64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.accountRepo.Withdraw(ctx, id, amount)
}
