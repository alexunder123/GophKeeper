package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func HashPasswd(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	dst := hash.Sum(nil)
	return hex.EncodeToString(dst)
}

func RandomID(n int) string {
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bts := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		bts[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bts)
}

func LynnCheckOrder(lynn []byte) bool {
	lynnArr := make([]int, len(lynn))
	for i, d := range lynn {
		j, err := strconv.Atoi(string(d))
		if err != nil {
			log.Error().Err(err).Msg("LynnCheckOrder strconv err")
			return false
		}
		lynnArr[i] = j
	}
	for i := len(lynnArr) - 2; i >= 0; i -= 2 {
		n := lynnArr[i] * 2
		if n >= 10 {
			n -= 9
		}
		lynnArr[i] = n
	}
	sum := 0
	for i := 0; i < len(lynnArr); i++ {
		sum += lynnArr[i]
	}
	return sum%10 == 0
}
