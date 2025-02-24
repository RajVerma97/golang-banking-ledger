package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	repo "github.com/RajVerma97/golang-banking-ledger/internal/repository/mocks"
	"github.com/RajVerma97/golang-banking-ledger/internal/service"
	queue_mocks "github.com/RajVerma97/golang-banking-ledger/pkg/queue/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTransactionService_Create(t *testing.T) {
	validAccountID := uuid.New()
	validTx := &models.Transaction{
		AccountID: validAccountID.String(),
		Type:      models.DEPOSIT,
		Amount:    100.0,
	}

	t.Run("Invalid Account ID", func(t *testing.T) {
		mockTxRepo := new(repo.MockTransactionRepository)
		mockAccRepo := new(repo.MockAccountRepository)
		mockPublisher := new(queue_mocks.MockPublisher)
		service := service.NewTransactionService(mockTxRepo, mockAccRepo, mockPublisher)

		invalidTx := &models.Transaction{AccountID: "invalid"}
		err := service.Create(context.Background(), invalidTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid account ID")
	})

	t.Run("Account Not Found", func(t *testing.T) {
		mockTxRepo := new(repo.MockTransactionRepository)
		mockAccRepo := new(repo.MockAccountRepository)
		mockPublisher := new(queue_mocks.MockPublisher)
		service := service.NewTransactionService(mockTxRepo, mockAccRepo, mockPublisher)

		mockAccRepo.On("GetByID", mock.Anything, validAccountID).
			Return(models.Account{}, assert.AnError)

		err := service.Create(context.Background(), validTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "account verification failed")
		mockAccRepo.AssertExpectations(t)
	})

	t.Run("Transaction Creation Failure", func(t *testing.T) {
		mockTxRepo := new(repo.MockTransactionRepository)
		mockAccRepo := new(repo.MockAccountRepository)
		mockPublisher := new(queue_mocks.MockPublisher)
		service := service.NewTransactionService(mockTxRepo, mockAccRepo, mockPublisher)

		mockAccRepo.On("GetByID", mock.Anything, validAccountID).
			Return(models.Account{}, nil)
		mockTxRepo.On("Create", mock.Anything, validTx).
			Return(assert.AnError)

		err := service.Create(context.Background(), validTx)
		assert.Error(t, err)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockTxRepo := new(repo.MockTransactionRepository)
		mockAccRepo := new(repo.MockAccountRepository)
		mockPublisher := new(queue_mocks.MockPublisher)
		service := service.NewTransactionService(mockTxRepo, mockAccRepo, mockPublisher)

		mockAccRepo.On("GetByID", mock.Anything, validAccountID).
			Return(models.Account{ID: validAccountID}, nil)
		mockTxRepo.On("Create", mock.Anything, validTx).
			Return(nil)
		mockPublisher.On("Publish", "", "transaction_queue", false, false, mock.Anything).
			Return(nil)

		err := service.Create(context.Background(), validTx)
		assert.NoError(t, err)
		mockPublisher.AssertExpectations(t)
	})
}
func TestTransactionService_GetByID(t *testing.T) {
	txID := uuid.New().String()
	expectedTx := &models.Transaction{ID: txID}

	t.Run("Success", func(t *testing.T) {
		mockTxRepo := new(repo.MockTransactionRepository)
		service := service.NewTransactionService(
			mockTxRepo,
			new(repo.MockAccountRepository),
			new(queue_mocks.MockPublisher),
		)

		
		mockTxRepo.On("GetByID", mock.Anything, txID).Return(expectedTx, nil)

		tx, err := service.GetByID(context.Background(), txID)
		require.NoError(t, err)
		assert.Equal(t, expectedTx, tx)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockTxRepo := new(repo.MockTransactionRepository)
		service := service.NewTransactionService(
			mockTxRepo,
			new(repo.MockAccountRepository),
			new(queue_mocks.MockPublisher),
		)

		mockTxRepo.On("GetByID", mock.Anything, txID).
			Return((*models.Transaction)(nil), assert.AnError)

		tx, err := service.GetByID(context.Background(), txID)
		assert.Error(t, err)
		assert.Nil(t, tx)
		mockTxRepo.AssertExpectations(t)
	})
}

func TestTransactionService_GetByAccountID(t *testing.T) {
	ctx := context.Background()
	mockAccRepo := new(repo.MockAccountRepository)
	mockPublisher := new(queue_mocks.MockPublisher)

	accountID := uuid.New().String()
	expectedTxs := []models.Transaction{
		{ID: uuid.New().String(), AccountID: accountID},
	}

	t.Run("Success", func(t *testing.T) {
		mockTxRepo := new(repo.MockTransactionRepository)
		service := service.NewTransactionService(mockTxRepo, mockAccRepo, mockPublisher)
		mockTxRepo.On("GetByAccountID", ctx, accountID).Return(expectedTxs, nil)

		txs, err := service.GetByAccountID(ctx, accountID)
		require.NoError(t, err)
		assert.Equal(t, expectedTxs, txs)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockTxRepo := new(repo.MockTransactionRepository)
		service := service.NewTransactionService(mockTxRepo, mockAccRepo, mockPublisher)
		mockTxRepo.On("GetByAccountID", ctx, accountID).Return([]models.Transaction{}, assert.AnError)

		txs, err := service.GetByAccountID(ctx, accountID)
		assert.Error(t, err)
		assert.Empty(t, txs)
		mockTxRepo.AssertExpectations(t)
	})
}
func TestTransactionService_PublishTransactionEvent(t *testing.T) {
	tx := &models.Transaction{
		ID:        uuid.New().String(),
		AccountID: uuid.New().String(),
		Amount:    100.0,
		CreatedAt: time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		
		mockTxRepo := new(repo.MockTransactionRepository)
		mockAccRepo := new(repo.MockAccountRepository)
		mockPublisher := new(queue_mocks.MockPublisher)
		service := service.NewTransactionService(mockTxRepo, mockAccRepo, mockPublisher)

		mockPublisher.On("Publish", "", "transaction_queue", false, false, mock.Anything).Return(nil)

		err := service.PublishTransactionEvent(context.Background(), tx)
		assert.NoError(t, err)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("Publish Error", func(t *testing.T) {
		
		mockTxRepo := new(repo.MockTransactionRepository)
		mockAccRepo := new(repo.MockAccountRepository)
		mockPublisher := new(queue_mocks.MockPublisher)
		service := service.NewTransactionService(mockTxRepo, mockAccRepo, mockPublisher)

		mockPublisher.On("Publish", "", "transaction_queue", false, false, mock.Anything).Return(assert.AnError)

		err := service.PublishTransactionEvent(context.Background(), tx)
		assert.Error(t, err)
		mockPublisher.AssertExpectations(t)
	})
}
