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
	dbDriver, dbSource := "", ""

	err := godotenv.Load("../../.env") // Since this test is run from the ./db/sqlc/main_test.go
	if err != nil {
		logger.Error("Failed to load db config from main_test.go", "error", err)
		dbDriver, dbSource = "postgres", "postgresql://root:123456@localhost:5432/gobank?sslmode=disable"
	} else {
		dbDriver, dbSource = os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE")
	}

	conn, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		logger.Error("Error creating test connection", "error", err)
		os.Exit(1)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
