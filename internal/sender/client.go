// Модуль предназначен для формирования и отправки запросов клиентом на сервер.
// Модуль принимает возвращаемые данные и производит первичную их обработку.
package sender

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"gophkeeper/internal/crypto"
	gkerrors "gophkeeper/internal/errors"
	pb "gophkeeper/internal/grpc/proto"
	"gophkeeper/internal/storage"
)

// GophKeeperClient поддерживает все необходимые методы клиента.
type GophKeeperClient struct {
	// pb.GophKeeperClient
	cc   pb.GophKeeperClient
	rsa  *crypto.UserSession
	Strg *storage.UserStorage
}

// NewGophKeeperClient генерирует структуру для gRPC клиента.
func NewGophKeeperClient(cc grpc.ClientConnInterface, rsa *crypto.UserSession, strg *storage.UserStorage) GophKeeperClient {
	cci := pb.NewGophKeeperClient(cc)
	return GophKeeperClient{cc: cci, rsa: rsa, Strg: strg}
}

// RefreshToken метод обновляет ключи сессии
func (c *GophKeeperClient) RefreshToken() error {
	return c.rsa.RefreshToken()
}

// ReqSessionID метод запрашивает у сервера сессию и обменивается открытыми ключами
func (c GophKeeperClient) ReqSessionID() error {
	var publicKeyBZ bytes.Buffer
	publicKey := c.rsa.GetPublicKey()
	enc := gob.NewEncoder(&publicKeyBZ)
	err := enc.Encode(&publicKey)
	if err != nil {
		log.Error().Err(err).Msg("ReqSessionID encoding publicKey error")
		return err
	}
	var request = pb.NewSessionIDRequest{UserPublicKeyBZ: publicKeyBZ.Bytes()}
	responce, err := c.cc.NewSessionID(context.Background(), &request)
	if err != nil {
		log.Error().Err(err).Msg("ReqSessionID encoding publicKey error")
		return err
	}
	sessionID := responce.SessionID
	var buf bytes.Buffer
	_, err = buf.Write(responce.PublicKeyBZ)
	if err != nil {
		log.Error().Err(err).Msg("NewSessionID writing to buf error")
		return err
	}
	var serverPublicKey = rsa.PublicKey{}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&serverPublicKey)
	if err != nil {
		log.Error().Err(err).Msg("NewSessionID decoding publicKey error")
		return err
	}
	c.rsa.WriteSessionID(sessionID, &serverPublicKey)
	return nil
}

// RegisterUser метод формирует и отправляет запрос на регистрацию нового пользователя
func (c *GophKeeperClient) RegisterUser(login, pass string) error {
	message, err := c.rsa.EncryptData(login+","+pass, []byte("login"))
	if err != nil {
		log.Error().Err(err).Msg("RegisterUser EncryptData error")
		return err
	}
	var request = pb.NewUserRequest{SessionID: c.rsa.GetSessionID(), NewUser: message}
	request.UserSign, err = c.rsa.UserSign()
	if err != nil {
		log.Error().Err(err).Msg("RegisterUser EncryptOAEP signing error")
		return err
	}
	responce, err := c.cc.NewUser(context.Background(), &request)
	if err != nil {
		return err
	}
	if c.rsa.CheckSign(responce.Sign) != nil {
		return gkerrors.ErrSignIncorrect
	}
	userID, err := c.rsa.DecryptData(responce.UserID, []byte("userID"))
	if err != nil {
		return err
	}
	symKey, err := c.rsa.DecryptData(responce.SymKey, []byte("key"))
	if err != nil {
		return err
	}
	c.rsa.WriteUserID(userID, symKey)
	c.Strg.TimeStamp, err = time.Parse(time.RFC3339, responce.TimeStamp)
	if err != nil {
		return err
	}
	return nil
}

// UserLogin метод формирует и отправляет запрос на авторизацию пользователя
func (c *GophKeeperClient) UserLogin(login, pass string) error {
	message, err := c.rsa.EncryptData(login+","+pass, []byte("login"))
	if err != nil {
		log.Error().Err(err).Msg("RegisterUser EncryptData error")
		return err
	}
	var request = pb.LoginUserRequest{SessionID: c.rsa.GetSessionID(), LoginUser: message}
	request.UserSign, err = c.rsa.UserSign()
	if err != nil {
		log.Error().Err(err).Msg("RegisterUser EncryptOAEP signing error")
		return err
	}
	responce, err := c.cc.LoginUser(context.Background(), &request)
	if err != nil {
		return err
	}
	if c.rsa.CheckSign(responce.Sign) != nil {
		return gkerrors.ErrSignIncorrect
	}
	userID, err := c.rsa.DecryptData(responce.UserID, []byte("userID"))
	if err != nil {
		return err
	}
	c.rsa.WriteUserID(userID, "")
	return nil
}

// CheckTimeStamp метод запрашивает и сравнивает время последнего сохранения данных пользователя.
func (c *GophKeeperClient) CheckTimeStamp() error {
	var request = pb.TimeStampRequest{SessionID: c.rsa.GetSessionID()}
	var err error
	request.UserSign, err = c.rsa.UserSign()
	if err != nil {
		log.Error().Err(err).Msg("RegisterUser EncryptOAEP signing error")
		return err
	}
	responce, err := c.cc.TimeStamp(context.Background(), &request)
	if err != nil {
		return err
	}
	if c.rsa.CheckSign(responce.Sign) != nil {
		return gkerrors.ErrSignIncorrect
	}
	if responce.TimeStamp != c.Strg.TimeStamp.Format(time.RFC3339) {
		fmt.Printf("Время последнего сохранения на сервере и клиенте не совпадают. На сервере = %s, на клиенте = %s", responce.TimeStamp, c.Strg.TimeStamp.Format(time.RFC3339))
	} else {
		fmt.Println("Время последнего сохранения на сервере и клиенте совпадают")
	}
	if responce.Locked {
		fmt.Printf("Данные на сервере заблокированы на изменение другим пользователем до: %s", responce.TimeLocked)
	}
	return nil
}

