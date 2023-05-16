package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"strconv"

	"github.com/rs/zerolog/log"
)

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

// generateKeys функция генерирует несимметричный ключ сессии пользователя
func generateKeys() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
