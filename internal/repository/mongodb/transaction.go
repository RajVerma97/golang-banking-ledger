package mongodb

import (
	"context"
	"fmt"
	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(db *mongo.Database) *TransactionRepository {
	return &TransactionRepository{
		collection: db.Collection("transactions"),
	}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	_, err := r.collection.InsertOne(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&tx)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("transaction not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}
	return &tx, nil
}

func (r *TransactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction

	cursor, err := r.collection.Find(ctx, bson.M{"accountID": accountID})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var tx models.Transaction
		if err := cursor.Decode(&tx); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return transactions, nil
}
