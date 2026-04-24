package api

import (
	db "Bankstore/db/sqlc"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type CreateTransferParams struct {
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	ToAccountID   int64 `json:"to_account_id" binding:"required"`
	Amount        int64 `json:"amount" binding:"required"`
}

type updateTransferRequest struct {
	ID          int64 `json:"id" binding:"required"`
	Amount      int64 `json:"amount" binding:"required"`
	ToAccountID int64 `json:"to_account_id" binding:"required"`
}

type createTransferResponse struct {
	AccountID int64 `json:"accountid"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req CreateTransferParams
	// получаем JSON и выполняем десериализацию
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	entry, err := server.store.CreateTransfer(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": entry.ID})
}

type delteTransferRQ struct {
	ID int64 `json:"id"`
}

func (server *Server) DeleteTransfer(ctx *gin.Context) {
	var req delteTransferRQ

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteTransfer(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": "Успешное удаление"})
}

func (server *Server) UpdateTransfer(ctx *gin.Context) {
	var req updateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateTransferParams{
		ID:          req.ID,
		Amount:      req.Amount,
		ToAccountID: req.ToAccountID,
	}
	entry, err := server.store.UpdateTransfer(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": entry.ID})
}
