// Модуль предназначен для долговременного хранения данных клиентов и их регистрационных данных на сервере.
// Модуль предназначен для хранения оперативных данных о блокировке данных пользователем на сервере.
package storage

import (
	"database/sql"
	"embed"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"

	gkerrors "gophkeeper/internal/errors"
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/crypto"
)

// Storager интерфейс базы данных сервера.
type Storager interface {
	CheckUser(string) (bool, error)
	RegisterUser(string, string) (string, string, string, error)
	AuthUser(string, string) (string, error)
	ChangeUserPassword(string, string, string) (bool, error)
	UsersData(string) ([]byte, string, string, error)
	UsersTimeStamp(string) (string, bool, string, error)
	UsersDataLock(string, string) (bool, string)
	UpdateUserData(string, string, string, []byte) (bool, string, error)
	CloseDB()
}

// Storage структура для хранения оперативных данных.
type Storage struct {
	cfg *config.Config
	db  *sql.DB
}

//go:embed migrate/*.sql
var embedMigrations embed.FS

// NewStorage метод генерирует хранилище оперативных данных.
func NewStorage(cfg *config.Config) (Storager, error) {
	db, err := sql.Open("pgx", cfg.SQLDatabase)
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	if err := goose.Up(db, "migrate"); err != nil {
		return nil, err
	}

	return &Storage{
		cfg: cfg,
		db:  db,
	}, nil
}

