package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/mongodb"
	"github.com/RajVerma97/golang-banking-ledger/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type Worker struct {
	rabbitMQChannel *amqp.Channel
	accountRepo     *postgres.AccountRepository
	transactionRepo *mongodb.TransactionRepository
}

func NewTransactionWorker(rabbitMQChannel *amqp.Channel, accountRepo *postgres.AccountRepository, transactionRepo *mongodb.TransactionRepository) *Worker {
	return &Worker{
		rabbitMQChannel: rabbitMQChannel,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (w *Worker) ProcessTransactions() {
	msgs, err := w.rabbitMQChannel.Consume(
		"transaction_queue", "", true, false, false, false, nil,
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
		return
	}

	log.Printf("Processing transaction: %s", tx.ID)
	w.handleTransaction(&tx)
}

func (w *Worker) handleTransaction(tx *models.Transaction) {
	ctx := context.Background()

	account, err := w.accountRepo.GetByID(ctx, uuid.MustParse(tx.AccountID))

	log.Println("account:")
	log.Println(account)
	if err != nil {
		log.Printf("Account not found: %s", tx.AccountID)
		tx.Status = models.FAILED
		w.updateTransaction(ctx, tx)
		return
	}

	if err := w.processTransactionLogic(tx, &account); err != nil {
		tx.Status = models.FAILED
	} else {
		tx.Status = models.SUCCESS
	}

	w.updateTransaction(ctx, tx)
}
func (w *Worker) processTransactionLogic(tx *models.Transaction, account *models.Account) error {
	log.Println("handle trnaasaaction logic")
	switch tx.Type {
	case models.DEPOSIT:
		account.Balance += tx.Amount

	case models.WITHDRAWL:
		if account.Balance < tx.Amount {
			log.Println("Insufficient funds")
			return errors.New("insufficient funds")
		}
		account.Balance -= tx.Amount

	default:
		log.Printf("Unknown transaction type: %s", tx.Type)
		return errors.New("invalid transaction type")
	}
	log.Printf("Balance after transaction: %f,accoutID:%d", account.Balance, account.ID)
	return w.accountRepo.UpdateBalance(context.Background(), account.ID, account.Balance)
}

func (w *Worker) updateTransaction(ctx context.Context, tx *models.Transaction) {
	tx.ProcessedAt = time.Now()
	if err := w.transactionRepo.Update(ctx, tx.ID, tx); err != nil {
		log.Printf("Failed to update transaction ledger: %v", err)
	}
}
