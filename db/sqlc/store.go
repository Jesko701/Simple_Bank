package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

// Implementing the querier interface
type Store interface {
	Querier
	TransferTx(context context.Context, arg TransferTxParams) (TransferTxResult, error)
	AccountDeleter
}

// This use to use at NewStore that return the Address of SQLStore to Store interface
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

type AccountDeleter interface {
	DeleteAccountDB(ctx context.Context, id int64) (int64, error)
}

// Delete account DB with  returning row effect
func (s *SQLStore) DeleteAccountDB(ctx context.Context, id int64) (int64, error) {
	result, err := s.db.ExecContext(ctx, "DELETE FROM accounts where id=$1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// execTx executes a function within a database transaction
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64   `json:"from_account_id"`
	ToAccountId   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		// Updating data based on TransferTxResult
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx,
				q,
				arg.FromAccountId,
				-arg.Amount,
				arg.ToAccountId,
				arg.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx,
				q,
				arg.ToAccountId,
				arg.Amount,
				arg.FromAccountId,
				-arg.Amount,
			)
		}
		return nil
	})

	if err != nil {
		return result, err
	}
	return result, nil
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 float64,
	accountID2 int64,
	amount2 float64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
