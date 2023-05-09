// Модуль предназначен для создания пользовательских ошибок.
package gkerrors

import "errors"

// Переменные для передачи хэндлену идентификатора ошибки.
var (
	ErrExpired        error = errors.New("sessionID expired")
	ErrLoginExist     error = errors.New("user with such login exists")
	ErrNoSuchUser     error = errors.New("user with such login not registered")
	ErrLoginIncorrect error = errors.New("users login contains error")
	ErrWrongPassword  error = errors.New("password incorrect")
	ErrNotAuth        error = errors.New("user isn't authenticated")
	ErrNoUserData     error = errors.New("user hasn't saved data")
	ErrTimeNotEqual   error = errors.New("users data timeStamp not equal to servers")
	ErrLocked         error = errors.New("users data changes locked by another user")
	ErrSignIncorrect  error = errors.New("incorrect sign encryption")
	ErrTooBig         error = errors.New("file is too big")
)
