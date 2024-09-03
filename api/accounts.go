package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"solo_simple-bank_tutorial/db/sqlc"
	"solo_simple-bank_tutorial/token"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string  `json:"currency" binding:"required,currency"`
	Balance  float64 `json:"balance" binding:"required,min=0"`
}

func (s *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(Authorization_Payload).(*token.Payload)
	arg := sqlc.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(Authorization_Payload).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belongs to user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    account,
	})
}

type listAccountRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5"`
}

func (s *Server) ListAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(Authorization_Payload).(*token.Payload)
	arg := sqlc.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  int32(req.PageSize),
		Offset: int32(req.PageID-1) * int32(req.PageSize),
	}

	account, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type updateAccountURI struct {
	ID int64 `uri:"id" binding:"required, min=1"`
}

type updateAccountJSON struct {
	Balance float64 `json:"balance" binding:"required, min=0"`
}

func (s *Server) UpdateAccount(ctx *gin.Context) {
	var reqUri updateAccountURI
	var reqBody updateAccountJSON

	// Uri checking
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	// Body checking
	if err := ctx.ShouldBindBodyWithJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := sqlc.UpdateAccountParams{
		ID:      reqUri.ID,
		Balance: reqBody.Balance,
	}

	updatedAccount, err := s.store.UpdateAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedAccount)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteAccountImpv(ctx *gin.Context, id int64, username string) (int64, error) {
	account, err := s.store.GetAccount(ctx, id)
	if err != nil {
		return 0, err
	}

	if account.Owner != username {
		return 0, fmt.Errorf("account doesn't belong to the user")
	}

	return s.store.DeleteAccountDB(ctx, id)
}

func (s *Server) DeleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, ok := ctx.MustGet(Authorization_Payload).(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid authorization payload",
		})
	}

	rows, err := s.deleteAccountImpv(ctx, req.ID, payload.Username)
	if err != nil {
		if err.Error() == "account doesn't belong to the user" {
			ctx.JSON(http.StatusForbidden, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if rows == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Account not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}
