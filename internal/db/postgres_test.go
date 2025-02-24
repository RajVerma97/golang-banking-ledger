package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/db"
	"github.com/RajVerma97/golang-banking-ledger/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInitPostgres(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		tc.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		tc.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")
	t.Cleanup(func() { pgContainer.Terminate(ctx) })

	connStr, err := pgContainer.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get connection string")

	t.Run("successful connection and migrations", func(t *testing.T) {
		t.Setenv("POSTGRES_URI", connStr)

		gormDB, err := db.InitPostgres()
		require.NoError(t, err)
		assert.NotNil(t, gormDB)

		var extExists bool
		gormDB.Raw(
			"SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp')",
		).Scan(&extExists)
		assert.True(t, extExists, "uuid-ossp extension should be installed")

		assert.True(t, gormDB.Migrator().HasTable(&models.Account{}),
			"accounts table should exist")
	})

	t.Run("missing POSTGRES_URI", func(t *testing.T) {
		os.Unsetenv("POSTGRES_URI")
		_, err := db.InitPostgres()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "POSTGRES_URI is not set")
	})
}

func TestPostgresConnectionFailure(t *testing.T) {
	t.Run("invalid connection string", func(t *testing.T) {
		t.Setenv("POSTGRES_URI", "invalid-connection-string")
		_, err := db.InitPostgres()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to PostgreSQL")
	})
}
