package db_test

import (
	"context"
	"os"
	"os/user"
	"testing"
	"time"

	"github.com/RajVerma97/golang-banking-ledger/internal/db"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInitMongo(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := mongodb.RunContainer(ctx,
		tc.WithImage("mongo:6"),
		mongodb.WithUsername("testuser"),
		mongodb.WithPassword("testpass"),
		tc.WithWaitStrategy(
			wait.ForLog("Waiting for connections").WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err, "Failed to start MongoDB container")
	t.Cleanup(func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	})

	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get connection string")

	t.Run("successful connection", func(t *testing.T) {
		t.Setenv("MONGO_URI", connStr)
		_, _, err := db.InitMongo()
		require.NoError(t, err)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		t.Setenv("MONGO_URI", "mongodb://wronguser:wrongpass@localhost:27017")
		_, _, err := db.InitMongo()
		require.Error(t, err)
	})
}

func TestMain(m *testing.M) {
	os.Clearenv()

	currentUser, err := user.Current()
	if err == nil {
		os.Setenv("HOME", currentUser.HomeDir)
	}

	os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")

	code := m.Run()
	os.Exit(code)
}
