// Основной модуль приложения сервера
package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gophkeeper/internal/config"
	"gophkeeper/internal/crypto"
	"gophkeeper/internal/grpc/proto"
	"gophkeeper/internal/handler"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/storage"
)

func main() {
	logger.Newlogger()
	log.Info().Msg("Start program")
	cnfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("NewConfig read environment error")
	}
	strg, err := storage.NewStorage(cnfg)
	if err != nil {
		log.Fatal().Err(err).Msg("NewStorage starting DB error")
	}
	log.Debug().Msg("storage init")
	rsa := crypto.NewSessions(cnfg)
	gRPCconf := handler.NewGophKeeperServer(cnfg, strg, rsa)
	log.Debug().Msg("handler init")
	listen, err := net.Listen("tcp", cnfg.RunAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("gRPC server announce error")
	}
	s := grpc.NewServer()
	proto.RegisterGophKeeperServer(s, gRPCconf)
	reflection.Register(s)
	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatal().Msgf("gRPC server failed: %s", err)
		}
	}()
	sigChan := make(chan os.Signal, 4)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan
	log.Info().Msgf("OS cmd received stop signal")
	s.GracefulStop()
	strg.CloseDB()
}
