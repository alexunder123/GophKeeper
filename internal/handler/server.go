package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"

	"gophkeeper/internal/config"
	"gophkeeper/internal/storage"

	pb "gophkeeper/internal/grpc/proto"

	"github.com/rs/zerolog/log"
)

// GophKeeperServer поддерживает все необходимые методы сервера.
type GophKeeperServer struct {
	pb.GophKeeperServer
	cfg  *config.Config
	strg *storage.Storage
}

// NewShortURLsServer генерирует структуру для gRPC сервера.
func NewGophKeeperServer(cfg *config.Config, strg *storage.Storage) *GophKeeperServer {
	return &GophKeeperServer{cfg: cfg, strg: strg}
}

// NewSessionID генерирует sessionID и rsa-ключ для нового подключения клиента.
func (s *GophKeeperServer) NewSessionID(ctx context.Context, in *pb.NewSessionIDRequest) (*pb.NewSessionIDResponce, error) {
	sessionID, publicKey, err := s.strg.NewSessionID()
	if err != nil {
		log.Error().Err(err).Msg("NewSessionID GenerateKeys error")
		return nil, err
	}
	var publicKeyBZ bytes.Buffer
	enc := gob.NewEncoder(&publicKeyBZ)
	err = enc.Encode(&publicKey)
	if err != nil {
		log.Error().Err(err).Msg("NewUser encoding publicKey error")
		return nil, err
	}
	publicKeyBZ.Bytes()
	return &pb.NewSessionIDResponce{SessionID: sessionID, PublicKeyBZ: publicKeyBZ.Bytes()}, nil
}

// NewUser создает нового пользователя.
func (s *GophKeeperServer) NewUser(ctx context.Context, in *pb.NewUserRequest) (*pb.NewUserResponce, error) {
	userID, symKey, err := s.strg.RegisterUser(in.SessionID, in.NewUser)
	if err != nil {
		log.Error().Err(err).Msg("NewUser RegisterUser error")
		return nil, err
	}
	var buf bytes.Buffer
	_, err = buf.Write(in.PublicKeyBZ)
	if err != nil {
		log.Error().Err(err).Msg("NewUser writing to buf error")
		return nil, err
	}
	var userPublicKey = rsa.PublicKey{}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&userPublicKey)
	if err != nil {
		log.Error().Err(err).Msg("NewUser decoding publicKey error")
		return nil, err
	}
	var responce pb.NewUserResponce
	hash := sha256.New()
	responce.UserID, err = rsa.EncryptOAEP(hash, rand.Reader, &userPublicKey, []byte(userID), []byte(`user`))
	if err != nil {
		log.Error().Err(err).Msg("NewUser EncryptOAEP userID error")
		return nil, err
	}
	responce.SymKey, err = rsa.EncryptOAEP(hash, rand.Reader, &userPublicKey, []byte(symKey), []byte(`key`))
	if err != nil {
		log.Error().Err(err).Msg("NewUser EncryptOAEP SymmetricalKey error")
		return nil, err
	}
	return &responce, nil
}

// LoginUser авторизует пользователя.
func (s *GophKeeperServer) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponce, error) {
	userID, err := s.strg.AuthUser(in.SessionID, in.LoginUser)
	if err != nil {
		log.Error().Err(err).Msg("LoginUser AuthUser error")
		return nil, err
	}
	var buf bytes.Buffer
	_, err = buf.Write(in.PublicKeyBZ)
	if err != nil {
		log.Error().Err(err).Msg("LoginUser writing to buf error")
		return nil, err
	}
	var userPublicKey = rsa.PublicKey{}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&userPublicKey)
	if err != nil {
		log.Error().Err(err).Msg("LoginUser decoding publicKey error")
		return nil, err
	}
	var responce pb.LoginUserResponce
	hash := sha256.New()
	responce.UserID, err = rsa.EncryptOAEP(hash, rand.Reader, &userPublicKey, []byte(userID), []byte(`user`))
	if err != nil {
		log.Error().Err(err).Msg("LoginUser EncryptOAEP userID error")
		return nil, err
	}
	return &responce, nil
}

// UserData передает клиенту сохраненные данные пользователя.
func (s *GophKeeperServer) UserData(ctx context.Context, in *pb.UserDataRequest) (*pb.UserDataResponce, error) {
	userData, timeStamp, symKey, userPublicKey, err := s.strg.UsersData(in.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("LoginUser AuthUser error")
		return nil, err
	}
	var responce = pb.UserDataResponce{UserData: userData, TimeStamp: timeStamp}
	hash := sha256.New()
	responce.SymKey, err = rsa.EncryptOAEP(hash, rand.Reader, &userPublicKey, []byte(symKey), []byte(`key`))
	if err != nil {
		log.Error().Err(err).Msg("NewUser EncryptOAEP SymmetricalKey error")
		return nil, err
	}
	return &responce, nil
}

// TimeStamp передает клиенту отметку времени о последних сохраненных данных пользователя.
func (s *GophKeeperServer) TimeStamp(ctx context.Context, in *pb.TimeStampRequest) (*pb.TimeStampResponce, error) {
	timeStamp, err := s.strg.UsersTimeStamp(in.SessionID)
	// if err == errors.ErrExpired {
	// 	log.Error().Err(err).Msg("TimeStamp AuthUser error")
	// 	return nil, err
	// }
	if err != nil {
		log.Error().Err(err).Msg("TimeStamp error")
		return nil, err
	}
	return &pb.TimeStampResponce{TimeStamp: timeStamp}, nil
}

// UpdateData принимает от клиента обновленные данные пользователя.
func (s *GophKeeperServer) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.UpdateDataResponce, error) {
	save, timeStamp, err := s.strg.UpdateUserData(in.SessionID, in.TimeStamp, in.UserData)
	if err != nil {
		log.Error().Err(err).Msg("UpdateData error")
		return nil, err
	}
	if !save {
		log.Error().Err(err).Msg("UpdateData error")
		return nil, err
	}
	return &pb.UpdateDataResponce{Status: save, TimeStamp: timeStamp}, nil
}

// LogOut закрывает сессию пользователя.
func (s *GophKeeperServer) LogOut(ctx context.Context, in *pb.LogOutRequest) (*pb.LogOutResponce, error) {
	status, err := s.strg.UserLogOut(in.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("LogOut error")
		return nil, err
	}
	return &pb.LogOutResponce{Status: status}, nil
}
