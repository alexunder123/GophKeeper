package storage

import (
	"database/sql"
	"errors"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"

	"gophkeeper/internal/config"
	"gophkeeper/internal/crypto"
	gkerrors "gophkeeper/internal/errors"
)

// Lock структура для хранения данных о временных блокировках данных пользователя на изменение другими
type Lock struct {
	session  string
	timeLock time.Time
	update   bool
}

// Storage структура для хранения оперативных данных.
type Storage struct {
	cfg   *config.Config
	db    *sql.DB
	locks map[string]Lock
	sync.RWMutex
}

// NewStorage метод генерирует хранилище оперативных данных.
func NewStorage(cfg *config.Config) (*Storage, error) {
	db, err := sql.Open("pgx", cfg.SQLDatabase)
	if err != nil {
		return nil, err
	}
	err = createDB(db)
	if err != nil {
		return nil, err
	}
	return &Storage{
		cfg:   cfg,
		db:    db,
		locks: make(map[string]Lock),
	}, nil
}

// CheckUser метод проверят занят ли такой логин в системе
func (s *Storage) CheckUser(userLogin string) (bool, error) {
	var login string
	err := s.db.QueryRow("SELECT login FROM GophKeeper WHERE login = $1", userLogin).Scan(&login)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, err
}

// RegisterUser метод регистрирует нового пользователя
func (s *Storage) RegisterUser(userLogin, userPass string) (string, string, string, error) {
	userID := crypto.RandomID(s.cfg.LenghtUserID)
	symKey := crypto.NewSymmetricalKey(userID)
	timeStamp := time.Now().Format(time.RFC3339)
	_, err := s.db.Exec("INSERT INTO GophKeeper(user_id, login, password, aeskey, time_stamp, file) VALUES($1, $2, $3, $4, $5, $6)", userID, userLogin, userPass, symKey, timeStamp, "-")
	if err != nil {
		return "", "", "", err
	}
	return userID, symKey, timeStamp, nil
}

// AuthUser метод авторизует пользователя в системе
func (s *Storage) AuthUser(userLogin, userPass string) (string, error) { //userID, err
	var userID, login, pass string
	err := s.db.QueryRow("SELECT user_id, login, password FROM GophKeeper WHERE login = $1", userLogin).Scan(&userID, &login, &pass)
	if errors.Is(err, sql.ErrNoRows) {
		return "", gkerrors.ErrNoSuchUser
	}
	if err != nil {
		return "", err
	}
	if userPass != pass {
		return "", gkerrors.ErrWrongPassword
	}
	return userID, nil
}

// UsersData метод возвращает пользователю его сохраненные данные.
func (s *Storage) UsersData(userID string) ([]byte, string, string, error) {
	var key, timeStamp, filePath string
	err := s.db.QueryRow("SELECT aeskey, time_stamp, file FROM GophKeeper WHERE user_id = $1", userID).Scan(&key, &timeStamp, &filePath)
	if err != nil {
		return nil, "", "", err
	}
	if filePath == "-" || filePath == "" {
		return nil, timeStamp, key, gkerrors.ErrNoUserData
	}
	fileBZ, err := readFile(filePath)
	if err != nil {
		return nil, "", "", err
	}
	return fileBZ, timeStamp, key, nil
}

// UsersTimeStamp метод возвращает пользователю время последнего сохранения данных и наличие текущей блокировки на изменение данных.
func (s *Storage) UsersTimeStamp(userID string) (string, bool, string, error) {
	var timeStamp string
	err := s.db.QueryRow("SELECT time_stamp FROM GophKeeper WHERE user_id = $1", userID).Scan(&timeStamp)
	if err != nil {
		return "", false, "", err
	}
	s.RLock()
	lock, ok := s.locks[userID]
	s.RUnlock()
	if ok && lock.timeLock.After(time.Now()) {
		return timeStamp, true, lock.timeLock.Format(time.RFC3339), nil
	}
	if ok {
		s.Lock()
		delete(s.locks, userID)
		s.Unlock()
	}
	return timeStamp, false, "", nil
}

// UsersDataLock метод устанавливает временную блокировку на изменение данных, кроме текущей сессии пользователя
func (s *Storage) UsersDataLock(userID, sessionID string) (bool, string) {
	s.RLock()
	lock, ok := s.locks[userID]
	s.RUnlock()
	if ok && lock.timeLock.After(time.Now()) {
		return false, lock.timeLock.Format(time.RFC3339)
	}
	lock.timeLock = time.Now().Add(time.Minute * time.Duration(s.cfg.LockingTime))
	lock.session = sessionID
	s.Lock()
	s.locks[userID] = lock
	s.Unlock()
	return true, lock.timeLock.Format(time.RFC3339)
}

// UpdateUserData метод обновляет данные пользователя в хранилище
func (s *Storage) UpdateUserData(userID, sessionID, userTimeStamp string, userData []byte) (bool, string, error) {
	s.RLock()
	lock, ok := s.locks[userID]
	s.RUnlock()
	if (ok && lock.timeLock.After(time.Now()) && lock.session != sessionID) || lock.update {
		return false, lock.timeLock.Format(time.RFC3339), gkerrors.ErrLocked
	}
	lock.update = true
	s.Lock()
	defer s.Unlock()
	s.locks[userID] = lock
	var timeStamp, filePath string
	err := s.db.QueryRow("SELECT time_stamp, file FROM GophKeeper WHERE user_id = $1", userID).Scan(&timeStamp, &filePath)
	if err != nil {
		lock.update = false
		s.locks[userID] = lock
		return false, "", err
	}
	if timeStamp != userTimeStamp {
		lock.update = false
		s.locks[userID] = lock
		return false, "", gkerrors.ErrTimeNotEqual
	}
	if filePath == "" || filePath == "-" {
		// filePath = filepath.Join(s.cfg.DatabaseDirectory, userID+".gksf")
		filePath = s.cfg.DatabaseDirectory + userID + ".gksf"
	}
	err = writeFile(filePath, userData)
	if err != nil {
		lock.update = false
		s.locks[userID] = lock
		return false, "", err
	}
	timeStamp = time.Now().Format(time.RFC3339)
	_, err = s.db.Exec("UPDATE GophKeeper SET time_stamp=$1, file=$2 WHERE user_id=$3", timeStamp, filePath, userID)
	if err != nil {
		lock.update = false
		s.locks[userID] = lock
		return false, "", err
	}
	log.Debug().Msgf("Запись об изменениях в БД обновлена")
	delete(s.locks, userID)
	return true, timeStamp, nil
}

// createDB метод создает таблицу в БД SQL
func createDB(db *sql.DB) error {
	// _, err := db.Exec("DROP TABLE IF EXISTS GophKeeper;")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("CreateDB drop table error")
	// }
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS GophKeeper(user_id text UNIQUE, login text UNIQUE, password text, aeskey text, time_stamp text, file text);")
	if err != nil {
		return err
	}
	return nil
}

// CloseDB метод закрывает соединение с БД SQL
func (s *Storage) CloseDB() {
	err := s.db.Close()
	if err != nil {
		log.Error().Err(err).Msg("CloseDB DB closing err")
	}
	log.Info().Msg("db closed")
}

// readFile метод считывает данные из файла
func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	var fileBZ = make([]byte, fi.Size())
	_, err = file.Read(fileBZ)
	if err != nil {
		log.Error().Err(err).Msg("readFile reading file err")
		return nil, err
	}
	return fileBZ, nil
}

// writeFile метод записывает данные в файл
func writeFile(filePath string, value []byte) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(value)
	if err != nil {
		return err
	}
	return nil
}
