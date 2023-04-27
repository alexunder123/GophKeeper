package config

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	RunAddress        string `json:"runaddress"`
	DatabaseDirectory string `json:"directory"`
	Expires           int    `json:"expires"`
	LenghtSesionID    int    `json:"lenght"`
}

func NewConfig() (*Config, error) {
	file, err := os.OpenFile("config.json", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	var fileBZ = make([]byte, 0)
	_, err = file.Read(fileBZ)
	if err != nil {
		log.Error().Err(err).Msg("ReadFile reading file err")
		return nil, err
	}
	var config Config
	if len(fileBZ) > 0 {
		err = json.Unmarshal(fileBZ, &config)
		if err != nil {
			log.Error().Err(err).Msg("ReadFile decoder err")
			return nil, err
		}
	}
	if config.RunAddress == "" {
		config.RunAddress = "127.0.0.1:3200"
	}
	if config.DatabaseDirectory == "" {
		config.DatabaseDirectory = "/users/"
	}
	if config.Expires == 0 {
		config.Expires = 24
	}
	if config.LenghtSesionID == 0 {
		config.LenghtSesionID = 16
	}
	return &config, nil
}
