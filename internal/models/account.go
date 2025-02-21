package models

import (
	"github.com/google/uuid"
	"time"
)

type Person struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
}
type Account struct {
	ID              uuid.UUID `json:"id"`
	AccountNumber   int       `json:"accountNumber"`
	AccountHolderID uuid.UUID `json:"accountHolderID"`
	Balance         float64   `json:"balance"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

 

