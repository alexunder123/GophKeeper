package crypto

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// UserSession структура для хранения данных одной сессии.
type UserSession struct {
	sessionID       string
	userID          string
	serverPublicKey *rsa.PublicKey
	privateKey      *rsa.PrivateKey
	symmetricalKey  string
}

// NewUserSession функция генерирует структуру хранения ключей сессии.
func NewUserSession() (*UserSession, error) {
	key, err := generateKeys()
	if err != nil {
		return nil, err
	}
	return &UserSession{privateKey: key}, nil
}

// RefreshToken метод обновляет структуру хранения ключей сессии.
func (u *UserSession) RefreshToken() error {
	var err error
	u = nil
	u, err = NewUserSession()
	if err != nil {
		return err
	}
	return nil
}

// GetPublicKey метод возвращает публичный ключ пользователя.
func (u *UserSession) GetPublicKey() rsa.PublicKey {
	return u.privateKey.PublicKey
}

// GetSessionID метод возвращает идентификатор сессии пользователя.
func (u *UserSession) GetSessionID() string {
	return u.sessionID
}

// WriteSessionID метод сохраняет идентификатор сессии пользователя и открытый ключ сервера.
func (u *UserSession) WriteSessionID(sessionID string, serverPublicKey *rsa.PublicKey) {
	u.sessionID = sessionID
	if serverPublicKey != nil { // сделал проверку для редактирования sessionID в тестах
		u.serverPublicKey = serverPublicKey
	}
}

// WriteUserID метод сохраняет идентификатор пользователя и ключ для шифрования данных.
func (u *UserSession) WriteUserID(userID, symKey string) {
	if userID != "" {
		u.userID = userID
	}
	if symKey != "" {
		u.symmetricalKey = symKey
	}
}

// EncryptData метод зашифровывает сообщение перед отправкой
func (u *UserSession) EncryptData(message string, label []byte) ([]byte, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, u.serverPublicKey, []byte(message), label)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptData метод расшифровывает данные сервера
func (u *UserSession) DecryptData(messageBZ, label []byte) (string, error) {
	hash := sha256.New()
	message, err := rsa.DecryptOAEP(hash, rand.Reader, u.privateKey, messageBZ, label)
	if err != nil {
		return "", err
	}
	return string(message), nil
}

// UserSign метод создает подпись клиента для отправки сообщений
func (u *UserSession) UserSign() ([]byte, error) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write([]byte(``))
	hashed := pssh.Sum(nil)
	sign, err := rsa.SignPSS(rand.Reader, u.privateKey, newhash, hashed, &opts)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

// CheckSign метод проверяет подпись сервера
func (u *UserSession) CheckSign(serverSignBZ []byte) error {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write([]byte(``))
	hashed := pssh.Sum(nil)
	err := rsa.VerifyPSS(u.serverPublicKey, newhash, hashed, serverSignBZ, &opts)
	if err != nil {
		return err
	}
	return nil
}

// EncryptUserData метод зашифровывает данные пользователя
func (u *UserSession) EncryptUserData(jsonBZ []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher([]byte(u.symmetricalKey[:len(u.symmetricalKey)-12]))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}
	nonce := u.symmetricalKey[len(u.symmetricalKey)-12:]
	messageBZ := aesgcm.Seal(nil, []byte(nonce), jsonBZ, nil)
	return messageBZ, nil
}

// DecryptUserData метод расшифровывает данные пользователя
func (u *UserSession) DecryptUserData(messageBZ []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher([]byte(u.symmetricalKey[:len(u.symmetricalKey)-12]))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := u.symmetricalKey[len(u.symmetricalKey)-12:]
	jsonBZ, err := aesgcm.Open(nil, []byte(nonce), messageBZ, nil)
	if err != nil {
		return nil, err
	}
	return jsonBZ, nil
}
