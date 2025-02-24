package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AccountNumber int       `json:"accountNumber" gorm:"unique;not null"`
	FirstName     string    `json:"firstName" gorm:"not null"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email" gorm:"unique;not null"`
	Phone         string    `json:"phone"`
	Balance       float64   `json:"balance" gorm:"not null;default:0"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}
type AccountCreate struct {
	FirstName string  `json:"firstName" validate:"required"`
	LastName  string  `json:"lastName,omitempty"`
	Email     string  `json:"email" validate:"required,email"`
	Phone     string  `json:"phone,omitempty"`
	Balance   float64 `json:"balance,omitempty"`
}
type AccountUpdate struct {
	FirstName *string  `json:"firstName,omitempty"`
	LastName  *string  `json:"lastName,omitempty"`
	Phone     *string  `json:"phone,omitempty"`
	Balance   *float64 `json:"balance"`
}

type Accounts []Account
