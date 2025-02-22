package models

import (
	"time"
)

type TransactionType string
type TransactionStatus string

const (
	DEPOSIT   TransactionType = "DEPOSIT"
	WITHDRAWL TransactionType = "WITHDRAWL"
)
const (
	SUCCESS TransactionStatus = "SUCCESS"
	FAILED  TransactionStatus = "FAILED"
	PENDING TransactionStatus = "PENDING"
)

type Transaction struct {
	ID          string            `json:"id" bson:"_id,omitempty"`
	Type        TransactionType   `json:"type" bson:"type"`
	Amount      float64           `json:"amount" bson:"amount"`
	AccountID   string            `json:"accountID" bson:"accountID"`
	Status      TransactionStatus `json:"status" bson:"status"`
	CreatedAt   time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt" bson:"updatedAt"`
	ProcessedAt time.Time         `json:"processedAt,omitempty" bson:"processedAt,omitempty"`
}
