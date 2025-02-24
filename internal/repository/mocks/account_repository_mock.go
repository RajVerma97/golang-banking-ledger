package mocks

import (
	"context"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetAll(ctx context.Context) (models.Accounts, error) {
	args := m.Called(ctx)
	return args.Get(0).(models.Accounts), args.Error(1)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Account), args.Error(1)
}

func (m *MockAccountRepository) Create(ctx context.Context, account *models.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) Update(ctx context.Context, id uuid.UUID, updates models.AccountUpdate) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
