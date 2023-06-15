// Модуль опередляет и конфигурирует внешний пакет для логирования работы сервиса.
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Newlogger конфигурирует формат и уровень отображения логов сервера.
func Newlogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	}
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// Newlogger конфигурирует формат и уровень отображения логов клиента.
func NewUserlogger() *os.File{
	zerolog.TimeFieldFormat = time.RFC3339

	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	}
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	file, err := os.OpenFile("logs.json", os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		log.Error().Err(err).Msg("Opening log file error. Writing log to console")
		return nil
	}
	log.Logger = zerolog.New(file).With().Timestamp().Logger()
	return file
}
