// Основной модуль приложения сервера
package main

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"gophkeeper/api/grpc/proto"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/crypto"
	"gophkeeper/internal/server/handler"
	"gophkeeper/internal/server/interceptor"
	"gophkeeper/internal/server/storage"
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

	path := filepath.Join("certificate/", "ca-cert.pem")
	pemClientCA, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not load SSL/TLS server key filepath = %s", path)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		log.Fatal().Msg("failed to add server CA's certificate")
	}

	certPath := filepath.Join("certificate/", "server-cert.pem")
	keyPath := filepath.Join("certificate/", "server-key.pem")
	serverCRT, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load SSL/TLS key")
	}
	configTLS := &tls.Config{
		Certificates: []tls.Certificate{serverCRT},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
	creds := credentials.NewTLS(configTLS)
	interceptor := interceptor.NewAuthInterceptor(rsa)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()), grpc.Creds(creds))
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
