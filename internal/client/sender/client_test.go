package sender

import (
	"testing"

	"gophkeeper/internal/client/config"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/client/storage"
	gkerrors "gophkeeper/internal/errors"
	"gophkeeper/internal/mocks"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestClient(t *testing.T) {

	// Конфигурируем клиент
	cnfg, err := config.NewUserConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("NewConfig read environment error")
	}
	strg := storage.NewUserStorage()
	rsa, err := crypto.NewUserSession()
	if err != nil {
		log.Fatal().Err(err).Msg("NewUserSession generating key error")
	}
	conn, err := grpc.Dial(cnfg.RunAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(mocks.Dialer()))
	if err != nil {
		log.Fatal().Err(err).Msg("gRPC connection error")
	}
	client := NewGophKeeperClient(conn, rsa, strg)

	// Начинаем тестирование

	// Запрос сессии
	err = client.ReqSessionID()
	require.NoError(t, err)

	// Регистрация пользователя
	checkRegister(t, &client)

	// Запрос сессии
	err = client.RefreshToken()
	require.NoError(t, err)
	err = client.ReqSessionID()
	require.NoError(t, err)

	// Авторизация пользователя
	checkAuth(t, &client)

	// Скачивание данных
	checkDownload(t, &client)

	// Проверка актуальности данных
	checkDataTimeStamp(t, &client)

	// Блокировка данных
	checkDataLock(t, &client)

	// Сохранение данных
	checkSaveData(t, &client)

	// Смена пароля
	checkChangePassword(t, &client)

	// Закрытие сессии
	checkLogOut(t, &client)
}

