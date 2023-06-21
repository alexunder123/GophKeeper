package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	gkerrors "gophkeeper/internal/errors"
	"gophkeeper/internal/server/config"
	"strings"
	"sync"
	"time"
)

// Session структура для хранения данных одной сессии.
type Session struct {
	userID        string
	userPublicKey *rsa.PublicKey
	privateKey    *rsa.PrivateKey
	expires       time.Time
}

// Sessions структура для хранения оперативных данных.
type Sessions struct {
	cfg      *config.Config
	sessions map[string]Session
	sync.RWMutex
}

// NewSessions функция генерирует структуру хранения оперативных данных клиентов.
func NewSessions(cfg *config.Config) *Sessions {
	s := Sessions{
		cfg:      cfg,
		sessions: make(map[string]Session),
	}
	s.sessionsCleaner(cfg.Expires)
	return &s
}

func (s *Sessions) sessionsCleaner(period int) {
	go func() {
		ticker := time.NewTicker(time.Hour * time.Duration(period/2))
		defer ticker.Stop()
		for {
			<-ticker.C
			s.Lock()
			for i, v := range s.sessions {
				if !v.expires.After(time.Now()) {
					delete(s.sessions, i)
				}
			}
			s.Unlock()
		}
	}()
}

// NewSessionID метод генерирует асимметричный ключ и сохраняет новую сессию
func (s *Sessions) NewSessionID(userKey *rsa.PublicKey) (string, *rsa.PublicKey, error) {
	sessionID := RandomID(s.cfg.LenghtSesionID)
	privateKey, err := GenerateKeys()
	if err != nil {
		return "", nil, err
	}
	publicKey := &privateKey.PublicKey
	s.Lock()
	s.sessions[sessionID] = Session{userPublicKey: userKey, privateKey: privateKey, expires: time.Now().Add(time.Hour * time.Duration(s.cfg.Expires))}
	s.Unlock()
	return sessionID, publicKey, nil
}

// CheckSign метод проверяет подпись и валидность сессии клиента
func (s *Sessions) CheckSign(sessionID string, userSignBZ []byte) (string, error) {
	s.RLock()
	session, ok := s.sessions[sessionID]
	s.RUnlock()
	if !ok {
		return "", gkerrors.ErrExpired
	}
	if !session.expires.After(time.Now()) {
		s.Lock()
		delete(s.sessions, sessionID)
		s.Unlock()
		return "", gkerrors.ErrExpired
	}
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write([]byte(``))
	hashed := pssh.Sum(nil)
	err := rsa.VerifyPSS(session.userPublicKey, newhash, hashed, userSignBZ, &opts)
	if err != nil {
		return "", gkerrors.ErrSignIncorrect
	}
	s.RLock()
	userID := s.sessions[sessionID].userID
	s.RUnlock()
	return userID, nil
}

func (s *Sessions) GetUserID(sessionID string) string {
	s.RLock()
	userID := s.sessions[sessionID].userID
	s.RUnlock()
	return userID
}

// DecryptLogin метод расшифровывает логин и пароль пользователя
func (s *Sessions) DecryptLogin(sessionID string, userLoginBZ []byte) (string, string, error) {
	s.RLock()
	privateKey := s.sessions[sessionID].privateKey
	s.RUnlock()
	hash := sha256.New()
	text, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, userLoginBZ, []byte(`login`))
	if err != nil {
		return "", "", err
	}
	login, pass, found := strings.Cut(string(text), ",")
	if !found {
		return "", "", gkerrors.ErrLoginIncorrect
	}
	return login, HashPasswd(pass), nil
}

// EncryptData метод зашифровывает сообщение перед отправкой
func (s *Sessions) EncryptData(sessionID, message string, label []byte) ([]byte, error) {
	s.RLock()
	userPublicKey := s.sessions[sessionID].userPublicKey
	s.RUnlock()
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, userPublicKey, []byte(message), label)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptPassword метод расшифровывает полученное сообщение
func (s *Sessions) DecryptPassword(sessionID string, messageBZ, label []byte) (string, error) {
	s.RLock()
	privateKey := s.sessions[sessionID].privateKey
	s.RUnlock()
	hash := sha256.New()
	message, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, messageBZ, label)
	if err != nil {
		return "", err
	}
	return HashPasswd(string(message)), nil
}

// SignData метод создает подпись сервера для отправки сообщений
func (s *Sessions) SignData(sessionID string) ([]byte, error) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write([]byte(``))
	hashed := pssh.Sum(nil)
	s.RLock()
	privateKey := s.sessions[sessionID].privateKey
	s.RUnlock()
	sign, err := rsa.SignPSS(rand.Reader, privateKey, newhash, hashed, &opts)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

// AddUserID метод добавляет userID в сессию клиента после его аутентификации
func (s *Sessions) AddUserID(sessionID, userID string) {
	s.Lock()
	user := s.sessions[sessionID]
	user.userID = userID
	s.sessions[sessionID] = user
	s.Unlock()
}

// UserLogOut метод удаляет сессию клиента
func (s *Sessions) UserLogOut(sessionID string) {
	s.Lock()
	delete(s.sessions, sessionID)
	s.Unlock()
}
