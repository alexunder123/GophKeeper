package config

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

// UserConfig хранит основные параметры конфигурации клиента.
type UserConfig struct {
	RunAddress string `json:"runaddress"` //Адрес запуска gRPC сервера
}

// NewConfig считывает основные параметры и генерирует структуру Config.
func NewUserConfig() (*UserConfig, error) {
	file, err := os.OpenFile("user_config.json", os.O_RDWR|os.O_CREATE, 0777)
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
		log.Error().Err(err).Msg("NewConfig reading file err")
		return nil, err
	}
	var userConfig UserConfig
	if len(fileBZ) > 0 {
		err = json.Unmarshal(fileBZ, &userConfig)
		if err != nil {
			log.Error().Err(err).Msg("NewConfig decoder err")
			return nil, err
		}
	}
	var newConf bool
	if userConfig.RunAddress == "" {
		userConfig.RunAddress = "127.0.0.1:3200"
		newConf = true
	}

	if newConf {
		bytes, err := json.Marshal(userConfig)
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

	return &userConfig, nil
}
