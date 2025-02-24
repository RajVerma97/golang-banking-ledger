package service

import (
	"context"
	"testing"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/mocks"
	queue_mocks "github.com/RajVerma97/golang-banking-ledger/pkg/queue/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_Create(t *testing.T) {
	accountID := uuid.New()

	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockAccountRepo := new(mocks.MockAccountRepository)
	mockPublisher := new(queue_mocks.MockPublisher)

	service := NewTransactionService(mockTransactionRepo, mockAccountRepo, mockPublisher)

	ctx := context.Background()
	tx := &models.Transaction{
		ID:        uuid.New().String(),
		Type:      models.DEPOSIT,
		Amount:    100.0,
		AccountID: accountID.String(),
		Status:    models.PENDING,
	}

	mockAccountRepo.On("GetByID", ctx, accountID).Return(models.Account{
		ID:      accountID,
		Balance: 500.0,
	}, nil)

	mockTransactionRepo.On("Create", ctx, tx).Return(nil)
	mockPublisher.On("Publish", "", "transaction_queue", false, false, mock.Anything).Return(nil)

	err := service.Create(ctx, tx)

	assert.NoError(t, err)
	mockTransactionRepo.AssertExpectations(t)
	mockAccountRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}
