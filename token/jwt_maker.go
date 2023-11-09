package token

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// 대칭키 JWT maker
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (MakerForJWT, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("secretKey is too short ! Must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewJWTPayload(username, duration)
	if err != nil {
		return "", nil
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*JWTPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, ErrTokenExpired
		} else if strings.Contains(err.Error(), "token is unverifiable") {
			return nil, ErrTokenInvalid
		} else {
			log.Println(err.Error())
			return nil, err
		}
	}

	payload, ok := jwtToken.Claims.(*JWTPayload)
	if !ok {
		return nil, ErrTokenInvalid
	}
	return payload, nil
}
