package token

import "time"

// JWT와 PASETO 생성
type MakerForJWT interface {
	CreateToken(username string, duration time.Duration) (string, *JWTPayload, error)

	VerifyToken(token string) (*JWTPayload, error)
}

type MakerForPaseto interface {
	CreateToken(username string, duration time.Duration) (string, *PasetoPayload, error)

	VerifyToken(token string) (*PasetoPayload, error)
}
