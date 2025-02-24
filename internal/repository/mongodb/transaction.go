package mongodb

import (
	"context"
	"fmt"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
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

func (r *TransactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&tx)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("transaction not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}
	return &tx, nil
}

func (r *TransactionRepository) GetByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error) {
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

func (r *TransactionRepository) Update(ctx context.Context, id string, tx *models.Transaction) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": tx}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}
