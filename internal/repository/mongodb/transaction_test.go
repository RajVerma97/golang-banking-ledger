package mongodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupRepo(t *testing.T) (*TransactionRepository, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(30 * time.Second),
	}

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := mongoContainer.Host(ctx)
	require.NoError(t, err)

	port, err := mongoContainer.MappedPort(ctx, "27017")
	require.NoError(t, err)

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	db := client.Database("testdb")

	err = db.Collection("transactions").Drop(ctx)
	if err != nil && err != mongo.ErrNoDocuments {
		require.NoError(t, err)
	}

	return NewTransactionRepository(db), func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}
}

func TestTransactionRepository_Create(t *testing.T) {
	repo, cleanup := setupRepo(t)
	defer cleanup()
	ctx := context.Background()

	tx := &models.Transaction{
		ID:        primitive.NewObjectID().Hex(),
		AccountID: "account123",
		Type:      "deposit",
		Amount:    100.50,
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, tx)
	require.NoError(t, err)

	var result models.Transaction
	err = repo.collection.FindOne(ctx, bson.M{"_id": tx.ID}).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, tx.AccountID, result.AccountID)
	assert.Equal(t, tx.Type, result.Type)
	assert.Equal(t, tx.Amount, result.Amount)
}

func TestTransactionRepository_GetByID(t *testing.T) {
	repo, cleanup := setupRepo(t)
	defer cleanup()
	ctx := context.Background()

	tx := &models.Transaction{
		ID:        primitive.NewObjectID().Hex(),
		AccountID: "account123",
		Amount:    200.0,
	}
	_, err := repo.collection.InsertOne(ctx, tx)
	require.NoError(t, err)

	t.Run("existing transaction", func(t *testing.T) {
		result, err := repo.GetByID(ctx, tx.ID)
		require.NoError(t, err)
		assert.Equal(t, tx.ID, result.ID)
		assert.Equal(t, tx.AccountID, result.AccountID)
	})

	t.Run("non-existing transaction", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "non-existing-id")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "transaction not found")
	})
}

func TestTransactionRepository_GetByAccountID(t *testing.T) {
	repo, cleanup := setupRepo(t)
	defer cleanup()
	ctx := context.Background()

	accountID := "account-789"
	transactions := []interface{}{
		&models.Transaction{AccountID: accountID, Amount: 100},
		&models.Transaction{AccountID: accountID, Amount: 200},
		&models.Transaction{AccountID: "other-account", Amount: 300},
	}

	_, err := repo.collection.InsertMany(ctx, transactions)
	require.NoError(t, err)

	results, err := repo.GetByAccountID(ctx, accountID)
	require.NoError(t, err)

	assert.Len(t, results, 2)
	for _, tx := range results {
		assert.Equal(t, accountID, tx.AccountID)
	}
}

func TestTransactionRepository_Update(t *testing.T) {
	repo, cleanup := setupRepo(t)
	defer cleanup()
	ctx := context.Background()

	originalTx := &models.Transaction{
		ID:        primitive.NewObjectID().Hex(),
		AccountID: "account456",
		Type:      models.WITHDRAWL,
		Amount:    500.0,
	}
	_, err := repo.collection.InsertOne(ctx, originalTx)
	require.NoError(t, err)

	update := &models.Transaction{
		ID:        originalTx.ID,
		AccountID: "account456",
		Type:      models.WITHDRAWL,
		Amount:    600.0,
		Status:    models.SUCCESS,
	}

	err = repo.Update(ctx, originalTx.ID, update)
	require.NoError(t, err)

	var updatedTx models.Transaction
	err = repo.collection.FindOne(ctx, bson.M{"_id": originalTx.ID}).Decode(&updatedTx)
	require.NoError(t, err)

	assert.Equal(t, 600.0, updatedTx.Amount)
	assert.Equal(t, models.SUCCESS, updatedTx.Status)
}

func TestTransactionRepository_Update_NonExisting(t *testing.T) {
	repo, cleanup := setupRepo(t)
	defer cleanup()
	ctx := context.Background()

	nonExistingID := primitive.NewObjectID().Hex()
	tx := &models.Transaction{
		ID:        nonExistingID,
		AccountID: "doesnt-exist",
		Amount:    100.0,
	}

	err := repo.Update(ctx, nonExistingID, tx)
	require.Error(t, err)
	assert.Equal(t, "transaction not found", err.Error())
}