// LockUserData метод запрашивает на сервере временную блокировку на изменение данных клиента другими пользователями.
func (c *GophKeeperClient) LockUserData() error {
	if c.Strg.Locked && c.Strg.TimeLocked.After(time.Now()) {
		return nil
	}
	var request = pb.DataLockRequest{SessionID: c.rsa.GetSessionID()}
	var err error
	request.UserSign, err = c.rsa.UserSign()
	if err != nil {
		log.Error().Err(err).Msg("RegisterUser EncryptOAEP signing error")
		return err
	}
	responce, err := c.cc.DataLock(context.Background(), &request)
	if err != nil {
		return err
	}
	if c.rsa.CheckSign(responce.Sign) != nil {
		return gkerrors.ErrSignIncorrect
	}
	if !responce.Locked {
		fmt.Printf("Данные на сервере заблокированы на изменение другим пользователем до: %s\n", responce.TimeLocked)
		return gkerrors.ErrLocked
	}
	c.Strg.TimeLocked, err = time.Parse(time.RFC3339, responce.TimeLocked)
	if err != nil {
		return err
	}
	c.Strg.Locked = true
	fmt.Printf("Данные на сервере успешно заблокированы на изменение до: %s\n", responce.TimeLocked)
	return nil
}

// Download метод запрашивает на сервере сохраненные данные пользователя.
func (c *GophKeeperClient) Download() error {
	var request = pb.UserDataRequest{SessionID: c.rsa.GetSessionID()}
	var err error
	request.UserSign, err = c.rsa.UserSign()
	if err != nil {
		log.Error().Err(err).Msg("RegisterUser EncryptOAEP signing error")
		return err
	}
	responce, err := c.cc.UserData(context.Background(), &request)
	if err != nil {
		return err
	}
	if c.rsa.CheckSign(responce.Sign) != nil {
		return gkerrors.ErrSignIncorrect
	}
	symKey, err := c.rsa.DecryptData(responce.SymKey, []byte("key"))
	if err != nil {
		return err
	}
	c.rsa.WriteUserID("", symKey)
	if len(responce.UserData) == 0 {
		c.Strg.TimeStamp, err = time.Parse(time.RFC3339, responce.TimeStamp)
		if err != nil {
			return err
		}
		fmt.Println("На сервере нет сохраненных данных клиента")
		return nil
	}
	jsonBZ, err := c.rsa.DecryptUserData(responce.UserData)
	if err != nil {
		return err
	}
	err = c.Strg.ImportUserData(jsonBZ, responce.TimeStamp)
	if err != nil {
		return err
	}
	fmt.Println("Данные успешно скачаны с сервера")
	if responce.Locked {
		fmt.Printf("Данные на сервере заблокированы на изменение другим пользователем до: %s\n", responce.TimeLocked)
	}
	return nil
}

// SaveData метод отправляет на сервер данные пользователя для сохранения.
func (c *GophKeeperClient) SaveData() error {
	jsonBZ, err := c.Strg.ExportUserData()
	if err != nil {
		return err
	}
	messageBZ, err := c.rsa.EncryptUserData(jsonBZ)
	if err != nil {
		return err
	}
	var request = pb.UpdateDataRequest{SessionID: c.rsa.GetSessionID(), TimeStamp: c.Strg.TimeStamp.Format(time.RFC3339), UserData: messageBZ}
	request.UserSign, err = c.rsa.UserSign()
	if err != nil {
		log.Error().Err(err).Msg("RegisterUser EncryptOAEP signing error")
		return err
	}
	log.Debug().Msgf("Данные для сохранения на сервер отправлены")
	responce, err := c.cc.UpdateData(context.Background(), &request)
	if errors.Is(err, gkerrors.ErrLocked) {
		fmt.Printf("Данные на сервере заблокированы на изменение другим пользователем до: %s\n", responce.TimeStamp)
		return nil
	}
	if err != nil {
		return err
	}
	if c.rsa.CheckSign(responce.Sign) != nil {
		return gkerrors.ErrSignIncorrect
	}
	c.Strg.TimeStamp, err = time.Parse(time.RFC3339, responce.TimeStamp)
	if err != nil {
		return err
	}
	c.Strg.Locked = false
	return nil
}

// UserLogOut метод очищает данные пользовательской сессии и отправляет на сервер запрос на удаление сессии.
func (c *GophKeeperClient) UserLogOut() {
	var request = pb.LogOutRequest{SessionID: c.rsa.GetSessionID()}
	var err error
	var noSend bool
	request.UserSign, err = c.rsa.UserSign()
	if err != nil {
		log.Error().Err(err).Msg("UserLogOut EncryptOAEP signing error")
		noSend = true
	}
	err = c.rsa.RefreshToken()
	if err != nil {
		log.Error().Err(err).Msg("UserLogOut RefreshToken error")
	}
	c.Strg = storage.NewUserStorage()
	if noSend {
		return
	}
	responce, err := c.cc.LogOut(context.Background(), &request)
	if err != nil {
		log.Error().Err(err).Msg("UserLogOut LogOut error")
		return
	}
	if !responce.Status {
		fmt.Println("Не удалось удалить сессию на сервере")
	}
}
