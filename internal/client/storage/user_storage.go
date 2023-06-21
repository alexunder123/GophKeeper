// Модуль предназначен для хранения в оперативной памяти данных пользователя на клиенте.

package storage

import (
	"encoding/json"
	"time"
)

// Password структура для хранения логинов и паролей клиента.
type Password struct {
	Name    string
	Login   string
	Pass    string
	Comment string
}

// Card структура для хранения данных карты клиента.
type Card struct {
	Name       string
	CardNumber string
	Comment    string
}

// Text структура для хранения текстовых данных клиента.
type Text struct {
	Name    string
	Data    string
	Comment string
}

// Binary структура для хранения произвольных данных клиента.
type Binary struct {
	Name    string
	Data    []byte
	Comment string
}

// UserStorage структура для хранения данных на клиенте.
type UserStorage struct {
	TimeStamp  time.Time  `json:"-"`
	Passwords  []Password `json:"passwords"`
	Cards      []Card     `json:"cards"`
	Texts      []Text     `json:"texts"`
	Binaries   []Binary   `json:"binaries"`
	Locked     bool       `json:"-"`
	TimeLocked time.Time  `json:"-"`
}

// NewUserStorage метод генерирует хранилище оперативных данных.
func NewUserStorage() *UserStorage {
	return &UserStorage{
		Passwords: make([]Password, 0),
		Cards:     make([]Card, 0),
		Texts:     make([]Text, 0),
		Binaries:  make([]Binary, 0),
	}
}

// ImportUserData метод раскодирует и сохраняет в хранилище данные пользователя.
func (s *UserStorage) ImportUserData(jsonBZ []byte, timeStamp string) error {
	var err error
	s.TimeStamp, err = time.Parse(time.RFC3339, timeStamp)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBZ, s)
	if err != nil {
		return err
	}
	return nil
}

// ExportUserData метод кодирует данные пользователя для отправки.
func (s *UserStorage) ExportUserData() ([]byte, error) {
	jsonBZ, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return jsonBZ, nil
}

// ListUserData метод возвращает количество сохраненных строк пользователя.
func (s *UserStorage) ListUserData() (passwords, cards, texts, binaries int) {
	return len(s.Passwords), len(s.Cards), len(s.Texts), len(s.Binaries)
}

// SliceUserData метод возвращает информацию о сохраненных строках пользователя в разделе паролей.
func (s *UserStorage) SliceUsersPasswords() []Password {
	return s.Passwords
}

// SliceUserData метод возвращает информацию о сохраненных строках пользователя в разделе карт.
func (s *UserStorage) SliceUsersCards() []Card {
	return s.Cards
}

// SliceUserData метод возвращает информацию о сохраненных строках пользователя в разделе текстов.
func (s *UserStorage) SliceUsersTexts() []Text {
	return s.Texts
}

// SliceUserData метод возвращает информацию о сохраненных строках пользователя в разделе двоичных данных.
func (s *UserStorage) SliceUsersBinaries() []Binary {
	return s.Binaries
}

// StringUsersPassword метод возвращает полную информацию о сохраненной строке записи с паролем.
func (s *UserStorage) StringUsersPassword(v int) *Password {
	if v >= len(s.Passwords) {
		return nil
	}
	return &s.Passwords[v]
}

// StringUsersCard метод возвращает полную информацию о сохраненной строке записи с картами.
func (s *UserStorage) StringUsersCard(v int) *Card {
	if v >= len(s.Cards) {
		return nil
	}
	return &s.Cards[v]
}

// StringUsersText метод возвращает полную информацию о сохраненной строке записи с текстовыми данными.
func (s *UserStorage) StringUsersText(v int) *Text {
	if v >= len(s.Texts) {
		return nil
	}
	return &s.Texts[v]
}

// StringUsersBinary метод возвращает полную информацию о сохраненной строке записи с двоичными данными.
func (s *UserStorage) StringUsersBinary(v int) *Binary {
	if v >= len(s.Binaries) {
		return nil
	}
	return &s.Binaries[v]
}

// AddUserData метод добавляет новую запись с паролем.
func (s *UserStorage) AddUsersPassword(pass *Password) {
	s.Passwords = append(s.Passwords, *pass)
}

// AddUserData метод добавляет новую запись с картами.
func (s *UserStorage) AddUsersCard(card *Card) {
	s.Cards = append(s.Cards, *card)
}

// AddUserData метод добавляет новую запись с текстом.
func (s *UserStorage) AddUsersText(text *Text) {
	s.Texts = append(s.Texts, *text)
}

// AddUserData метод добавляет новую запись с двоичными данными.
func (s *UserStorage) AddUsersBinary(binary *Binary) {
	s.Binaries = append(s.Binaries, *binary)
}

// EditUsersPassword метод редактирует существующую запись с данными пароля.
func (s *UserStorage) EditUsersPassword(i int, password *Password) {
	s.Passwords[i] = *password
}

// EditUsersCard метод редактирует существующую запись с данными карты.
func (s *UserStorage) EditUsersCard(i int, card *Card) {
	s.Cards[i] = *card
}

// EditUsersCard метод редактирует существующую запись с данными текста.
func (s *UserStorage) EditUsersText(i int, text *Text) {
	s.Texts[i] = *text
}

// EditUsersBinary метод редактирует существующую запись с двоичными данными.
func (s *UserStorage) EditUsersBinary(i int, binary *Binary) {
	s.Binaries[i] = *binary
}
