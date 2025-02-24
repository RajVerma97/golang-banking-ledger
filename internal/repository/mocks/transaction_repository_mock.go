package mocks

import (
	"context"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, id string, tx *models.Transaction) error {
	args := m.Called(ctx, id, tx)
	return args.Error(0)
}
