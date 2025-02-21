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
	ID          uuid.UUID         `json:"id"`
	Type        TransactionType   `json:"type"`
	Amount      float64           `json:"amount"`
	AccountID   uuid.UUID         `json:"accountID"`
	Status      TransactionStatus `json:"status"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	ProcessedAt time.Time         `json:"proccessedAt"`
}
