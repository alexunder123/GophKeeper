// Модуль предназначен для приема сервером сообщений от пользователей, обработки их и отправки ответов пользователям.
package handler

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/gob"
	"errors"

	"gophkeeper/internal/config"
	"gophkeeper/internal/crypto"
	gkerrors "gophkeeper/internal/errors"
	"gophkeeper/internal/storage"

	pb "gophkeeper/internal/grpc/proto"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GophKeeperServer поддерживает все необходимые методы сервера.
type GophKeeperServer struct {
	pb.GophKeeperServer
	cfg  *config.Config
	strg *storage.Storage
	rsa  *crypto.Sessions
}

// NewGophKeeperServer генерирует структуру для gRPC сервера.
func NewGophKeeperServer(cfg *config.Config, strg *storage.Storage, rsa *crypto.Sessions) *GophKeeperServer {
	return &GophKeeperServer{cfg: cfg, strg: strg, rsa: rsa}
}

// NewSessionID генерирует sessionID и rsa-ключ для нового подключения клиента.
func (s *GophKeeperServer) NewSessionID(ctx context.Context, in *pb.NewSessionIDRequest) (*pb.NewSessionIDResponce, error) {
	var buf bytes.Buffer
	_, err := buf.Write(in.UserPublicKeyBZ)
	if err != nil {
		log.Error().Err(err).Msg("NewSessionID writing to buf error")
		return nil, status.Error(codes.Internal, "writing to buf error")
	}

	var userPublicKey = rsa.PublicKey{}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&userPublicKey)
	if err != nil {
		log.Error().Err(err).Msg("NewSessionID decoding publicKey error")
		return nil, status.Error(codes.Internal, "decoding publicKey erro")
	}

	sessionID, publicKey, err := s.rsa.NewSessionID(&userPublicKey)
	if err != nil {
		log.Error().Err(err).Msg("NewSessionID GenerateKeys error")
		return nil, status.Error(codes.Internal, "GenerateKeys error")
	}

	var publicKeyBZ bytes.Buffer
	enc := gob.NewEncoder(&publicKeyBZ)
	err = enc.Encode(&publicKey)
	if err != nil {
		log.Error().Err(err).Msg("NewUser encoding publicKey error")
		return nil, status.Error(codes.Internal, "encoding publicKey error")
	}

	return &pb.NewSessionIDResponce{SessionID: sessionID, PublicKeyBZ: publicKeyBZ.Bytes()}, nil
}

