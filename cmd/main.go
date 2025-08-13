package main

import (
	"database/sql"
	"gobank/api"
	db "gobank/db/sqlc"
	"gobank/util"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load config from .env
	config, err := util.LoadConfig(".")
	if err != nil {
		logger.Error("Failed to load configurations from .env", "error", err)
		return
	}

	// Connect to database
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		logger.Error("Error creating test connection", "error", err)
		return
	}

	// Create a server
	svr := api.NewServer(db.NewStore(conn), logger)
	if err := svr.Start(config.Domain, config.Port); err != nil {
		logger.Error("Error: server shutdown", "error", err)
	}
}
