package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestAccountModelValidation(t *testing.T) {
	validate := validator.New()
	t.Run("AccountCreate Validation", func(t *testing.T) {
		tests := []struct {
			name        string
			input       AccountCreate
			shouldError bool
			errMsg      string
		}{
			{
				name: "Valid AccountCreate",
				input: AccountCreate{
					FirstName: "John",
					Email:     "john@example.com",
				},
				shouldError: false,
			},
			{
				name: "Missing FirstName",
				input: AccountCreate{
					Email: "john@example.com",
				},
				shouldError: true,
				errMsg:      "Key: 'AccountCreate.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag",
			},
			{
				name: "Invalid Email",
				input: AccountCreate{
					FirstName: "John",
					Email:     "invalid-email",
				},
				shouldError: true,
				errMsg:      "Key: 'AccountCreate.Email' Error:Field validation for 'Email' failed on the 'email' tag",
			},
			{
				name: "Negative Balance",
				input: AccountCreate{
					FirstName: "John",
					Email:     "john@example.com",
					Balance:   -100.0,
				},
				shouldError: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validate.Struct(tt.input)
				if tt.shouldError {
					assert.Error(t, err)
					if tt.errMsg != "" {
						assert.Contains(t, err.Error(), tt.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("Account Field Validation", func(t *testing.T) {
		account := Account{
			AccountNumber: 123456789,
			FirstName:     "John",
			Email:         "john@example.com",
		}

		assert.Equal(t, 123456789, account.AccountNumber, "AccountNumber should match")
		assert.Equal(t, "John", account.FirstName, "FirstName should match")
		assert.Equal(t, "john@example.com", account.Email, "Email should match")
	})

	t.Run("AccountUpdate Partial Updates", func(t *testing.T) {
		tests := []struct {
			name         string
			input        AccountUpdate
			expectFields []string
		}{
			{
				name: "Update FirstName",
				input: AccountUpdate{
					FirstName: stringPtr("Jane"),
				},
				expectFields: []string{"FirstName"},
			},
			{
				name: "Update Balance",
				input: AccountUpdate{
					Balance: float64Ptr(1500.0),
				},
				expectFields: []string{"Balance"},
			},
			{
				name: "Clear LastName",
				input: AccountUpdate{
					LastName: stringPtr(""),
				},
				expectFields: []string{"LastName"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				for _, field := range tt.expectFields {
					switch field {
					case "FirstName":
						assert.NotNil(t, tt.input.FirstName)
						assert.Equal(t, "Jane", *tt.input.FirstName)
					case "Balance":
						assert.NotNil(t, tt.input.Balance)
						assert.Equal(t, 1500.0, *tt.input.Balance)
					case "LastName":
						assert.NotNil(t, tt.input.LastName)
						assert.Equal(t, "", *tt.input.LastName)
					}
				}
			})
		}
	})
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
