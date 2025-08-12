package db

import (
	"context"
	"database/sql"
	"gobank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createAccountMock(t *testing.T) Account {
	/*
	 * Since we need to create account for various testing action, we create a separate, independent method for
	 * create account. Since the method didn't start with Test, go test wouldn't mistake this function as unit test
	 * It's a good practice to make the unit test as loosely coupling as possible, so we separate this method as such
	 */

	// Currencies list
	currencies := []string{"EUR", "VND", "USD"}

	// Create test data
	arg := CreateAccountParams{
		Owner:    util.RandomString(7),
		Balance:  util.RandomInt(1, 10000),
		Currency: currencies[util.RandomInt(0, int64(len(currencies)-1))],
	}

	// Run the function in test
	account, err := testQueries.CreateAccount(context.Background(), arg)

	// Check if no error occur and no empty value get return
	require.NoError(t, err)
	require.NotEmpty(t, account)

	// Check if the return value match the expected output
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// Check for auto generated values
	require.NotZero(t, account.AccountID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createAccountMock(t)
}

func TestGetAccount(t *testing.T) {
	// Create account
	mock := createAccountMock(t)

	// Test get mock account by ID
	account, err := testQueries.GetAccount(context.Background(), mock.AccountID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, mock.AccountID, account.AccountID)
	require.Equal(t, mock.Owner, account.Owner)
	require.Equal(t, mock.Balance, account.Balance)
	require.Equal(t, mock.Currency, account.Currency)
	require.WithinDuration(t, mock.CreatedAt.Time, account.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	// Create mock account
	mock := createAccountMock(t)

	// Create args
	arg := UpdateAccountParams{
		AccountID: mock.AccountID,
		Balance:   util.RandomInt(1, 10000),
	}

	// Test update account
	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, mock.AccountID, account.AccountID)
	require.Equal(t, mock.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance) // Here, the expected value should be arg.Balance
	require.Equal(t, mock.Currency, account.Currency)
	require.WithinDuration(t, mock.CreatedAt.Time, account.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	// Create mock account
	mock := createAccountMock(t)

	// Test delete account
	err := testQueries.DeleteAccount(context.Background(), mock.AccountID)
	require.NoError(t, err)

	// Try getting the mock account, if fail (err == sql.ErrNoRows and account is empty), we successfully delete
	account, err := testQueries.GetAccount(context.Background(), mock.AccountID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccount(t *testing.T) {
	// Create a list of mock account
	for range 10 {
		createAccountMock(t)
	}

	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, int(arg.Limit))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
