package service

import (
	"context"
	"errors"
	"testing"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAccountService_GetAll(t *testing.T) {
	mockRepo := new(mocks.MockAccountRepository)
	service := NewAccountService(mockRepo)

	ctx := context.Background()
	expectedAccounts := models.Accounts{
		{ID: uuid.New(), FirstName: "John", Email: "john@example.com"},
		{ID: uuid.New(), FirstName: "Jane", Email: "jane@example.com"},
	}

	mockRepo.On("GetAll", ctx).Return(expectedAccounts, nil)

	accounts, err := service.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedAccounts, accounts)
	mockRepo.AssertExpectations(t)
}

func TestAccountService_GetByID(t *testing.T) {
	mockRepo := new(mocks.MockAccountRepository)
	service := NewAccountService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	expectedAccount := models.Account{ID: id, FirstName: "John", Email: "john@example.com"}

	mockRepo.On("GetByID", ctx, id).Return(expectedAccount, nil)

	account, err := service.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, expectedAccount, account)
	mockRepo.AssertExpectations(t)
}

func TestAccountService_Create(t *testing.T) {
	mockRepo := new(mocks.MockAccountRepository)
	service := NewAccountService(mockRepo)

	ctx := context.Background()
	account := &models.Account{
		FirstName: "John",
		Email:     "john@example.com",
	}

	mockRepo.On("Create", ctx, account).Return(nil)

	err := service.Create(ctx, account)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAccountService_Update(t *testing.T) {
	mockRepo := new(mocks.MockAccountRepository)
	service := NewAccountService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	updates := models.AccountUpdate{
		FirstName: stringPtr("Jane"),
	}

	mockRepo.On("Update", ctx, id, updates).Return(nil)

	err := service.Update(ctx, id, updates)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAccountService_Delete(t *testing.T) {
	mockRepo := new(mocks.MockAccountRepository)
	service := NewAccountService(mockRepo)

	ctx := context.Background()
	id := uuid.New()

	mockRepo.On("Delete", ctx, id).Return(nil)

	err := service.Delete(ctx, id)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAccountService_ErrorCases(t *testing.T) {
	mockRepo := new(mocks.MockAccountRepository)
	service := NewAccountService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	testErr := errors.New("test error")

	t.Run("GetAll Error", func(t *testing.T) {
		mockRepo.On("GetAll", ctx).Return(models.Accounts{}, testErr)

		_, err := service.GetAll(ctx)
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	})

	t.Run("GetByID Error", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(models.Account{}, testErr)

		_, err := service.GetByID(ctx, id)
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	})

	t.Run("Create Error", func(t *testing.T) {
		account := &models.Account{FirstName: "John", Email: "john@example.com"}
		mockRepo.On("Create", ctx, account).Return(testErr)

		err := service.Create(ctx, account)
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	})

	t.Run("Update Error", func(t *testing.T) {
		updates := models.AccountUpdate{FirstName: stringPtr("Jane")}
		mockRepo.On("Update", ctx, id, updates).Return(testErr)

		err := service.Update(ctx, id, updates)
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	})

	t.Run("Delete Error", func(t *testing.T) {
		mockRepo.On("Delete", ctx, id).Return(testErr)

		err := service.Delete(ctx, id)
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	})

	mockRepo.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
