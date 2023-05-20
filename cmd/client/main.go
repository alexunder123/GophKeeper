// Основной модуль приложения клиента
// go:build -ldflags "-X main.buildVersion=v0.0.1 -X main.buildDate=dd.mm.yyyy"
package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gophkeeper/internal/client/config"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/client/menu"
	"gophkeeper/internal/client/sender"
	"gophkeeper/internal/client/storage"
	"gophkeeper/internal/logger"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
)


func main() {
	logger.Newlogger()
	log.Info().Msg("Start client")
	cnfg, err := config.NewUserConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("NewConfig read environment error")
	}
	strg := storage.NewUserStorage()
	rsa, err := crypto.NewUserSession()
	if err != nil {
		log.Fatal().Err(err).Msg("NewUserSession generating key error")
	}
	conn, err := grpc.Dial(cnfg.RunAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("gRPC connection error")
	}
	sndr := sender.NewGophKeeperClient(conn, rsa, strg)
	fmt.Printf("Менеджер паролей GophKeeper. Версия клиента: %s, Дата сборки: %s\n", buildVersion, buildDate)
	fmt.Println("Клиент запущен, устанавливаю соединение с сервером")
	for {
		if !menu.EnteringMenu(sndr) {
			break
		}
		if !menu.AuthMenu(sndr) {
			break
		}
	}
	fmt.Println("Приложение закрывается. Нажмите клавишу Enter")
	var temp string
	fmt.Scanf("%s", &temp)
}
