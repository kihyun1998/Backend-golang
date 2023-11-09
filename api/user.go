package api

import (
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username              string    `json:"username"`
	FullName              string    `json:"full_name"`
	Email                 string    `json:"email"`
	PasswordLastChangedAt time.Time `json:"password_last_changed_at"`
	CreatedAt             time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest

	// 입력값 유효성 검사
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 인자 생성
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// 유저 생성
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 서버 에러
		return
	}

	res := CreateUserResponse{
		Username:              user.Username,
		FullName:              user.FullName,
		Email:                 user.Email,
		PasswordLastChangedAt: user.PasswordLastChangedAt,
		CreatedAt:             user.CreatedAt,
	}
	// 성공 시
	ctx.JSON(http.StatusOK, res)
}
