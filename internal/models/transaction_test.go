package models

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestTransactionModel(t *testing.T) {
	validate := validator.New()
	now := time.Now()

	t.Run("Transaction Validation", func(t *testing.T) {
		tests := []struct {
			name        string
			transaction Transaction
			expectError bool
			errContains []string
		}{
			{
				name: "Valid Deposit Transaction",
				transaction: Transaction{
					Type:      DEPOSIT,
					Amount:    100.50,
					AccountID: "acc_123",
					Status:    SUCCESS,
					CreatedAt: now,
					UpdatedAt: now,
				},
				expectError: false,
			},
			{
				name: "Invalid Transaction Type",
				transaction: Transaction{
					Type:      "INVALID",
					Amount:    100.50,
					AccountID: "acc_123",
					Status:    SUCCESS,
					CreatedAt: now,
					UpdatedAt: now,
				},
				expectError: true,
				errContains: []string{"Type", "oneof"},
			},
			{
				name: "Negative Amount",
				transaction: Transaction{
					Type:      WITHDRAWL,
					Amount:    -50.00,
					AccountID: "acc_123",
					Status:    SUCCESS,
					CreatedAt: now,
					UpdatedAt: now,
				},
				expectError: true,
				errContains: []string{"Amount", "gt"},
			},
			{
				name: "Missing Account ID",
				transaction: Transaction{
					Type:      DEPOSIT,
					Amount:    100.50,
					Status:    SUCCESS,
					CreatedAt: now,
					UpdatedAt: now,
				},
				expectError: true,
				errContains: []string{"AccountID", "required"},
			},
			{
				name: "Invalid Status",
				transaction: Transaction{
					Type:      DEPOSIT,
					Amount:    100.50,
					AccountID: "acc_123",
					Status:    "INVALID",
					CreatedAt: now,
					UpdatedAt: now,
				},
				expectError: true,
				errContains: []string{"Status", "oneof"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validate.Struct(tt.transaction)
				if tt.expectError {
					assert.Error(t, err)
					if len(tt.errContains) > 0 {
						errMsg := err.Error()
						for _, fragment := range tt.errContains {
							assert.Contains(t, errMsg, fragment)
						}
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Timestamp Behavior", func(t *testing.T) {
		tx := Transaction{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.False(t, tx.CreatedAt.IsZero(), "CreatedAt should be initialized")
		assert.False(t, tx.UpdatedAt.IsZero(), "UpdatedAt should be initialized")
		assert.True(t, tx.ProcessedAt.IsZero(), "ProcessedAt should be empty by default")
	})
}
