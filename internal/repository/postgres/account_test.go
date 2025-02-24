package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupRepo(t *testing.T) (*AccountRepository, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "banking_ledger",
			"POSTGRES_USER":     "admin",
			"POSTGRES_PASSWORD": "secret",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForSQL("5432/tcp", "pgx", func(host string, port nat.Port) string {
				return fmt.Sprintf("host=%s port=%s user=admin password=secret dbname=banking_ledger sslmode=disable",
					host, port.Port())
			}).WithStartupTimeout(10*time.Second),
		),
		Cmd: []string{"postgres", "-c", "listen_addresses=*"},
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	
	dsn := fmt.Sprintf("host=%s port=%s user=admin password=secret dbname=banking_ledger sslmode=disable",
		host, port.Port())

	
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        dsn,
	}), &gorm.Config{})
	require.NoError(t, err)

	
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Ping())

	
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	require.NoError(t, err, "Failed to create UUID extension")

	
	err = db.AutoMigrate(&models.Account{})
	require.NoError(t, err)

	
	t.Logf("Successfully connected to PostgreSQL at: %s:%s", host, port.Port())

	return NewAccountRepository(db), func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}
}
func TestAccountRepository_GetAll(t *testing.T) {
	t.Run("no accounts", func(t *testing.T) {
		repo, cleanup := setupRepo(t)
		defer cleanup()
		ctx := context.Background()

		accounts, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.Empty(t, accounts)
	})

	t.Run("with accounts", func(t *testing.T) {
		repo, cleanup := setupRepo(t)
		defer cleanup()
		ctx := context.Background()

		account1 := models.Account{
			ID:            uuid.New(),
			AccountNumber: 1001,
			FirstName:     "John",
			LastName:      "Doe",
			Email:         "john@example.com",
			Phone:         "1234567890",
			Balance:       100.0,
		}
		account2 := models.Account{
			ID:            uuid.New(),
			AccountNumber: 1002,
			FirstName:     "Jane",
			LastName:      "Smith",
			Email:         "jane@example.com",
			Phone:         "0987654321",
			Balance:       200.0,
		}

		require.NoError(t, repo.Create(ctx, &account1))
		require.NoError(t, repo.Create(ctx, &account2))

		accounts, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, accounts, 2)

		ids := []uuid.UUID{accounts[0].ID, accounts[1].ID}
		assert.Contains(t, ids, account1.ID)
		assert.Contains(t, ids, account2.ID)
	})
}

func TestAccountRepository_GetByID(t *testing.T) {
	t.Run("existing account", func(t *testing.T) {
		repo, cleanup := setupRepo(t)
		defer cleanup()
		ctx := context.Background()

		account := models.Account{
			ID:        uuid.New(),
			FirstName: "Alice",
			LastName:  "Brown",
			Email:     "alice@example.com",
			Phone:     "5555555555",
			Balance:   300.0,
		}
		require.NoError(t, repo.Create(ctx, &account))

		retrieved, err := repo.GetByID(ctx, account.ID)
		require.NoError(t, err)
		assert.Equal(t, account.ID, retrieved.ID)
		assert.Equal(t, account.FirstName, retrieved.FirstName)
		assert.Equal(t, account.LastName, retrieved.LastName)
		assert.Equal(t, account.Email, retrieved.Email)
		assert.Equal(t, account.Phone, retrieved.Phone)
		assert.Equal(t, account.Balance, retrieved.Balance)
	})

	t.Run("non-existing account", func(t *testing.T) {
		repo, cleanup := setupRepo(t)
		defer cleanup()
		ctx := context.Background()

		_, err := repo.GetByID(ctx, uuid.New())
		require.Error(t, err)
		assert.EqualError(t, err, "account not found")
	})
}

