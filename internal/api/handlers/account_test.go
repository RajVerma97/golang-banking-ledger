package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountHandler_GetAccounts(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockAccountService)
		handler := NewAccountHandler(mockService)

		expectedAccounts := models.Accounts{
			{ID: uuid.New(), FirstName: "John"},
		}

		
		mockService.On("GetAll", mock.Anything).Return(expectedAccounts, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/accounts", nil)

		handler.GetAccounts(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Accounts
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedAccounts, response)
		mockService.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockService := new(mocks.MockAccountService)
		handler := NewAccountHandler(mockService)

		
		mockService.On("GetAll", mock.Anything).Return(models.Accounts{}, assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/accounts", nil)

		handler.GetAccounts(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestAccountHandler_GetAccountByID(t *testing.T) {
	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	t.Run("InvalidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/accounts/invalid", nil)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		handler.GetAccountByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		accountID := uuid.New()
		mockService.On("GetByID", mock.Anything, accountID).Return(models.Account{}, assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/accounts/"+accountID.String(), nil)
		c.Params = []gin.Param{{Key: "id", Value: accountID.String()}}

		handler.GetAccountByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		accountID := uuid.New()
		expectedAccount := models.Account{ID: accountID, FirstName: "John"}
		mockService.On("GetByID", mock.Anything, accountID).Return(expectedAccount, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/accounts/"+accountID.String(), nil)
		c.Params = []gin.Param{{Key: "id", Value: accountID.String()}}

		handler.GetAccountByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Account
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedAccount, response)
		mockService.AssertExpectations(t)
	})
}
func TestAccountHandler_CreateAccount(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockAccountService)
		handler := NewAccountHandler(mockService)

		newAccount := models.AccountCreate{
			FirstName: "John",
			Email:     "john@example.com",
		}
		body, _ := json.Marshal(newAccount)

		
		mockService.On("Create",
			mock.Anything, 
			mock.MatchedBy(func(acc *models.Account) bool {
				return acc.FirstName == "John" &&
					acc.Email == "john@example.com" &&
					acc.AccountNumber > 0 
			}),
		).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))

		handler.CreateAccount(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})
}
func TestAccountHandler_UpdateAccount(t *testing.T) {
	accountID := uuid.New()

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mocks.MockAccountService)
		handler := NewAccountHandler(mockService)

		updateData := models.AccountUpdate{FirstName: stringPtr("NewName")}
		body, _ := json.Marshal(updateData)

		mockService.On("Update",
			mock.Anything,
			accountID,
			updateData,
		).Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: accountID.String()}}
		c.Request = httptest.NewRequest("PUT", "/accounts/"+accountID.String(), bytes.NewBuffer(body))

		handler.UpdateAccount(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockAccountService)
		handler := NewAccountHandler(mockService)

		updateData := models.AccountUpdate{FirstName: stringPtr("NewName")}
		body, _ := json.Marshal(updateData)

		mockService.On("Update",
			mock.Anything,
			accountID,
			updateData,
		).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: accountID.String()}}
		c.Request = httptest.NewRequest("PUT", "/accounts/"+accountID.String(), bytes.NewBuffer(body))

		handler.UpdateAccount(c)
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}
func TestAccountHandler_DeleteAccount(t *testing.T) {
	accountID := uuid.New()

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mocks.MockAccountService)
		handler := NewAccountHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}
		c.Request = httptest.NewRequest("DELETE", "/accounts/invalid", nil)

		handler.DeleteAccount(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mocks.MockAccountService)
		handler := NewAccountHandler(mockService)

		mockService.On("Delete", 
			mock.Anything, 
			accountID,
		).Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: accountID.String()}}
		c.Request = httptest.NewRequest("DELETE", "/accounts/"+accountID.String(), nil)

		handler.DeleteAccount(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockAccountService)
		handler := NewAccountHandler(mockService)

		mockService.On("Delete",
			mock.Anything, 
			accountID,
		).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: accountID.String()}}
		c.Request = httptest.NewRequest("DELETE", "/accounts/"+accountID.String(), nil)

		handler.DeleteAccount(c)
		assert.Equal(t, http.StatusOK, w.Code)
		
		
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "account deleted successfully", response["message"])
		
		mockService.AssertExpectations(t)
	})
}

func stringPtr(s string) *string {
	return &s
}