// NewUser создает нового пользователя.
func (s *GophKeeperServer) NewUser(ctx context.Context, in *pb.NewUserRequest) (*pb.NewUserResponce, error) {
	_, err := s.rsa.CheckSign(in.SessionID, in.UserSign)
	if err != nil {
		log.Error().Err(err).Msg("NewUser CheckSign error")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}

	userLogin, userPass, err := s.rsa.DecryptLogin(in.SessionID, in.NewUser)
	if errors.Is(err, gkerrors.ErrLoginIncorrect) {
		return nil, status.Error(codes.InvalidArgument, "users login contains error")
	}
	if err != nil {
		log.Error().Err(err).Msg("NewUser EncryptLogin error")
		return nil, status.Error(codes.Internal, "EncryptLogin error")
	}

	exist, err := s.strg.CheckUser(userLogin)
	if exist {
		log.Error().Msgf("NewUser exists error")
		return nil, status.Error(codes.AlreadyExists, "user with such login exists")
	}
	if err != nil {
		log.Error().Err(err).Msg("NewUser CheckUser error")
		return nil, status.Error(codes.Internal, "CheckUser error")
	}

	userID, symKey, timeStamp, err := s.strg.RegisterUser(userLogin, userPass)
	if err != nil {
		log.Error().Err(err).Msg("NewUser RegisterUser error")
		return nil, status.Error(codes.Internal, "RegisterUser error")
	}
	s.rsa.AddUserID(in.SessionID, userID)

	var responce = pb.NewUserResponce{TimeStamp: timeStamp}
	responce.UserID, err = s.rsa.EncryptData(in.SessionID, userID, []byte(`userID`))
	if err != nil {
		log.Error().Err(err).Msg("NewUser EncryptData error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}

	responce.SymKey, err = s.rsa.EncryptData(in.SessionID, symKey, []byte(`key`))
	if err != nil {
		log.Error().Err(err).Msg("NewUser EncryptOAEP SymmetricalKey error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}

	responce.Sign, err = s.rsa.SignData(in.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("NewUser EncryptOAEP signing error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}

	return &responce, nil
}

// LoginUser авторизует пользователя.
func (s *GophKeeperServer) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponce, error) {
	_, err := s.rsa.CheckSign(in.SessionID, in.UserSign)
	if err != nil {
		log.Error().Err(err).Msg("NewUser CheckSign error")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}

	userLogin, userPass, err := s.rsa.DecryptLogin(in.SessionID, in.LoginUser)
	if errors.Is(err, gkerrors.ErrLoginIncorrect) {
		return nil, status.Error(codes.InvalidArgument, "users login contains error")
	}
	if err != nil {
		log.Error().Err(err).Msg("NewUser DecryptLogin error")
		return nil, status.Error(codes.Internal, "DecryptLogin error")
	}

	userID, err := s.strg.AuthUser(userLogin, userPass)
	if errors.Is(err, gkerrors.ErrNoSuchUser) {
		return nil, status.Error(codes.NotFound, "user with such login not registered")
	}
	if errors.Is(err, gkerrors.ErrWrongPassword) {
		return nil, status.Error(codes.InvalidArgument, "password incorrect")
	}
	if err != nil {
		log.Error().Err(err).Msg("LoginUser AuthUser error")
		return nil, status.Error(codes.Internal, "AuthUser error")
	}

	s.rsa.AddUserID(in.SessionID, userID)
	var responce pb.LoginUserResponce
	responce.UserID, err = s.rsa.EncryptData(in.SessionID, userID, []byte(`userID`))
	if err != nil {
		log.Error().Err(err).Msg("LoginUser EncryptData error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}

	responce.Sign, err = s.rsa.SignData(in.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("LoginUser EncryptOAEP signing error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}

	return &responce, nil
}

// UserData передает клиенту сохраненные данные пользователя.
func (s *GophKeeperServer) UserData(ctx context.Context, in *pb.UserDataRequest) (*pb.UserDataResponce, error) {
	userID, err := s.rsa.CheckSign(in.SessionID, in.UserSign)
	if err != nil {
		log.Error().Err(err).Msg("UserData CheckSign error")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	if userID == "" {
		log.Error().Err(err).Msg("UserData userID empty")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	var responce pb.UserDataResponce
	userData, timeStamp, symKey, err := s.strg.UsersData(userID)
	if errors.Is(err, gkerrors.ErrNoUserData) {
		responce.TimeStamp = timeStamp
		responce.SymKey, err = s.rsa.EncryptData(in.SessionID, symKey, []byte(`key`))
		if err != nil {
			log.Error().Err(err).Msg("UserData EncryptOAEP SymmetricalKey error")
			return nil, status.Error(codes.Internal, "EncryptData error")
		}
		responce.Sign, err = s.rsa.SignData(in.SessionID)
		if err != nil {
			log.Error().Err(err).Msg("UserData EncryptOAEP signing error")
			return nil, status.Error(codes.Internal, "EncryptData error")
		}
		return &responce, nil //status.Error(codes.NotFound, "user hasn't saved data")
	}
	if err != nil {
		log.Error().Err(err).Msg("UserData getData error")
		return nil, status.Error(codes.Internal, "getData error")
	}
	responce.UserData = userData
	responce.TimeStamp = timeStamp
	responce.SymKey, err = s.rsa.EncryptData(in.SessionID, symKey, []byte(`key`))
	if err != nil {
		log.Error().Err(err).Msg("UserData EncryptOAEP SymmetricalKey error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}

	responce.Sign, err = s.rsa.SignData(in.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("UserData EncryptOAEP signing error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}
	return &responce, nil
}

// TimeStamp передает клиенту отметку времени о последних сохраненных данных пользователя.
func (s *GophKeeperServer) TimeStamp(ctx context.Context, in *pb.TimeStampRequest) (*pb.TimeStampResponce, error) {
	userID, err := s.rsa.CheckSign(in.SessionID, in.UserSign)
	if err != nil {
		log.Error().Err(err).Msg("TimeStamp CheckSign error")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	if userID == "" {
		log.Error().Err(err).Msg("UserData userID empty")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}

	timeStamp, locked, timeLocked, err := s.strg.UsersTimeStamp(userID)
	if err != nil {
		log.Error().Err(err).Msg("TimeStamp error")
		return nil, status.Error(codes.Internal, "TimeStamp error")
	}

	var responce = pb.TimeStampResponce{TimeStamp: timeStamp, Locked: locked, TimeLocked: timeLocked}
	responce.Sign, err = s.rsa.SignData(in.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("UserData EncryptOAEP signing error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}

	return &responce, nil
}

// DataLock помечает данные клиента, как заблокированные на изменение другими пользователями
func (s *GophKeeperServer) DataLock(ctx context.Context, in *pb.DataLockRequest) (*pb.DataLockResponce, error) {
	userID, err := s.rsa.CheckSign(in.SessionID, in.UserSign)
	if err != nil {
		log.Error().Err(err).Msg("TimeStamp CheckSign error")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	if userID == "" {
		log.Error().Err(err).Msg("UserData userID empty")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}

	locked, timeLocked := s.strg.UsersDataLock(userID, in.SessionID)
	var responce = pb.DataLockResponce{Locked: locked, TimeLocked: timeLocked}
	responce.Sign, err = s.rsa.SignData(in.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("UserData EncryptOAEP signing error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}

	return &responce, nil
}

// UpdateData принимает от клиента обновленные данные пользователя.
func (s *GophKeeperServer) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.UpdateDataResponce, error) {
	log.Debug().Msgf("Данные для сохранения пришли на сервер")
	userID, err := s.rsa.CheckSign(in.SessionID, in.UserSign)
	if err != nil {
		log.Error().Err(err).Msg("UpdateData CheckSign error")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	if userID == "" {
		log.Error().Err(err).Msg("UserData userID empty")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}

	save, timeStamp, err := s.strg.UpdateUserData(userID, in.SessionID, in.TimeStamp, in.UserData)
	if errors.Is(err, gkerrors.ErrLocked) {
		return nil, status.Error(codes.PermissionDenied, "users data changes locked by another user")
	}
	if errors.Is(err, gkerrors.ErrTimeNotEqual) {
		return nil, status.Error(codes.FailedPrecondition, "users data timeStamp not equal to servers")
	}
	if err != nil {
		log.Error().Err(err).Msg("UpdateData error")
		return nil, status.Error(codes.Internal, "UpdateData error")
	}

	var responce = pb.UpdateDataResponce{Status: save, TimeStamp: timeStamp}
	responce.Sign, err = s.rsa.SignData(in.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("UserData EncryptOAEP signing error")
		return nil, status.Error(codes.Internal, "EncryptData error")
	}

	return &responce, nil
}

// LogOut закрывает сессию пользователя.
func (s *GophKeeperServer) LogOut(ctx context.Context, in *pb.LogOutRequest) (*pb.LogOutResponce, error) {
	_, err := s.rsa.CheckSign(in.SessionID, in.UserSign)
	if err != nil {
		log.Error().Err(err).Msg("LogOut CheckSign error")
		return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
	}
	s.rsa.UserLogOut(in.SessionID)
	return &pb.LogOutResponce{Status: true}, nil
}
