package config

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

// Config хранит основные параметры конфигурации сервиса.
type Config struct {
	RunAddress        string `json:"runaddress"`      //Адрес запуска gRPC сервера
	DatabaseDirectory string `json:"directory"`       //Путь к каталогу с файлами пользователей
	SQLDatabase       string `json:"database"`        //Адрес подключения SQL-сервера
	Expires           int    `json:"expires"`         //Время жизни токена SessionID, в часах
	LenghtSesionID    int    `json:"lenghtsessionid"` //Длина токена SessionID
	LenghtUserID      int    `json:"lenghtuserid"`    //Длина идентификатора userID
	LockingTime       int    `json:"lockingtime"`     //Время блокировки на запись данных пользователем, в минутах
}

// NewConfig считывает основные параметры и генерирует структуру Config.
func NewConfig() (*Config, error) {
	file, err := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var fileBZ = make([]byte, 0)
	_, err = file.Read(fileBZ)
	if err != nil {
		log.Error().Err(err).Msg("NewConfig reading file err")
		return nil, err
	}
	var config Config
	if len(fileBZ) > 0 {
		err = json.Unmarshal(fileBZ, &config)
		if err != nil {
			log.Error().Err(err).Msg("NewConfig decoder err")
			return nil, err
		}
	}
	newConf := false
	if config.RunAddress == "" {
		config.RunAddress = "127.0.0.1:3200"
		newConf = true
	}
	if config.DatabaseDirectory == "" {
		config.DatabaseDirectory = ""
		newConf = true
	}
	if config.SQLDatabase == "" {
		config.SQLDatabase = "user=postgres password=1 host=localhost port=5432 database=postgres sslmode=disable" //postgres://postgres:1@localhost:5432/postgres?sslmode=disable
		newConf = true
	}
	if config.Expires == 0 {
		config.Expires = 2
		newConf = true
	}
	if config.LenghtSesionID == 0 {
		config.LenghtSesionID = 16
		newConf = true
	}
	if config.LenghtUserID == 0 {
		config.LenghtUserID = 12
		newConf = true
	}
	if config.LockingTime == 0 {
		config.LockingTime = 15
		newConf = true
	}

	if newConf {
		bytes, err := json.Marshal(config)
		if err != nil {
			log.Error().Err(err).Msg("NewConfig encoding to file err")
			return nil, err
		}
		_, err = file.Write(bytes)
		if err != nil {
			log.Error().Err(err).Msg("NewConfig writing to file err")
			return nil, err
		}
	}

	return &config, nil
}
