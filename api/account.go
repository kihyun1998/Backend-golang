package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	// 입력값 유효성 검사
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // 사용자 에러
		return
	}

	// 해당 키의 값을 authPayload에 저장
	// 인터페이스로 반환하기에 형변환함
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.PasetoPayload)
	// 인자 생성
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	// 계정 생성
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 서버 에러
		return
	}

	// 성공 시
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err)) //ID 없을 때 404
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 데이터베이스 서버 에러
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.PasetoPayload)
	if authPayload.Username != account.Owner {
		err := errors.New("[ERR] IT IS NOT YOUR ACCOUNT")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.PasetoPayload)

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	account, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 데이터베이스 서버 에러
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type updateAccountURIRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type updateAccountJSONRequest struct {
	Balance int64 `json:"balance" binding:"required,min=0"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var reqURI updateAccountURIRequest
	var reqJSON updateAccountJSONRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 업데이트 하려는 아이디 정보 조회
	updateAccount, err := server.store.GetAccount(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err)) //ID 없을 때 404
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 데이터베이스 서버 에러
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.PasetoPayload)
	if authPayload.Username != updateAccount.Owner {
		err := errors.New("[ERR] IT IS NOT YOUR ACCOUNT")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      reqURI.ID,
		Balance: reqJSON.Balance,
	}
	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deleteAccount, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err)) //ID 없을 때 404
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 데이터베이스 서버 에러
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.PasetoPayload)
	if authPayload.Username != deleteAccount.Owner {
		err := errors.New("[ERR] IT IS NOT YOUR ACCOUNT")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, db.Account{})
}
