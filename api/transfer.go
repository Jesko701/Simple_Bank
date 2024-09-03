package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"solo_simple-bank_tutorial/db/sqlc"
	"solo_simple-bank_tutorial/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64   `json:"to_account_id" binding:"required,min=1"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency" binding:"required,currency"`
}

func (s *Server) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	fromAccount, valid := s.validAccount(ctx, req.FromAccountId, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(Authorization_Payload).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := fmt.Errorf("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	// check the target account transfer
	_, valid = s.validAccount(ctx, req.ToAccountId, req.Currency)
	if !valid {
		return
	}

	arg := sqlc.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
	}

	result, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// validAccount validating the current user
func (s *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (sqlc.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}

	// check the current currency with database currency
	if account.Currency != currency {
		err := fmt.Errorf("account |%v| currency mismatch: %v vs %v", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return account, false
	}

	return account, true
}
