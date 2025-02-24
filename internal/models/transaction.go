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
	ID          string            `json:"id" bson:"_id,omitempty" validate:"omitempty,uuid4"`
	Type        TransactionType   `json:"type" bson:"type" validate:"required,oneof=DEPOSIT WITHDRAWL"`
	Amount      float64           `json:"amount" bson:"amount" validate:"required,gt=0"`
	AccountID   string            `json:"accountID" bson:"accountID" validate:"required"`
	Status      TransactionStatus `json:"status" bson:"status" validate:"required,oneof=SUCCESS FAILED PENDING"`
	CreatedAt   time.Time         `json:"createdAt" bson:"createdAt" validate:"required"`
	UpdatedAt   time.Time         `json:"updatedAt" bson:"updatedAt" validate:"required"`
	ProcessedAt time.Time         `json:"processedAt,omitempty" bson:"processedAt,omitempty"`
}
