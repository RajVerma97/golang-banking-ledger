package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/mongodb"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/postgres"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	rabbitMQChannel *amqp.Channel
	accountRepo     *postgres.AccountRepository
	transactionRepo *mongodb.TransactionRepository
}

func NewTransactionWorker(rabbitMQChannel *amqp.Channel,
	accountRepo *postgres.AccountRepository,
	transactionRepo *mongodb.TransactionRepository) *Worker {
	return &Worker{
		rabbitMQChannel: rabbitMQChannel,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (w *Worker) ProcessTransactions() {
	msgs, err := w.rabbitMQChannel.Consume(
		"transaction_queue", "", false, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for msg := range msgs {
		w.processMessage(msg)
	}
}

func (w *Worker) processMessage(msg amqp.Delivery) {
	var tx models.Transaction
	if err := json.Unmarshal(msg.Body, &tx); err != nil {
		log.Printf("Error decoding transaction message: %v", err)
		msg.Nack(false, true)
		return
	}

	log.Printf("Processing transaction: %s", tx.ID)

	err := w.handleTransaction(&tx)

	if err != nil {
		log.Printf("Transaction processing failed: %v", err)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
}

func (w *Worker) handleTransaction(tx *models.Transaction) error {
	ctx := context.Background()

	existingTx, err := w.transactionRepo.GetByID(ctx, tx.ID)
	if err == nil && existingTx.Status != models.PENDING {
		log.Printf("Transaction %s already processed with status %s", tx.ID, existingTx.Status)
		return nil
	}
	accountID, err := uuid.Parse(tx.AccountID)
	if err != nil {
		log.Printf("Invalid account ID: %s, error: %v", tx.AccountID, err)
		tx.Status = models.FAILED
		w.updateTransaction(ctx, tx)
		return fmt.Errorf("invalid account ID: %w", err)
	}

	account, err := w.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		log.Printf("Account not found: %s", accountID)
		tx.Status = models.FAILED
		w.updateTransaction(ctx, tx)
		return errors.New("account not found")
	}

	if err := w.processTransactionLogic(tx, &account); err != nil {
		tx.Status = models.FAILED
		log.Printf("Error processing transaction: %v", err)
		w.updateTransaction(ctx, tx)
		return err
	}

	tx.Status = models.SUCCESS
	w.updateTransaction(ctx, tx)

	return nil
}
func (w *Worker) processTransactionLogic(tx *models.Transaction, account *models.Account) error {

	switch tx.Type {
	case models.DEPOSIT:
		account.Balance += tx.Amount
	case models.WITHDRAWL:
		if account.Balance < tx.Amount {
			log.Printf("Insufficient funds: current balance %f, withdrawal amount %f", account.Balance, tx.Amount)
			return errors.New("insufficient funds")
		}
		account.Balance -= tx.Amount
	default:
		log.Printf("Unknown transaction type: %s", tx.Type)
		return errors.New("invalid transaction type")
	}

	updates := models.AccountUpdate{
		Balance: &account.Balance,
	}

	if err := w.accountRepo.Update(context.Background(), account.ID, updates); err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	log.Printf("Successfully updated balance to %f for account %s", account.Balance, account.ID)
	return nil
}

func (w *Worker) updateTransaction(ctx context.Context, tx *models.Transaction) {
	tx.ProcessedAt = time.Now()
	if err := w.transactionRepo.Update(ctx, tx.ID, tx); err != nil {
		log.Printf("Failed to update transaction ledger: %v", err)
	}
}
