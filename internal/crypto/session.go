package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	math "math/rand"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// hashPasswd функция хэширует пароль клиента для хранения на сервере
func hashPasswd(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	dst := hash.Sum(nil)
	return hex.EncodeToString(dst)
}

// RandomID функция генерирует рандомный идентификатор требуемой длины
func RandomID(n int) string {
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bts := make([]byte, n)
	math.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		bts[i] = letterBytes[math.Intn(len(letterBytes))]
	}
	return string(bts)
}

// LynnCheckOrder функция проверяет корректность введенного номера банковской карты
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

// NewSymmetricalKey функция создает симметричный ключ клиента
func NewSymmetricalKey(userID string) string {
	if len(userID) >= 28 {
		userID = userID[:16]
	}
	salt := RandomID(32 - len(userID))
	key := []byte(userID + salt)
	nonce := RandomID(12)
	return string(key) + nonce
}

// generateKeys функция генерирует несимметричный ключ сессии пользователя
func generateKeys() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
