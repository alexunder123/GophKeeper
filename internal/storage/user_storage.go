package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/internal/crypto"
	gkerrors "gophkeeper/internal/errors"
	"os"
	"time"

	"github.com/rs/zerolog/log"
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

// ExportUserData метод раскодирует и сохраняет в хранилище данные пользователя.
func (s *UserStorage) ExportUserData() ([]byte, error) {
	jsonBZ, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return jsonBZ, nil
}

// ListUserData метод выводит на экран количество сохраненных строк пользователя.
func (s *UserStorage) ListUserData() {
	fmt.Printf(`В базе содержится следующее количество записей:
	Количество сохраненных паролей: %d
	Количество сохраненных карт: %d
	Количество сохраненных произвольных текстов: %d
	Количество сохраненных бинарных данных: %d
	`, len(s.Passwords), len(s.Cards), len(s.Texts), len(s.Binaries))
}

// SliceUserData метод выводит на экран информацию о сохраненных строках пользователя в соответствующем разделе.
func (s *UserStorage) SliceUserData(i int) {
	switch i {
	case 1:
		for i, val := range s.Passwords {
			fmt.Printf("Номер: %d, Имя: %s\n", i+1, val.Name)
		}
	case 2:
		for i, val := range s.Cards {
			fmt.Printf("Номер: %d, Имя: %s\n", i+1, val.Name)
		}
	case 3:
		for i, val := range s.Texts {
			fmt.Printf("Номер: %d, Имя: %s\n", i+1, val.Name)
		}
	case 4:
		for i, val := range s.Binaries {
			fmt.Printf("Номер: %d, Имя: %s\n", i+1, val.Name)
		}
	}
}

