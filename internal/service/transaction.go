package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/pkg/queue"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TransactionService struct {
	transactionRepo   TransactionRepository
	accountRepo       AccountRepository
	rabbitMQPublisher queue.Publisher
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	GetByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error)
}

func NewTransactionService(transactionRepo TransactionRepository, accountRepo AccountRepository, rabbitMQPublisher queue.Publisher) *TransactionService {
	return &TransactionService{
		transactionRepo:   transactionRepo,
		accountRepo:       accountRepo,
		rabbitMQPublisher: rabbitMQPublisher,
	}
}


func (ts *TransactionService) Create(ctx context.Context, tx *models.Transaction) error {
	
	accountID, err := uuid.Parse(tx.AccountID)
	if err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	
	_, err = ts.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account verification failed: %w", err)
	}

	
	err = ts.transactionRepo.Create(ctx, tx)
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
func (ts *TransactionService) PublishTransactionEvent(ctx context.Context, transaction *models.Transaction) error {
	body, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	err = ts.rabbitMQPublisher.Publish(
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
