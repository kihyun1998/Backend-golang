package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// 처음 한번 실행 like 초기화
func init() {

	// 보통 현재 시간으로 Random Seed 정함
	// 그리고 UnixNano 하는 이유는 rand.Seed가 int64이기 때문
	// rand.Seed(time.Now().UnixNano())
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generated a radom integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // Int63n: min-max 사이 반환
}

// RandomStr generated a random string of length n
func RandomStr(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)] // Intn: 0-k-1
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomStr(6)
}

func RandomMoney() int64 {
	return RandomInt(1, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "KRW"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomAmount() int64 {
	return RandomInt(-1000, 1000)
}
