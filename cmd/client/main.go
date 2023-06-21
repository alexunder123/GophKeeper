// Основной модуль приложения клиента
// go:build -ldflags "-X main.buildVersion=v0.0.1 -X main.buildDate=dd.mm.yyyy"
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"gophkeeper/internal/client/config"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/client/interceptor"
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
	logfile := logger.NewUserlogger()
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

	path := filepath.Join("certificate/", "ca-cert.pem")
	pemServerCA, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not load SSL/TLS client key filepath = %s", path)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		log.Fatal().Msg("failed to add server CA's certificate")
	}
	certPath := filepath.Join("certificate/", "client-cert.pem")
	keyPath := filepath.Join("certificate/", "client-key.pem")
	clientCRT, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load SSL/TLS key")
	}
	configTLS := &tls.Config{
		Certificates: []tls.Certificate{clientCRT},
		RootCAs:      certPool,
	}

	credsTLS := credentials.NewTLS(configTLS)

	auth := interceptor.NewAuthClient(rsa)

	conn, err := grpc.Dial(cnfg.RunAddress, grpc.WithTransportCredentials(credsTLS), grpc.WithUnaryInterceptor(auth.Unary()))
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
	if logfile != nil {
		logfile.Close()
	}
	var temp string
	fmt.Scanf("%s", &temp)
}
