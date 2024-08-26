package api

import (
	"database/sql"
	"errors"
	"net/http"
	"solo_simple-bank_tutorial/db"
	"solo_simple-bank_tutorial/db/sqlc"
	"solo_simple-bank_tutorial/util"
	"time"

	"github.com/gin-gonic/gin"
)

type createUserParams struct {
	Username       string `json:"username" binding:"required"`
	HashedPassword string `json:"password" binding:"required,min=6"`
	FullName       string `json:"fullname" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username" `
	FullName          string    `json:"full_name" `
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user sqlc.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) CreateUser(ctx *gin.Context) {
	var req createUserParams

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashed_password, err := util.HashPassword(req.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}

	arg := sqlc.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashed_password,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	createdUser, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrUniqueViolation) {
			ctx.JSON(http.StatusForbidden, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	response := newUserResponse(createdUser)
	ctx.JSON(http.StatusCreated, response)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = util.CheckPassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, err := s.token.CreateToken(
		user.Username,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}
