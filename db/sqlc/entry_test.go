package db

import (
	"context"
	"database/sql"
	"gobank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createEntryMock(t *testing.T) Entry {
	// First create an account since entry needs a valid account_id
	account := createAccountMock(t)

	arg := CreateEntryParams{
		AccountID: account.AccountID,
		Amount:    util.RandomInt(-1000, 1000), // Allow both positive and negative amounts
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.EntryID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createEntryMock(t)
}

func TestGetEntry(t *testing.T) {
	mock := createEntryMock(t)

	entry, err := testQueries.GetEntry(context.Background(), mock.EntryID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, mock.EntryID, entry.EntryID)
	require.Equal(t, mock.AccountID, entry.AccountID)
	require.Equal(t, mock.Amount, entry.Amount)
	require.WithinDuration(t, mock.CreatedAt.Time, entry.CreatedAt.Time, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	mock := createEntryMock(t)

	arg := UpdateEntryParams{
		EntryID: mock.EntryID,
		Amount:  util.RandomInt(-1000, 1000),
	}

	entry, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, mock.EntryID, entry.EntryID)
	require.Equal(t, mock.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.WithinDuration(t, mock.CreatedAt.Time, entry.CreatedAt.Time, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	mock := createEntryMock(t)

	err := testQueries.DeleteEntry(context.Background(), mock.EntryID)
	require.NoError(t, err)

	entry, err := testQueries.GetEntry(context.Background(), mock.EntryID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry)
}

func TestListEntry(t *testing.T) {
	for range 10 {
		createEntryMock(t)
	}

	arg := ListEntryParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, int(arg.Limit))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
