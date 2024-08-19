package sqlc

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func floatEqual(a, b float64) bool {
	const tolerance = 1e-10
	return math.Abs(a-b) <= tolerance
}

func TestTransferTx(t *testing.T) {
	// Because the method attach to the SQLStore struct
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	n := 10
	amount := float64(10)

	final_errs := make(chan error)
	final_result := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// Running concurently
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: acc1.ID,
				ToAccountId:   acc2.ID,
				Amount:        amount,
			})

			final_errs <- err
			final_result <- result
		}()
	}

	// Check the results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-final_errs
		require.NoError(t, err)

		result := <-final_result

		// Check transfers
		transfers := result.Transfer
		require.NotEmpty(t, transfers)
		require.Equal(t, acc1.ID, transfers.FromAccountID)
		require.Equal(t, acc2.ID, transfers.ToAccountID)
		require.Equal(t, amount, transfers.Amount)
		require.NotZero(t, transfers.ID)
		require.NotZero(t, transfers.CreatedAt)
		// Check the value of transfers is exists
		_, err = store.GetTransfer(context.Background(), transfers.ID)
		require.NoError(t, err)

		// Check Entries (from account and to account)
		from_account_entries := result.FromEntry
		require.NotEmpty(t, from_account_entries)
		require.Equal(t, acc1.ID, from_account_entries.AccountID)
		require.Equal(t, -amount, from_account_entries.Amount)
		require.NotZero(t, from_account_entries.ID)
		require.NotZero(t, from_account_entries.CreatedAt)
		_, err = store.GetEntries(context.Background(), from_account_entries.ID)
		require.NoError(t, err)

		to_account_entries := result.ToEntry
		require.NotEmpty(t, to_account_entries)
		require.Equal(t, acc2.ID, to_account_entries.AccountID)
		require.Equal(t, amount, to_account_entries.Amount)
		require.NotEmpty(t, to_account_entries.ID)
		require.NotEmpty(t, to_account_entries.CreatedAt)
		_, err = store.GetEntries(context.Background(), to_account_entries.ID)
		require.NoError(t, err)

		// Check the account (Consider that the balanace has updated)
		// 1. id of the account --> 2. balance of the account
		from_account := result.FromAccount
		require.NotEmpty(t, from_account)
		require.Equal(t, acc1.ID, from_account.ID)

		to_account := result.ToAccount
		require.NotEmpty(t, to_account)
		require.Equal(t, acc2.ID, to_account.ID)

		diff1 := acc1.Balance - from_account.Balance
		diff2 := to_account.Balance - acc2.Balance
		require.True(t, floatEqual(diff1, diff2))
		require.True(t, diff1 > 0)
		require.True(t, math.Mod(diff1, amount) == 0)

		// Check the loop of transaction
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check final balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, acc1.Balance-float64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance+float64(n)*amount, updatedAccount2.Balance)
}
