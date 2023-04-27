package errors

import "errors"

// Переменные для передачи хэндлену идентификатора ошибки.
var (
	ErrExpired error = errors.New("SessionID expired")
)
