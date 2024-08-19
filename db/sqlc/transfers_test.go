package sqlc

import (
	"context"
	"solo_simple-bank_tutorial/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1 Account, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	ac1 := createRandomAccount(t)
	ac2 := createRandomAccount(t)
	createRandomTransfer(t, ac1, ac2)
}

func TestGetTransfer(t *testing.T) {
	ac1 := createRandomAccount(t)
	ac2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, ac1, ac2)

	result, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Validate the transfer data
	require.Equal(t, transfer1.FromAccountID, result.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, result.ToAccountID)
	require.Equal(t, transfer1.ID, result.ID)
	require.Equal(t, transfer1.Amount, result.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, result.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	ac1 := createRandomAccount(t)
	ac2 := createRandomAccount(t)

	// ac1 to ac2 = 5, ac2 to ac1 = 5.
	for i := 0; i < 5; i++ {
		createRandomTransfer(t, ac1, ac2)
		createRandomTransfer(t, ac2, ac1)
	}

	arg := ListTransfersParams{
		FromAccountID: ac2.ID,
		ToAccountID:   ac1.ID,
		Limit:         2,
		Offset:        1,
	}

	result, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, result, 2)

	for _, data := range result {
		require.NotEmpty(t, data)
		require.True(t, data.FromAccountID == ac2.ID || data.ToAccountID == ac1.ID)
		// Change this if the FromAccountID was ac1
		// require.True(t, data.FromAccountID == ac1.ID || data.ToAccountID == ac2.ID)
	}
}

/*
Decide the require.Len was based on the looping data input (line 63)
* Creating offset and limit based on data result
*/