func TestAccountRepository_Create(t *testing.T) {
	repo, cleanup := setupRepo(t)
	defer cleanup()
	ctx := context.Background()

	account := models.Account{
		ID:        uuid.New(),
		FirstName: "Bob",
		LastName:  "Green",
		Email:     "bob@example.com",
		Phone:     "4444444444",
		Balance:   400.0,
	}

	err := repo.Create(ctx, &account)
	require.NoError(t, err)

	var retrieved models.Account
	err = repo.db.WithContext(ctx).First(&retrieved, "id = ?", account.ID).Error
	require.NoError(t, err)
	assert.Equal(t, account.ID, retrieved.ID)
}

func TestAccountRepository_Update(t *testing.T) {
	t.Run("full update", func(t *testing.T) {
		repo, cleanup := setupRepo(t)
		defer cleanup()
		ctx := context.Background()

		account := models.Account{
			ID:        uuid.New(),
			FirstName: "Original",
			LastName:  "Name",
			Email:     "original@example.com",
			Phone:     "1111111111",
			Balance:   100.0,
		}
		require.NoError(t, repo.Create(ctx, &account))

		newFirstName := "Updated"
		newLastName := "Surname"
		newPhone := "2222222222"
		newBalance := 200.0
		updates := models.AccountUpdate{
			FirstName: &newFirstName,
			LastName:  &newLastName,
			Phone:     &newPhone,
			Balance:   &newBalance,
		}

		err := repo.Update(ctx, account.ID, updates)
		require.NoError(t, err)

		updatedAccount, err := repo.GetByID(ctx, account.ID)
		require.NoError(t, err)
		assert.Equal(t, newFirstName, updatedAccount.FirstName)
		assert.Equal(t, newLastName, updatedAccount.LastName)
		assert.Equal(t, newPhone, updatedAccount.Phone)
		assert.Equal(t, newBalance, updatedAccount.Balance)
		assert.False(t, updatedAccount.UpdatedAt.IsZero())
	})

	t.Run("partial update", func(t *testing.T) {
		repo, cleanup := setupRepo(t)
		defer cleanup()
		ctx := context.Background()

		account := models.Account{
			ID:        uuid.New(),
			FirstName: "Partial",
			LastName:  "Update",
			Email:     "partial@example.com",
			Phone:     "3333333333",
			Balance:   150.0,
		}
		require.NoError(t, repo.Create(ctx, &account))

		newFirstName := "PartiallyUpdated"
		newBalance := 250.0
		updates := models.AccountUpdate{
			FirstName: &newFirstName,
			Balance:   &newBalance,
		}

		err := repo.Update(ctx, account.ID, updates)
		require.NoError(t, err)

		updatedAccount, err := repo.GetByID(ctx, account.ID)
		require.NoError(t, err)
		assert.Equal(t, newFirstName, updatedAccount.FirstName)
		assert.Equal(t, "Update", updatedAccount.LastName)  
		assert.Equal(t, "3333333333", updatedAccount.Phone) 
		assert.Equal(t, newBalance, updatedAccount.Balance)
	})

	t.Run("non-existing account", func(t *testing.T) {
		repo, cleanup := setupRepo(t)
		defer cleanup()
		ctx := context.Background()

		updates := models.AccountUpdate{FirstName: strPtr("NonExistent")}
		err := repo.Update(ctx, uuid.New(), updates)
		require.Error(t, err)
		assert.EqualError(t, err, "account not found")
	})
}

func TestAccountRepository_Delete(t *testing.T) {
	t.Run("existing account", func(t *testing.T) {
		repo, cleanup := setupRepo(t)
		defer cleanup()
		ctx := context.Background()

		account := models.Account{
			ID:        uuid.New(),
			FirstName: "Delete",
			LastName:  "Me",
			Email:     "delete@example.com",
			Phone:     "9999999999",
			Balance:   500.0,
		}
		require.NoError(t, repo.Create(ctx, &account))

		err := repo.Delete(ctx, account.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, account.ID)
		require.Error(t, err)
		assert.EqualError(t, err, "account not found")
	})

	t.Run("non-existing account", func(t *testing.T) {
		repo, cleanup := setupRepo(t)
		defer cleanup()
		ctx := context.Background()

		err := repo.Delete(ctx, uuid.New())
		require.Error(t, err)
		assert.EqualError(t, err, "account not found")
	})
}

func strPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
