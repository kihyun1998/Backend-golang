package api

import (
	"errors"
	"fmt"
	"net/http"
	"simplebank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.MakerForPaseto) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// 클라이언트에서 authorization 헤더 제공 안하는 경우
		if len(authorizationHeader) == 0 {
			err := errors.New("[ERR] AUTHORIZATION HEADER IS NOT PROVIDED")

			// 요청 중단하고 상태코드와 json 보냄
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// authorization 값 공백으로 구분해야 함
		// Bearer abcdefg... 로 키가 입력받아서 fields[0]에는 Bearer fields[1]에는 토큰값이 들어갈 것임.
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("[ERR] INVALID AUTHORIZATION HEADER FORMAT")

			// 요청 중단하고 상태코드와 json 보냄
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// 인증 타입 비교
		authorizationType := strings.ToLower(fields[0])
		// 지원하지 않는 인증 타입일 경우
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("[ERR] UNSUPPORTED AUTHORIZATION TYPE %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// 액세스 토큰 추출
		accessToken := fields[1]
		// 토큰 유효 검증
		payload, err := tokenMaker.VerifyToken(accessToken)
		// 토큰 유효하지 않으면 에러
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// context에 payload 저장
		ctx.Set(authorizationPayloadKey, payload)
		// 다음 핸들러 호출
		ctx.Next()

	}
}