// CheckUser метод проверят занят ли такой логин в системе
func (s *Storage) CheckUser(userLogin string) (bool, error) {
	var login string
	err := s.db.QueryRow("SELECT login FROM GophKeeper WHERE login = $1", userLogin).Scan(&login)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// RegisterUser метод регистрирует нового пользователя
func (s *Storage) RegisterUser(userLogin, userPass string) (string, string, string, error) {
	userID := crypto.RandomID(s.cfg.LenghtUserID)
	symKey := crypto.NewSymmetricalKey(userID)
	timeStamp := time.Now().Format(time.RFC3339)
	_, err := s.db.Exec("INSERT INTO GophKeeper(user_id, login, password, aeskey, time_stamp) VALUES($1, $2, $3, $4, $5)", userID, userLogin, userPass, symKey, timeStamp)
	if err != nil {
		return "", "", "", err
	}
	return userID, symKey, timeStamp, nil
}

// AuthUser метод авторизует пользователя в системе
func (s *Storage) AuthUser(userLogin, userPass string) (string, error) {
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

// AuthUser метод авторизует пользователя в системе
func (s *Storage) ChangeUserPassword(userID, oldPass, newPass string) (bool, error) {
	var pass string
	err := s.db.QueryRow("SELECT password FROM GophKeeper WHERE user_id = $1", userID).Scan(&pass)
	if errors.Is(err, sql.ErrNoRows) {
		return false, gkerrors.ErrNoSuchUser
	}
	if err != nil {
		return false, err
	}
	if oldPass != pass {
		return false, gkerrors.ErrWrongPassword
	}
	_, err = s.db.Exec("UPDATE GophKeeper SET password=$1 WHERE user_id=$2", newPass, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UsersData метод возвращает пользователю его сохраненные данные.
func (s *Storage) UsersData(userID string) ([]byte, string, string, error) {
	var key, timeStamp string
	var fileBZ []byte
	err := s.db.QueryRow("SELECT aeskey, time_stamp, user_data FROM GophKeeper WHERE user_id = $1", userID).Scan(&key, &timeStamp, &fileBZ)
	if err != nil {
		return nil, "", "", err
	}
	if len(fileBZ) == 0 {
		return nil, timeStamp, key, gkerrors.ErrNoUserData
	}
	return fileBZ, timeStamp, key, nil
}

// UsersTimeStamp метод возвращает пользователю время последнего сохранения данных и наличие текущей блокировки на изменение данных.
func (s *Storage) UsersTimeStamp(userID string) (string, bool, string, error) {
	var timeStamp, timeLock string
	err := s.db.QueryRow("SELECT time_stamp FROM GophKeeper WHERE user_id = $1", userID).Scan(&timeStamp)
	if err != nil {
		return "", false, "", err
	}

	err = s.db.QueryRow("SELECT time_lock FROM GophKeeperLocks WHERE user_id = $1", userID).Scan(&timeLock)
	if err != nil {
		return timeStamp, false, "", err
	}
	lock, err := time.Parse(time.RFC3339, timeLock)
	if err != nil {
		return timeStamp, false, "", err
	}
	if lock.After(time.Now()) {
		return timeStamp, true, timeLock, nil
	}
	s.db.Exec("DELETE FROM GophKeeperLocks WHERE user_id=$1", userID)

	return timeStamp, false, "", nil
}

// UsersDataLock метод устанавливает временную блокировку на изменение данных, кроме текущей сессии пользователя
func (s *Storage) UsersDataLock(userID, sessionID string) (bool, string) {
	var timeLock string
	err := s.db.QueryRow("SELECT time_lock FROM GophKeeperLocks WHERE user_id = $1", userID).Scan(&timeLock)
	if err == nil {
		lock, err := time.Parse(time.RFC3339, timeLock)
		if err != nil {
			log.Error().Err(err).Msgf("UsersDataLock parsing timeLock error. userID = %s, timeLock = %s", userID, timeLock)
		} else if lock.After(time.Now()) {
			return true, timeLock
		}
		_, err = s.db.Exec("DELETE FROM GophKeeperLocks WHERE user_id=$1", userID)
		if err != nil {
			log.Error().Err(err).Msg("UsersDataLock deleting lock from DB error")
		}
	}

	timeLock = time.Now().Add(time.Minute * time.Duration(s.cfg.LockingTime)).Format(time.RFC3339)
	_, err = s.db.Exec("INSERT INTO GophKeeperLocks(user_id, sessionID, time_lock) VALUES($1, $2, $3)", userID, sessionID, timeLock)
	if err != nil {
		log.Error().Err(err).Msgf("UsersDataLock inserting DB error. userID = %s", userID)
		return false, ""
	}
	return true, timeLock
}

// UpdateUserData метод обновляет данные пользователя в хранилище
func (s *Storage) UpdateUserData(userID, sessionID, userTimeStamp string, userData []byte) (bool, string, error) {
	var lockedSessionID, timeLock string
	err := s.db.QueryRow("SELECT sessionID, time_lock FROM GophKeeperLocks WHERE user_id = $1", userID).Scan(&lockedSessionID, &timeLock)
	if err == nil {
		lock, err := time.Parse(time.RFC3339, lockedSessionID)
		if err != nil {
			log.Error().Err(err).Msgf("UpdateUserData parsing timeLock error. userID = %s, timeLock = %s", userID, timeLock)
		} else if lock.After(time.Now()) && sessionID != lockedSessionID {
			return true, timeLock, gkerrors.ErrLocked
		}
	}

	var timeStamp string
	err = s.db.QueryRow("SELECT time_stamp FROM GophKeeper WHERE user_id = $1", userID).Scan(&timeStamp)
	if err != nil {
		return false, "", err
	}
	if timeStamp != userTimeStamp {
		return false, "", gkerrors.ErrTimeNotEqual
	}
	timeStamp = time.Now().Format(time.RFC3339)
	_, err = s.db.Exec("UPDATE GophKeeper SET time_stamp=$1, user_data=$2 WHERE user_id=$3", timeStamp, userData, userID)
	if err != nil {
		return false, "", err
	}
	log.Debug().Msgf("Запись об изменениях в БД обновлена")
	_, err = s.db.Exec("DELETE FROM GophKeeperLocks WHERE user_id=$1", userID)
	if err != nil {
		log.Error().Err(err).Msgf("UsersDataLock deleting lock from DB error. userID = %s", userID)
	}
	return true, timeStamp, nil
}

// CloseDB метод закрывает соединение с БД SQL
func (s *Storage) CloseDB() {
	err := s.db.Close()
	if err != nil {
		log.Error().Err(err).Msg("CloseDB DB closing err")
	}
	log.Info().Msg("db closed")
}
