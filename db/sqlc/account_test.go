package sqlc

import (
	"context"
	"database/sql"
	"solo_simple-bank_tutorial/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, user.Username, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotEmpty(t, account.CreatedAt)
	require.NotEmpty(t, account.ID)

	return account
}

func createRandomAccountWithOwner(t *testing.T, owner string) Account {
	arg := CreateAccountParams{
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotEmpty(t, account.CreatedAt)
	require.NotEmpty(t, account.ID)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestAddAccountBalance(t *testing.T) {
	account := createRandomAccount(t)
	arg := AddAccountBalanceParams{
		Amount: util.RandomMoney(),
		ID:     account.ID,
	}

	result, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	expectedBalance := account.Balance + arg.Amount
	require.Equal(t, arg.ID, result.ID)
	require.Equal(t, expectedBalance, result.Balance)
	require.NotEmpty(t, result.CreatedAt)
	require.NotEmpty(t, result.Currency)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	result, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, account.ID, result.ID)
	require.Equal(t, account.Owner, result.Owner)
	require.Equal(t, account.Balance, result.Balance)

	require.NotEmpty(t, result.Currency)
	require.NotEmpty(t, result.CreatedAt)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 9; i++ {
		createRandomUser(t)
	}

	// Skip 5 accounts and get 4 accounts
	arg := ListAccountsParams{
		Limit:  4,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 4) // based on the limit

	// Ensure the account output is 4
	for _, data := range accounts {
		require.NotEmpty(t, data)
	}
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney() + account.Balance,
	}

	result, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, arg.ID, result.ID)
	require.Equal(t, arg.Balance, result.Balance)
	require.Equal(t, account.Owner, result.Owner)
	require.Equal(t, account.Currency, result.Currency)
	require.WithinDuration(t, account.CreatedAt, result.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	findAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, findAccount)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
