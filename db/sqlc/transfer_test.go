package db

import (
	"context"
	"database/sql"
	"gobank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createTransferMock(t *testing.T) Transfer {
	// Create two accounts for transfer
	account1 := createAccountMock(t)
	account2 := createAccountMock(t)

	arg := CreateTransactionParams{
		FromAccountID: account1.AccountID,
		ToAccountID:   account2.AccountID,
		Amount:        util.RandomInt(1, 1000),
	}

	transfer, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.TransferID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createTransferMock(t)
}

func TestGetTransfer(t *testing.T) {
	mock := createTransferMock(t)

	transfer, err := testQueries.GetTransaction(context.Background(), mock.TransferID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, mock.TransferID, transfer.TransferID)
	require.Equal(t, mock.FromAccountID, transfer.FromAccountID)
	require.Equal(t, mock.ToAccountID, transfer.ToAccountID)
	require.Equal(t, mock.Amount, transfer.Amount)
	require.WithinDuration(t, mock.CreatedAt.Time, transfer.CreatedAt.Time, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	mock := createTransferMock(t)

	arg := UpdateTransactionParams{
		TransferID: mock.TransferID,
		Amount:     util.RandomInt(1, 1000),
	}

	transfer, err := testQueries.UpdateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, mock.TransferID, transfer.TransferID)
	require.Equal(t, mock.FromAccountID, transfer.FromAccountID)
	require.Equal(t, mock.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.WithinDuration(t, mock.CreatedAt.Time, transfer.CreatedAt.Time, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	mock := createTransferMock(t)

	err := testQueries.DeleteTransaction(context.Background(), mock.TransferID)
	require.NoError(t, err)

	transfer, err := testQueries.GetTransaction(context.Background(), mock.TransferID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer)
}

func TestListTransfer(t *testing.T) {
	for range 10 {
		createTransferMock(t)
	}

	arg := ListTransactionParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, int(arg.Limit))

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
