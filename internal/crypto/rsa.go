package crypto

import (
	"crypto/rand"
	"crypto/rsa"
)

func GenerateKeys() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
