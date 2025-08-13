package db

import (
	"database/sql"
	"gobank/util"
	"log/slog"
	"os"
	"testing"

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
	config, err := util.LoadConfig("../..")
	if err != nil {
		logger.Error("Failed to load db config from main_test.go", "error", err)
	}

	conn, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		logger.Error("Error creating test connection", "error", err)
		os.Exit(1)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
