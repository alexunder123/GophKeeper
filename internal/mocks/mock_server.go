package mocks

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/gob"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "gophkeeper/api/grpc/proto"
	servConf "gophkeeper/internal/server/config"
	servRsa "gophkeeper/internal/server/crypto"
)

type mockServer struct {
	pb.GophKeeperServer
	cfg *servConf.Config
	rsa *servRsa.Sessions
}

var timeStamp, realSessionID string

func Dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	servCnfg, err := servConf.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("NewConfig read environment error")
	}
	servRsa := servRsa.NewSessions(servCnfg)
	server := grpc.NewServer()

	pb.RegisterGophKeeperServer(server, &mockServer{cfg: servCnfg, rsa: servRsa})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal().Msgf("gRPC server failed: %s", err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (s *mockServer) NewSessionID(ctx context.Context, in *pb.NewSessionIDRequest) (*pb.NewSessionIDResponce, error) {
	var buf bytes.Buffer
	buf.Write(in.UserPublicKeyBZ)

	var userPublicKey = rsa.PublicKey{}
	dec := gob.NewDecoder(&buf)
	dec.Decode(&userPublicKey)
	var publicKey *rsa.PublicKey
	realSessionID, publicKey, _ = s.rsa.NewSessionID(&userPublicKey)

	var publicKeyBZ bytes.Buffer
	enc := gob.NewEncoder(&publicKeyBZ)
	enc.Encode(&publicKey)
	return &pb.NewSessionIDResponce{SessionID: realSessionID, PublicKeyBZ: publicKeyBZ.Bytes()}, nil
}

func (s *mockServer) NewUser(ctx context.Context, in *pb.NewUserRequest) (*pb.NewUserResponce, error) {
	userLogin, _, err := s.rsa.DecryptLogin(in.SessionID, in.NewUser)
	if err != nil {
		log.Error().Err(err).Msg("NewUser EncryptLogin error")
		return nil, status.Error(codes.Internal, "EncryptLogin error")
	}
	if userLogin == "123" {
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	if userLogin == "234" {
		return nil, status.Error(codes.AlreadyExists, "user with such login exists")
	}
	userID := servRsa.RandomID(s.cfg.LenghtUserID)
	symKey := servRsa.NewSymmetricalKey(userID)
	timeStamp := time.Now().Format(time.RFC3339)
	var responce = pb.NewUserResponce{TimeStamp: timeStamp}
	responce.UserID, _ = s.rsa.EncryptData(in.SessionID, userID, []byte(`userID`))
	responce.SymKey, _ = s.rsa.EncryptData(in.SessionID, symKey, []byte(`key`))

	if userLogin == "345" {
		return &responce, nil
	}

	responce.Sign, _ = s.rsa.SignData(in.SessionID)

	return &responce, nil
}

func (s *mockServer) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponce, error) {
	userLogin, _, err := s.rsa.DecryptLogin(in.SessionID, in.LoginUser)
	if err != nil {
		log.Error().Err(err).Msg("NewUser DecryptLogin error")
		return nil, status.Error(codes.Internal, "DecryptLogin error")
	}

	if userLogin == "123" {
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	userID := "1234567890"
	var responce pb.LoginUserResponce
	responce.UserID, _ = s.rsa.EncryptData(in.SessionID, userID, []byte(`userID`))

	if userLogin == "345" {
		return &responce, nil
	}

	responce.Sign, _ = s.rsa.SignData(in.SessionID)

	return &responce, nil
}

func (s *mockServer) UserData(ctx context.Context, in *pb.UserDataRequest) (*pb.UserDataResponce, error) {
	if in.SessionID == "9876543210" {
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}

	var responce pb.UserDataResponce
	symKey := servRsa.NewSymmetricalKey("1234567890")
	timeStamp = time.Now().Format(time.RFC3339)
	responce.UserData = nil
	responce.TimeStamp = timeStamp
	responce.SymKey, _ = s.rsa.EncryptData(realSessionID, symKey, []byte(`key`))

	if in.SessionID == "147852369" {
		return &responce, nil
	}

	responce.Sign, _ = s.rsa.SignData(realSessionID)
	return &responce, nil
}

func (s *mockServer) TimeStamp(ctx context.Context, in *pb.TimeStampRequest) (*pb.TimeStampResponce, error) {
	if in.SessionID == "9876543210" {
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}

	var responce = pb.TimeStampResponce{TimeStamp: time.Now().Add(time.Minute).Format(time.RFC3339), Locked: true, TimeLocked: time.Now().Add(time.Minute).Format(time.RFC3339)}

	if in.SessionID == "147852369" {
		return &responce, nil
	}

	responce.Sign, _ = s.rsa.SignData(realSessionID)

	if in.SessionID == "123654789" {
		return &responce, nil
	}

	responce.Locked = false
	responce.TimeStamp = timeStamp

	return &responce, nil
}

func (s *mockServer) DataLock(ctx context.Context, in *pb.DataLockRequest) (*pb.DataLockResponce, error) {
	if in.SessionID == "9876543210" {
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}

	locked := false
	timeLocked := time.Now().Add(time.Minute).Format(time.RFC3339)

	var responce = pb.DataLockResponce{Locked: locked, TimeLocked: timeLocked}

	if in.SessionID == "147852369" {
		return &responce, nil
	}

	responce.Sign, _ = s.rsa.SignData(realSessionID)

	if in.SessionID == "123654789" {
		return &responce, nil
	}

	responce.Locked = true

	return &responce, nil
}

func (s *mockServer) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.UpdateDataResponce, error) {
	if in.SessionID == "9876543210" {
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	if in.SessionID == "123654789" {
		return nil, status.Error(codes.PermissionDenied, "users data changes locked by another user")
	}
	if in.SessionID == "321654987" {
		return nil, status.Error(codes.FailedPrecondition, "users data timeStamp not equal to servers")
	}

	var responce = pb.UpdateDataResponce{Status: true, TimeStamp: time.Now().Add(time.Minute).Format(time.RFC3339)}

	if in.SessionID == "147852369" {
		return &responce, nil
	}

	responce.Sign, _ = s.rsa.SignData(realSessionID)

	return &responce, nil
}

func (s *mockServer) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponce, error) {
	old, _ := s.rsa.DecryptPassword(realSessionID, in.OldPassword, []byte("oldPass"))
	if old == servRsa.HashPasswd("234") {
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}

	var responce = pb.ChangePasswordResponce{Status: true}

	if old == servRsa.HashPasswd("345") {
		return &responce, nil
	}

	responce.Sign, _ = s.rsa.SignData(realSessionID)

	return &responce, nil
}

func (s *mockServer) LogOut(ctx context.Context, in *pb.LogOutRequest) (*pb.LogOutResponce, error) {
	if in.SessionID == "9876543210" {
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	responce := pb.LogOutResponce{Status: false}
	if in.SessionID == "147852369" {
		return &responce, nil
	}
	responce.Status = true
	return &responce, nil
}
