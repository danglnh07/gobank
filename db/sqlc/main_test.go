package db

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// MAIN entry of the package level test

// Package level variables
var (
	testQueries *Queries
	conn        *sql.DB
)

func TestMain(m *testing.M) {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load db config from .env
	err := godotenv.Load("../../.env") // Since this test is run from the ./db/sqlc/main_test.go
	if err != nil {
		logger.Error("Failed to load db config from main_test.go", "error", err)
		os.Exit(1)
	}

	conn, err = sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE"))
	if err != nil {
		logger.Error("Error creating test connection", "error", err)
		os.Exit(1)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
