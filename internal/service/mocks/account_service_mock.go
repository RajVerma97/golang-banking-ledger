package mocks

import (
	"context"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

var _ service.AccountServiceInterface = (*MockAccountService)(nil)

func (m *MockAccountService) GetAll(ctx context.Context) (models.Accounts, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(models.Accounts), args.Error(1)
}

func (m *MockAccountService) GetByID(ctx context.Context, id uuid.UUID) (models.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Account), args.Error(1)
}

func (m *MockAccountService) Create(ctx context.Context, account *models.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountService) Update(ctx context.Context, id uuid.UUID, updates models.AccountUpdate) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockAccountService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
