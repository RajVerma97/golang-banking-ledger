package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/streadway/amqp"
)

type TransactionService struct {
	transactionRepo TransactionRepository
	accountRepo     AccountRepository
	rabbitMQChannel *amqp.Channel
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	GetByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error)
	// Update(ctx context.Context, id string, tx *models.Transaction) error
	UpdateBalance(ctx context.Context, id string, amount float64) error
}

func NewTransactionService(transactionRepo TransactionRepository, accountRepo AccountRepository, rabbitMQChannel *amqp.Channel) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		rabbitMQChannel: rabbitMQChannel,
	}
}

func (ts *TransactionService) Create(ctx context.Context, tx *models.Transaction) error {
	err := ts.transactionRepo.Create(ctx, tx)
	if err != nil {
		return err
	}

	return ts.PublishTransactionEvent(ctx, tx)
}

func (ts *TransactionService) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	return ts.transactionRepo.GetByID(ctx, id)
}

func (ts *TransactionService) GetByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error) {
	return ts.transactionRepo.GetByAccountID(ctx, accountID)
}

func (ts *TransactionService) Update(ctx context.Context, id string, amount float64) error {
	return ts.transactionRepo.Update(ctx, id, amount)
}
func (ts *TransactionService) PublishTransactionEvent(ctx context.Context, transaction *models.Transaction) error {
	body, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	err = ts.rabbitMQChannel.Publish(
		"",
		"transaction_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Println("Transaction published to RabbitMQ:", transaction.ID)
	return nil
}
