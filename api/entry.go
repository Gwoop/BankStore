package api

import (
	db "Bankstore/db/sqlc"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type createEntryRequest struct {
	AccountID int64 `json:"accountid" binding:"required"`
	Amount    int64 `json:"amount" binding:"required"`
}

type updateEntryRequest struct {
	AmountID int64 `json:"amountid" binding:"required"`
	Amount   int64 `json:"amount" binding:"required"`
}

type createEntryResponse struct {
	AccountID int64 `json:"accountid"`
}

func (server *Server) CreateEntry(ctx *gin.Context) {
	var req createEntryRequest
	// получаем JSON и выполняем десериализацию
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateEntryParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	}
	// создаем пользователя в БД и обрабатывем ошибки(уникальность и БД)
	entry, err := server.store.CreateEntry(ctx, arg)
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

	// сформировать структуру для ответа
	rsp := createEntryResponse{
		AccountID: entry.AccountID,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type delteEntryRQ struct {
	AccountID int64 `json:"accountid"`
}

func (server *Server) DeleteEntry(ctx *gin.Context) {
	var req delteEntryRQ

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteEntry(ctx, req.AccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// сформировать структуру для ответа

	ctx.JSON(http.StatusOK, gin.H{"response": "Успешное удаление"})
}

func (server *Server) UpdateEntry(ctx *gin.Context) {
	var req updateEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateEntryParams{
		ID:     req.AmountID,
		Amount: req.Amount,
	}
	entry, err := server.store.UpdateEntry(ctx, arg)
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
