package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	math "math/rand"
	"time"
)

// HashPasswd функция хэширует пароль клиента для хранения на сервере
func HashPasswd(password string) string {
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

// GenerateKeys функция генерирует несимметричный ключ сессии пользователя
func GenerateKeys() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
