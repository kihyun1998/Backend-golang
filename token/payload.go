package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrTokenExpired = errors.New("[ERR] TOKEN HAS EXPIRED")
var ErrTokenInvalid = errors.New("[ERR] TOKEN IS INVALID")

type Payload struct {
	jwt.RegisteredClaims
}

// 페이로드 생성
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		jwt.RegisteredClaims{
			ID:        tokenID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    username,
		},
	}
	return payload, nil
}