// StringUserData метод выводит на экран полную информацию о сохраненной записи.
func (s *UserStorage) StringUserData(i, v int) {
	switch i {
	case 1:
		if v >= len(s.Passwords) {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		val := s.Passwords[v]
		fmt.Printf(`Данные строки:
		Имя: %s Логин: %s Пароль: %s Примечание: %s
		`, val.Name, val.Login, val.Pass, val.Comment)
	case 2:
		if v >= len(s.Cards) {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		val := s.Cards[v]
		fmt.Printf(`Данные строки:
		Имя: %s Номер карты: %s Примечание %s
		`, val.Name, val.CardNumber, val.Comment)
	case 3:
		if v >= len(s.Texts) {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		val := s.Texts[v]
		fmt.Printf(`Данные строки:
		Имя: %s текст: %s Примечание %s
		`, val.Name, val.Data, val.Comment)
	case 4:
		if v >= len(s.Binaries) {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		val := s.Texts[v]
		fmt.Printf(`Данные строки:
		Имя: %s размер данных: %d символов Примечание %s
		`, val.Name, len(val.Data), val.Comment)
	}
}

// AddUserData метод добавляет новую запись в соответствующем разделе.
func (s *UserStorage) AddUserData(i int) {
	switch i {
	case 1:
		var pass = Password{}
		fmt.Print("Введите имя записи: ")
		fmt.Scanln(&pass.Name)
		fmt.Print("Введите логин: ")
		fmt.Scanln(&pass.Login)
		fmt.Print("Введите пароль: ")
		fmt.Scanln(&pass.Pass)
		fmt.Print("Введите примечание: ")
		fmt.Scanln(&pass.Comment)
		s.Passwords = append(s.Passwords, pass)
		fmt.Println("Данные успешно добавлены")
	case 2:
		var card = Card{}
		fmt.Print("Введите имя записи: ")
		fmt.Scanln(&card.Name)
		for {
			fmt.Print("Введите номер карты: ")
			fmt.Scanln(&card.CardNumber)
			if crypto.LynnCheckOrder([]byte(card.CardNumber)) {
				break
			}
			fmt.Println("Номер карты содержит ошибку. Попробуйте еще раз")
		}
		fmt.Print("Введите примечание: ")
		fmt.Scanln(&card.Comment)
		s.Cards = append(s.Cards, card)
		fmt.Println("Данные успешно добавлены")
	case 3:
		var text = Text{}
		fmt.Print("Введите имя записи: ")
		fmt.Scanln(&text.Name)
		fmt.Print("Введите запись: ")
		fmt.Scanln(&text.Data)
		fmt.Print("Введите примечание: ")
		fmt.Scanln(&text.Comment)
		s.Texts = append(s.Texts, text)
		fmt.Println("Данные успешно добавлены")
	case 4:
		var file string
		fmt.Print("Введите путь к файлу, размер файла не должен превышать 64kB: ")
		fmt.Scanln(&file)
		var binary = Binary{}
		var err error
		binary.Name, binary.Data, err = readUserFile(file)
		if errors.Is(err, gkerrors.ErrTooBig) {
			fmt.Println("Размер файла превышает допустимые 64кБ")
			return
		}
		if err != nil {
			fmt.Println("Ошибка чтения файла")
			return
		}
		fmt.Print("Введите примечание: ")
		fmt.Scanln(&binary.Comment)
		s.Binaries = append(s.Binaries, binary)
		fmt.Println("Данные успешно добавлены")
	}
}

// readFile метод считывает данные из файла
func readUserFile(path string) (string, []byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return "", nil, err
	}
	if fi.Size() > 65536 {
		return "", nil, gkerrors.ErrTooBig
	}
	defer file.Close()
	var fileBZ = make([]byte, fi.Size())
	_, err = file.Read(fileBZ)
	if err != nil {
		log.Error().Err(err).Msg("readFile reading file err")
		return "", nil, err
	}
	return fi.Name(), fileBZ, nil
}

// EditUserData метод редактирует существующую запись в соответствующем разделе.
func (s *UserStorage) EditUserData(i, v int) {
	switch i {
	case 1:
		if v >= len(s.Passwords) {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		fmt.Print("Введите имя записи: ")
		fmt.Scanln(&s.Passwords[v].Name)
		fmt.Print("Введите логин: ")
		fmt.Scanln(&s.Passwords[v].Login)
		fmt.Print("Введите пароль: ")
		fmt.Scanln(&s.Passwords[v].Pass)
		fmt.Print("Введите примечание: ")
		fmt.Scanln(&s.Passwords[v].Comment)
		fmt.Println("Данные успешно изменены")
	case 2:
		if v >= len(s.Cards) {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		fmt.Print("Введите имя записи: ")
		fmt.Scanln(&s.Cards[v].Name)
		for {
			fmt.Print("Введите номер карты: ")
			fmt.Scanln(&s.Cards[v].CardNumber)
			if crypto.LynnCheckOrder([]byte(s.Cards[v].CardNumber)) {

				break
			}
			fmt.Println("Номер карты содержит ошибку. Попробуйте еще раз")
		}
		fmt.Print("Введите примечание: ")
		fmt.Scanln(&s.Cards[v].Comment)
		fmt.Println("Данные успешно изменены")
	case 3:
		if v >= len(s.Texts) {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		fmt.Print("Введите имя записи: ")
		fmt.Scanln(&s.Texts[v].Name)
		fmt.Print("Введите запись: ")
		fmt.Scanln(&s.Texts[v].Data)
		fmt.Print("Введите примечание: ")
		fmt.Scanln(&s.Texts[v].Comment)
		fmt.Println("Данные успешно изменены")
	case 4:
		if v >= len(s.Binaries) {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		var file string
		fmt.Print("Введите путь к файлу, размер файла не должен превышать 64kB: ")
		fmt.Scanln(&file)
		var binary = Binary{}
		var err error
		binary.Name, binary.Data, err = readUserFile(file)
		if errors.Is(err, gkerrors.ErrTooBig) {
			fmt.Println("Размер файла превышает допустимые 64кБ")
			return
		}
		if err != nil {
			fmt.Println("Ошибка чтения файла")
			return
		}
		fmt.Print("Введите примечание: ")
		fmt.Scanln(&binary.Comment)
		s.Binaries[v] = binary
		fmt.Println("Данные успешно изменены")
	}
}