func checkRegister(t *testing.T, client *GophKeeperClient) {
	tests := []struct {
		name    string
		login   string
		pass    string
		isError bool
		errCode codes.Code
		err     error
	}{
		{
			name:    "Регистрация. Ошибка подписи",
			login:   "123",
			pass:    "123",
			isError: true,
			errCode: codes.Unauthenticated,
		},
		{
			name:    "Регистрация. Логин занят",
			login:   "234",
			pass:    "123",
			isError: true,
			errCode: codes.AlreadyExists,
		},
		{
			name:    "Регистрация. Нет подписи сервера",
			login:   "345",
			pass:    "123",
			isError: true,
			err:     gkerrors.ErrSignIncorrect,
		},
		{
			name:    "Регистрация. Успешно",
			login:   "userName1",
			pass:    "123",
			isError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.RegisterUser(tt.login, tt.pass)
			if tt.isError {
				st, ok := status.FromError(err)
				if ok {
					require.Equal(t, tt.errCode, st.Code())
				} else {
					require.Equal(t, tt.err, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func checkAuth(t *testing.T, client *GophKeeperClient) {
	tests := []struct {
		name    string
		login   string
		pass    string
		isError bool
		errCode codes.Code
		err     error
	}{
		{
			name:    "Авторизация. Ошибка",
			login:   "123",
			pass:    "123",
			isError: true,
			errCode: codes.Unauthenticated,
		},
		{
			name:    "Авторизация. Нет подписи сервера",
			login:   "345",
			pass:    "123",
			isError: true,
			err:     gkerrors.ErrSignIncorrect,
		},
		{
			name:    "Авторизация. Успешно",
			login:   "userName1",
			pass:    "123",
			isError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.UserLogin(tt.login, tt.pass)
			if tt.isError {
				st, ok := status.FromError(err)
				if ok {
					require.Equal(t, tt.errCode, st.Code())
				} else {
					require.Equal(t, tt.err, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func checkDownload(t *testing.T, client *GophKeeperClient) {
	tests := []struct {
		name      string
		sessionID string
		isError   bool
		errCode   codes.Code
		err       error
	}{
		{
			name:      "Скачивание данных. Ошибка",
			sessionID: "9876543210",
			isError:   true,
			errCode:   codes.Unauthenticated,
		},
		{
			name:      "Скачивание данных. Нет подписи сервера",
			sessionID: "147852369",
			isError:   true,
			err:       gkerrors.ErrSignIncorrect,
		},
		{
			name:      "Скачивание данных. Успешно",
			sessionID: "1234567890",
			isError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.rsa.WriteSessionID(tt.sessionID, nil)
			err := client.Download()
			if tt.isError {
				st, ok := status.FromError(err)
				if ok {
					require.Equal(t, tt.errCode, st.Code())
				} else {
					require.Equal(t, tt.err, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func checkDataTimeStamp(t *testing.T, client *GophKeeperClient) {
	tests := []struct {
		name      string
		sessionID string
		isError   bool
		errCode   codes.Code
		err       error
	}{
		{
			name:      "Время сохранения данных. Ошибка",
			sessionID: "9876543210",
			isError:   true,
			errCode:   codes.Unauthenticated,
		},
		{
			name:      "Время сохранения данных. Нет подписи сервера",
			sessionID: "147852369",
			isError:   true,
			err:       gkerrors.ErrSignIncorrect,
		},
		{
			name:      "Время сохранения данных. Разное время + блокировка",
			sessionID: "123654789",
			isError:   false,
		},
		{
			name:      "Время сохранения данных. Успешно",
			sessionID: "123654890",
			isError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.rsa.WriteSessionID(tt.sessionID, nil)
			err := client.CheckTimeStamp()
			if tt.isError {
				st, ok := status.FromError(err)
				if ok {
					require.Equal(t, tt.errCode, st.Code())
				} else {
					require.Equal(t, tt.err, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func checkDataLock(t *testing.T, client *GophKeeperClient) {
	tests := []struct {
		name      string
		sessionID string
		isError   bool
		errCode   codes.Code
		err       error
	}{
		{
			name:      "Блокировка данных. Ошибка",
			sessionID: "9876543210",
			isError:   true,
			errCode:   codes.Unauthenticated,
		},
		{
			name:      "Блокировка данных. Нет подписи сервера",
			sessionID: "147852369",
			isError:   true,
			err:       gkerrors.ErrSignIncorrect,
		},
		{
			name:      "Блокировка данных. Заблокировано другим пользователем",
			sessionID: "123654789",
			isError:   true,
			err:       gkerrors.ErrLocked,
		},
		{
			name:      "Блокировка данных. Успешно",
			sessionID: "123654890",
			isError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.rsa.WriteSessionID(tt.sessionID, nil)
			err := client.LockUserData()
			if tt.isError {
				st, ok := status.FromError(err)
				if ok {
					require.Equal(t, tt.errCode, st.Code())
				} else {
					require.Equal(t, tt.err, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func checkSaveData(t *testing.T, client *GophKeeperClient) {
	tests := []struct {
		name      string
		sessionID string
		isError   bool
		errCode   codes.Code
		err       error
	}{
		{
			name:      "Сохранение данных. Ошибка",
			sessionID: "9876543210",
			isError:   true,
			errCode:   codes.Unauthenticated,
		},
		{
			name:      "Сохранение данных. Нет подписи сервера",
			sessionID: "147852369",
			isError:   true,
			err:       gkerrors.ErrSignIncorrect,
		},
		{
			name:      "Сохранение данных. Успешно",
			sessionID: "123654890",
			isError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.rsa.WriteSessionID(tt.sessionID, nil)
			err := client.SaveData()
			if tt.isError {
				st, ok := status.FromError(err)
				if ok {
					require.Equal(t, tt.errCode, st.Code())
				} else {
					require.Equal(t, tt.err, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func checkChangePassword(t *testing.T, client *GophKeeperClient) {
	tests := []struct {
		name    string
		oldPass string
		newPass string
		isError bool
		errCode codes.Code
		err     error
	}{
		{
			name:    "Смена пароля. Ошибка",
			oldPass: "234",
			newPass: "456",
			isError: true,
			errCode: codes.Unauthenticated,
		},
		{
			name:    "Смена пароля. Нет подписи сервера",
			oldPass: "345",
			newPass: "456",
			isError: true,
			err:     gkerrors.ErrSignIncorrect,
		},
		{
			name:    "Смена пароля. Успешно",
			oldPass: "123",
			newPass: "456",
			isError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.ChangePassword(tt.oldPass, tt.newPass)
			if tt.isError {
				st, ok := status.FromError(err)
				if ok {
					require.Equal(t, tt.errCode, st.Code())
				} else {
					require.Equal(t, tt.err, err)
				}
			} else {
				require.Equal(t, true, result)
			}
		})
	}
}

func checkLogOut(t *testing.T, client *GophKeeperClient) {
	tests := []struct {
		name      string
		sessionID string
		isError   bool
		errCode   codes.Code
		err       error
	}{
		{
			name:      "Удаление сессии. Ошибка",
			sessionID: "9876543210",
			isError:   true,
			errCode:   codes.Unauthenticated,
		},
		{
			name:      "Удаление сессии. Сессия уже удалена",
			sessionID: "147852369",
			isError:   false,
		},
		{
			name:      "Сохранение данных. Успешно",
			sessionID: "123654890",
			isError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.rsa.WriteSessionID(tt.sessionID, nil)
			client.UserLogOut()
		})
	}
}
