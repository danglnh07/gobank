package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Store provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// Constructor method for Store struct
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// Method to execute a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Create transaction object
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Create queries object and run transaction
	q := New(tx)
	err = fn(q)

	// If queries fail, then we try to rollback
	if err != nil {
		if rollbacnErr := tx.Rollback(); rollbacnErr != nil {
			return fmt.Errorf("transaction error: %v\nrollback error: %v", err, rollbacnErr)
		}
		return err
	}

	// If success, we commit the transaction
	return tx.Commit()
}

// Parameter struct for transfer money action
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// Result struct return after transferring money
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Method to perform transfer money action
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Execute database transaction
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create a transfer record in database
		result.Transfer, err = q.CreateTransaction(ctx, CreateTransactionParams(arg))
		if err != nil {
			return err
		}

		// Add an entry for the from account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // Since the money go out, it should be minus
		})
		if err != nil {
			return err
		}

		// Add an entry for the to account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update account balance
		// Here, we always keep the order of operation fix (always update the account with lower ID first)
		// to prevent deadlock (prevent circular wait by establish a total ordering)
		fromArg := AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		}
		toArg := AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		}

		if arg.FromAccountID < arg.ToAccountID {
			updateFromAccountBalance(q, &result, fromArg)
			updateToAccountBalance(q, &result, toArg)
		} else {
			updateToAccountBalance(q, &result, toArg)
			updateFromAccountBalance(q, &result, fromArg)
		}

		return nil
	})

	return result, err
}

// Helper method: update the from_account balance
func updateFromAccountBalance(q *Queries, result *TransferTxResult, arg AddAccountBalanceParams) error {
	var err error
	result.FromAccount, err = q.AddAccountBalance(context.Background(), arg)
	return err
}

// Helper method: update the to_account balance
func updateToAccountBalance(q *Queries, result *TransferTxResult, arg AddAccountBalanceParams) error {
	var err error
	result.ToAccount, err = q.AddAccountBalance(context.Background(), arg)
	return err
}
