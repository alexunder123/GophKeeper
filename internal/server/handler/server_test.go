package handler

import (
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"gophkeeper/api/grpc/proto"
	clientCNFG "gophkeeper/internal/client/config"
	clientCRPT "gophkeeper/internal/client/crypto"
	"gophkeeper/internal/client/sender"
	clientSTRG "gophkeeper/internal/client/storage"
	gkerrors "gophkeeper/internal/errors"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/mocks"
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/crypto"
)

func TestServer(t *testing.T) {
	logger.Newlogger()

	// Конфигурируем сервер
	cnfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("NewConfig read environment error")
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	strg := mocks.NewMockStorager(ctrl)
	rsa := crypto.NewSessions(cnfg)
	gRPCconf := NewGophKeeperServer(cnfg, strg, rsa)
	listen, err := net.Listen("tcp", cnfg.RunAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("gRPC server announce error")
	}
	server := grpc.NewServer()
	proto.RegisterGophKeeperServer(server, gRPCconf)
	reflection.Register(server)
	go func() {
		if err := server.Serve(listen); err != nil {
			log.Fatal().Msgf("gRPC server failed: %s", err)
		}
	}()

	// Конфигурируем клиент
	clientCnfg, err := clientCNFG.NewUserConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("NewConfig read environment error")
	}
	clientStrg := clientSTRG.NewUserStorage()
	clientRsa, err := clientCRPT.NewUserSession()
	if err != nil {
		log.Fatal().Err(err).Msg("NewUserSession generating key error")
	}
	conn, err := grpc.Dial(clientCnfg.RunAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("gRPC connection error")
	}
	client := sender.NewGophKeeperClient(conn, clientRsa, clientStrg)

	// Начинаем тестирование

	// Проверка актуальности данных без запроса сессии
	err = client.CheckTimeStamp()
	require.Error(t, err)

	// Проверка на запрос данных без запроса сессии
	err = client.Download()
	require.Error(t, err)

	// Попытка блокировки данных без запроса сессии
	err = client.LockUserData()
	require.Error(t, err)

	// Проверка регистрации пользователя
	// Запрос сессии
	err = client.ReqSessionID()
	require.NoError(t, err)

	// Проверка регистрации при занятом логине
	strg.EXPECT().CheckUser("userName1").Return(true, gkerrors.ErrLoginExist)
	err = client.RegisterUser("userName1", "123")
	require.Error(t, err)

	// Успешная регистрация
	userID := crypto.RandomID(12)
	symKey := crypto.NewSymmetricalKey(userID)
	timeStamp := time.Now().Format(time.RFC3339)
	timeLock := time.Now().Add(time.Minute * 5).Format(time.RFC3339)
	first := strg.EXPECT().CheckUser("userName2").Return(false, nil).MaxTimes(1)
	strg.EXPECT().RegisterUser("userName2", crypto.HashPasswd("123")).Return(userID, symKey, timeStamp, nil).After(first)
	err = client.RegisterUser("userName2", "123")
	require.NoError(t, err)

	// Проверка на скачивание данных, при их отсутствии
	strg.EXPECT().UsersData(userID).Return(nil, timeStamp, symKey, gkerrors.ErrNoUserData)
	err = client.Download()
	require.NoError(t, err)

	// Закрытие сессии
	client.UserLogOut()

	// Проверка авторизации пользователя и обмена данными с пользователем
	// Запрос сессии
	err = client.ReqSessionID()
	require.NoError(t, err)

	// Авторизация несуществующим пользователем
	strg.EXPECT().AuthUser("userName3", crypto.HashPasswd("123")).Return("", gkerrors.ErrNoSuchUser)
	err = client.UserLogin("userName3", "123")
	require.Error(t, err)

	// Проверка на запрос без авторизации
	err = client.Download()
	require.Error(t, err)

	// Проверка актуальности данных без авторизации
	err = client.CheckTimeStamp()
	require.Error(t, err)

	// Попытка блокировки данных без авторизации
	err = client.LockUserData()
	require.Error(t, err)

	// Попытка записи без авторизации
	err = client.SaveData()
	require.Error(t, err)

	// Попытка смены пароля без авторизации
	status, err := client.ChangePassword("123", "456")
	require.Error(t, err)
	require.Equal(t, false, status)

	// Проверка авторизации с неверным паролем
	strg.EXPECT().AuthUser("userName4", crypto.HashPasswd("123")).Return("", gkerrors.ErrWrongPassword)
	err = client.UserLogin("userName4", "123")
	require.Error(t, err)

	// Успешная авторизация
	strg.EXPECT().AuthUser("userName5", crypto.HashPasswd("234")).Return("1234567890", nil)
	err = client.UserLogin("userName5", "234")
	require.NoError(t, err)

	// Успешное скачивание данных
	strg.EXPECT().UsersData("1234567890").Return(nil, timeStamp, symKey, gkerrors.ErrNoUserData)
	err = client.Download()
	require.NoError(t, err)

	// Проверка актуальности данных при наличии блокировки
	strg.EXPECT().UsersTimeStamp("1234567890").Return(timeStamp, true, timeLock, nil)
	err = client.CheckTimeStamp()
	require.NoError(t, err)

	// Попытка блокировки данных, при ее наличии
	strg.EXPECT().UsersDataLock("1234567890", clientRsa.GetSessionID()).Return(false, timeLock)
	err = client.LockUserData()
	require.Error(t, err)

	jsonBZ, _ := client.Strg.ExportUserData()
	messageBZ, _ := clientRsa.EncryptUserData(jsonBZ)
	// Попытка записи при наличии блокировки данных
	strg.EXPECT().UpdateUserData("1234567890", clientRsa.GetSessionID(), timeStamp, messageBZ).Return(false, timeLock, gkerrors.ErrLocked)
	err = client.SaveData()
	require.Error(t, err)

	// Попытка записи при неактуальных данных
	strg.EXPECT().UpdateUserData("1234567890", clientRsa.GetSessionID(), timeStamp, messageBZ).Return(false, timeLock, gkerrors.ErrTimeNotEqual)
	err = client.SaveData()
	require.Error(t, err)

	// Проверка актуальности данных при отсутствии блокировки
	strg.EXPECT().UsersTimeStamp("1234567890").Return(timeStamp, false, "", nil)
	err = client.CheckTimeStamp()
	require.NoError(t, err)

	// Успешная блокировка данных
	strg.EXPECT().UsersDataLock("1234567890", clientRsa.GetSessionID()).Return(true, timeLock)
	err = client.LockUserData()
	require.NoError(t, err)

	// Успешная запись данных
	strg.EXPECT().UpdateUserData("1234567890", clientRsa.GetSessionID(), timeStamp, messageBZ).Return(true, timeStamp, nil)
	err = client.SaveData()
	require.NoError(t, err)

	// Попытка смены пароля при неверном пароле
	strg.EXPECT().ChangeUserPassword("1234567890", crypto.HashPasswd("456"), crypto.HashPasswd("123")).Return(false, gkerrors.ErrWrongPassword)
	status, err = client.ChangePassword("456", "123")
	require.Error(t, err)
	require.Equal(t, false, status)

	// Успешная смена пароля
	strg.EXPECT().ChangeUserPassword("1234567890", crypto.HashPasswd("123"), crypto.HashPasswd("456")).Return(true, nil)
	status, err = client.ChangePassword("123", "456")
	require.NoError(t, err)
	require.Equal(t, true, status)

	// Закрытие сессии
	client.UserLogOut()

	// Закрытие удаленной сессии (возврат ошибки)
	client.UserLogOut()

	// Оканчиваем тестирование

	strg.EXPECT().CloseDB()
	server.GracefulStop()
	strg.CloseDB()
}
