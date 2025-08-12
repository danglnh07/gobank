package db

import (
	"context"
	"gobank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	// Create a store
	store := NewStore(conn)

	// Create 2 mock account for testing
	acc1 := createAccountMock(t)
	acc2 := createAccountMock(t)

	// Run the test in concurrency
	n := 10
	amount := util.RandomInt(1, 200)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for range n {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.AccountID,
				ToAccountID:   acc2.AccountID,
				Amount:        amount,
			})

			// Send the error and result to the channel
			errs <- err
			results <- result
		}()
	}

	// Check for errors and results
	for range n {
		// Check if no error occurs
		err := <-errs
		require.NoError(t, err)

		// Check if result not empty
		result := <-results
		require.NotEmpty(t, result)

		// Check if transfer is correct
		require.NotEmpty(t, result.Transfer)
		require.NotZero(t, result.Transfer.TransferID)
		require.Equal(t, acc1.AccountID, result.Transfer.FromAccountID)
		require.Equal(t, acc2.AccountID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransaction(context.Background(), result.Transfer.TransferID)
		require.NoError(t, err)

		// Check if from entry is correct
		require.NotEmpty(t, result.FromEntry)
		require.NotZero(t, result.FromEntry.EntryID)
		require.Equal(t, acc1.AccountID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), result.FromEntry.EntryID)
		require.NoError(t, err)

		// Check if to entry is correct
		require.NotEmpty(t, result.ToEntry)
		require.NotZero(t, result.ToEntry.EntryID)
		require.Equal(t, acc2.AccountID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), result.ToEntry.EntryID)
		require.NoError(t, err)

		// Check if from account is correct
		require.NotEmpty(t, result.FromAccount)
		require.NotZero(t, result.FromAccount.AccountID)
		require.Equal(t, acc1.Owner, result.FromAccount.Owner)
		require.Equal(t, acc1.Balance-amount, result.FromAccount.Balance)
		require.Equal(t, acc1.Currency, result.FromAccount.Currency)
		require.WithinDuration(t, acc1.CreatedAt.Time, result.FromAccount.CreatedAt.Time, time.Second)

		// Update the account1 balance for it to work in the next test
		acc1.Balance -= amount

		// Check if to account is correct
		require.NotEmpty(t, result.ToAccount)
		require.NotZero(t, result.ToAccount.AccountID)
		require.Equal(t, acc2.Owner, result.ToAccount.Owner)
		require.Equal(t, acc2.Balance+amount, result.ToAccount.Balance)
		require.Equal(t, acc2.Currency, result.ToAccount.Currency)
		require.WithinDuration(t, acc2.CreatedAt.Time, result.ToAccount.CreatedAt.Time, time.Second)

		// Update the account2 balance for it to work in the next test
		acc2.Balance += amount
	}
}

func TestTransferTxDeadlock(t *testing.T) {
	// Create a store
	store := NewStore(conn)

	// Create 2 mock account for testing
	acc1 := createAccountMock(t)
	acc2 := createAccountMock(t)

	// Run the test in concurrency
	n := 10
	amount := util.RandomInt(1, 200)
	errs := make(chan error)

	for i := range n {
		// Swap between acc1 and acc2 to test deadlock
		fromAccount := acc1
		toAccount := acc2

		if i%2 == 1 {
			fromAccount = acc2
			toAccount = acc1
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.AccountID,
				ToAccountID:   toAccount.AccountID,
				Amount:        amount,
			})

			// Send the error and result to the channel
			errs <- err
		}()
	}

	for range n {
		err := <-errs
		require.NoError(t, err)
	}

	// Since we perform transferring the same amount of mony with the same amount of times,
	// The final balance of 2 accounts should be the same
	res1, err := store.GetAccount(context.Background(), acc1.AccountID)
	require.NoError(t, err)
	require.NotEmpty(t, res1)

	res2, err := store.GetAccount(context.Background(), acc2.AccountID)
	require.NoError(t, err)
	require.NotEmpty(t, res2)

	require.Equal(t, acc1.Balance, res1.Balance)
	require.Equal(t, acc2.Balance, res2.Balance)
}
