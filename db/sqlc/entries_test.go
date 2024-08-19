package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntries(t *testing.T, account Account) Entry {
	arg := CreateEntriesParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	}

	result, err := testQueries.CreateEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.ID)

	require.Equal(t, account.ID, result.AccountID)
	require.Equal(t, account.Balance, result.Amount)
	require.NotEmpty(t, result.CreatedAt)
	return result
}
func TestEntries(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntries(t, account)
}

func TestGetEntries(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntries(t, account)

	result, err := testQueries.GetEntries(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, entry.ID, result.ID)
	require.Equal(t, entry.Amount, result.Amount)
	require.Equal(t, entry.AccountID, result.AccountID)
	require.WithinDuration(t, entry.CreatedAt, result.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	// Create a random account
	account := createRandomAccount(t)

	// Create a number of random entries for this account
	for i := 0; i < 10; i++ {
		createRandomEntries(t, account)
	}

	// Set up parameters for listing entries
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	// Call the ListEntries function
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	// Check each returned entry
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}
