package models

import (
	"github.com/google/uuid"
	"time"
)

type TransactionType uint32
type TransactionStatus uint32

const (
	DEBIT TransactionType = iota
	CREDIT
)
const (
	SUCCESS TransactionStatus = iota
	FAILED
	PENDING
)

type Transaction struct {
	ID          uuid.UUID         `json:"id" bson:"_id,omitempty"`
	Type        TransactionType   `json:"type" bson:"type"`
	Amount      float64           `json:"amount" bson:"amount"`
	AccountID   uuid.UUID         `json:"accountID" bson:"accountID"`
	Status      TransactionStatus `json:"status" bson:"status"`
	CreatedAt   time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt" bson:"updatedAt"`
	ProcessedAt time.Time         `json:"processedAt,omitempty" bson:"processedAt,omitempty"`
}
