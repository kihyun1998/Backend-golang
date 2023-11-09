package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store    // db 상호작용 위해
	router *gin.Engine // 라우팅 위해
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// 유효성 검사 로직 등록 과정
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// 계정 생성
	router.POST("/accounts", server.createAccount)
	// 계정 조회 by ID
	router.GET("/accounts/:id", server.getAccount)
	// 여러 계정 조회
	router.GET("/accounts", server.listAccount)
	// 계정 업데이트
	router.PUT("/accounts/:id", server.updateAccount)
	// 계정 삭제
	router.DELETE("/accounts/:id", server.deleteAccount)
	// 송금
	router.POST("/transfers", server.createTransfer)
	// 유저 생성
	router.POST("/users", server.createUser)

	server.router = router
	return server
}

// 요청 서버에서 HTTP 서버 실행, API 송수신 시작
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
