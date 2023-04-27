package storage

import (
	"crypto/rsa"
	"encoding/json"
	"os"
	"sync"
	"time"

	"gophkeeper/internal/config"
	"gophkeeper/internal/crypto"
)

// Session структура для хранения данных одной сессии.
type Session struct {
	userID     string
	privateKey *rsa.PrivateKey
	expires    time.Time
}

// Storage структура для хранения оперативных данных.
type Storage struct {
	cfg      *config.Config
	sessions map[string]Session
	sync.RWMutex
}

// NewStorage метод генерирует хранилище оперативных данных.
func NewStorage(cfg *config.Config) *Storage {
	return &Storage{
		cfg:      cfg,
		sessions: make(map[string]Session),
	}
}

func (s *Storage) NewSessionID() (string, *rsa.PublicKey, error) {
	sessionID := crypto.RandomID(s.cfg.LenghtSesionID)
	privateKey, err := crypto.GenerateKeys()
	if err != nil {
		return "", nil, err
	}
	publicKey := &privateKey.PublicKey
	s.Lock()
	s.sessions[sessionID] = Session{privateKey: privateKey, expires: time.Now().Add(time.Hour * time.Duration(s.cfg.Expires))}
	s.Unlock()
	return sessionID, publicKey, nil
}





type readerFile struct {
	file    *os.File
	decoder *json.Decoder
}

// func readStorage(cfg *config.Config, fs *Storage) {
// 	file, err := newReaderFile(cfg)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("ReadStorage NewWriterFile err")
// 	}
// 	defer file.close()
// 	file.readFile(fs)
// }

func newReaderFile(cfg *config.Config) (*readerFile, error) {
	file, err := os.OpenFile(cfg.DatabaseDirectory, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &readerFile{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// func (r *readerFile) readFile(fs *Storage) {
// 	var fileBZ = make([]byte, 0)
// 	_, err := r.file.Read(fileBZ)
// 	if err != nil {
// 		log.Error().Err(err).Msg("ReadFile reading file err")
// 		return
// 	}
// 	for r.decoder.More() {
// 		var t storageStruct
// 		err := r.decoder.Decode(&t)
// 		if err != nil {
// 			log.Error().Err(err).Msg("ReadFile decoder err")
// 			return
// 		}
// 		fs.Lock()
// 		fs.baseURL[t.Key] = t.Value
// 		fs.userURL[t.Key] = t.UserID
// 		fs.deletedURL[t.Key] = t.Deleted
// 		fs.Unlock()
// 	}
// }

func (r *readerFile) close() error {
	return r.file.Close()
}

type writerFile struct {
	file    *os.File
	encoder *json.Encoder
}

func newWriterFile(cfg *config.Config) (*writerFile, error) {
	file, err := os.OpenFile(cfg.DatabaseDirectory, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &writerFile{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// func (w *writerFile) writeFile(key, userID, value string) error {
// 	t := storageStruct{UserID: userID, Key: key, Value: value, Deleted: false}
// 	return w.encoder.Encode(&t)
// }

func (w *writerFile) close() error {
	return w.file.Close()
}
