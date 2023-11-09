package token

import (
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomStr(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Issuer)
	require.WithinDuration(t, issuedAt, payload.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt.Time, time.Second)
}

// JWT 상위 버전에서는 알아서 걸러주고 Fatal 나서 테스트 진행이 안된다.
// TestExpiredJWTToken과 TestInvalidJWTToken불가
// func TestExpiredJWTToken(t *testing.T) {
// 	maker, err := NewJWTMaker(util.RandomStr(32))
// 	require.NoError(t, err)

// 	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
// 	log.Print(token)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)

// 	payload, err := maker.VerifyToken(token)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrTokenExpired.Error())
// 	require.Nil(t, payload)
// }

// 이 서명 잘못된거 테스트도 에러가 잘 나오는건 확인이 됐다.
// func TestInvalidJWTTokenAlgNone(t *testing.T) {
// 	payload, err := NewPayload(util.RandomOwner(), time.Minute)
// 	require.NoError(t, err)

// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
// 	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
// 	require.NoError(t, err)

// 	maker, err := NewJWTMaker(util.RandomStr(32))
// 	require.NoError(t, err)

// 	payload, err = maker.VerifyToken(token)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrTokenInvalid.Error())
// 	require.Nil(t, payload)
// }
