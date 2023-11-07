package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store   // db 상호작용 위해
	router *gin.Engine // 라우팅 위해
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// 계정 생성
	router.POST("/accounts", server.createAccount)
	// 계정 조회 by ID
	router.GET("/accounts/:id", server.getAccount)
	// 여러 계정 조회
	router.GET("/accounts", server.listAccount)

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
