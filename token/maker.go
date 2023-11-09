package token

import "time"

// JWT와 PASETO 생성
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}

type MakerWantPaseto interface {
	CreateToken(username string, duration time.Duration) (string, error)

	VerifyToken(token string) (*PasetoPayload, error)
}
