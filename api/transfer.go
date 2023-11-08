package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	// 입력값 유효성 검사
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // 사용자 에러
		return
	}
	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	// 인자 생성
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// 생성
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 서버 에러
		return
	}

	// 성공 시
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	// 아이디 유효 확인
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err)) //ID 없을 때 404
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 데이터베이스 서버 에러
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf("[ERR] account [%d] currency mismatch : %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
