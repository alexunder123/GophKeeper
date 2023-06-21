package menu

import (
	"errors"
	"fmt"
	"os"

	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/client/sender"
	"gophkeeper/internal/client/storage"
	gkerrors "gophkeeper/internal/errors"

	"github.com/rs/zerolog/log"
)

// dataType - определяем тип данных пользователя для работы с меню.
type dataType int

// Определяем константы для выбора типов данных.
const (
	usersPasswords dataType = iota
	usersCards
	usersTexts
	usersBinaries
)

// printSliceUserData метод выводит на экран информацию о сохраненных строках пользователя в соответствующем разделе.
func printSliceUserData(i dataType, sndr sender.GophKeeperClient) {
	switch i {
	case usersPasswords:
		passwords := sndr.Strg.SliceUsersPasswords()
		for i, val := range passwords {
			fmt.Printf("Номер: %d, Имя: %s\n", i+1, val.Name)
		}
	case usersCards:
		cards := sndr.Strg.SliceUsersCards()
		for i, val := range cards {
			fmt.Printf("Номер: %d, Имя: %s\n", i+1, val.Name)
		}
	case usersTexts:
		texts := sndr.Strg.SliceUsersTexts()
		for i, val := range texts {
			fmt.Printf("Номер: %d, Имя: %s\n", i+1, val.Name)
		}
	case usersBinaries:
		binaries := sndr.Strg.SliceUsersBinaries()
		for i, val := range binaries {
			fmt.Printf("Номер: %d, Имя: %s\n", i+1, val.Name)
		}
	}
}

// printStringUserData метод выводит на экран полную информацию о сохраненной записи.
func printStringUserData(i dataType, number int, sndr sender.GophKeeperClient) {
	switch i {
	case usersPasswords:
		val := sndr.Strg.StringUsersPassword(number)
		if val == nil {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		fmt.Printf(`Данные строки:
		Имя: %s Логин: %s Пароль: %s Примечание: %s
		`, val.Name, val.Login, val.Pass, val.Comment)
	case usersCards:
		val := sndr.Strg.StringUsersCard(number)
		if val == nil {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		fmt.Printf(`Данные строки:
		Имя: %s Номер карты: %s Примечание %s
		`, val.Name, val.CardNumber, val.Comment)
	case usersTexts:
		val := sndr.Strg.StringUsersText(number)
		if val == nil {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		fmt.Printf(`Данные строки:
		Имя: %s текст: %s Примечание %s
		`, val.Name, val.Data, val.Comment)
	case usersBinaries:
		val := sndr.Strg.StringUsersBinary(number)
		if val == nil {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		fmt.Printf(`Данные строки:
		Имя: %s размер данных: %d символов Примечание %s
		`, val.Name, len(val.Data), val.Comment)
	}
}

// AddUserData метод добавляет новую запись в соответствующем разделе.
func addUsersData(i dataType, sndr sender.GophKeeperClient) {
	switch i {
	case usersPasswords:
		var pass = storage.Password{}
		inputUsersPasswords(&pass)
		sndr.Strg.AddUsersPassword(&pass)
		fmt.Println("Данные успешно добавлены")
	case usersCards:
		var card = storage.Card{}
		inputUsersCards(&card)
		sndr.Strg.AddUsersCard(&card)
		fmt.Println("Данные успешно добавлены")
	case usersTexts:
		var text = storage.Text{}
		inputUsersTexts(&text)
		sndr.Strg.AddUsersText(&text)
		fmt.Println("Данные успешно добавлены")
	case usersBinaries:
		var binary = storage.Binary{}
		inputUsersBinaries(&binary)
		sndr.Strg.AddUsersBinary(&binary)
		fmt.Println("Данные успешно добавлены")
	}
}

// editUsersData метод редактирует существующую запись в соответствующем разделе.
func editUsersData(i dataType, v int, sndr sender.GophKeeperClient) {
	switch i {
	case usersPasswords:
		val := sndr.Strg.StringUsersPassword(v)
		if val == nil {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		var pass = storage.Password{}
		inputUsersPasswords(&pass)
		sndr.Strg.EditUsersPassword(v, &pass)
		fmt.Println("Данные успешно изменены")
	case usersCards:
		val := sndr.Strg.StringUsersCard(v)
		if val == nil {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		var card = storage.Card{}
		inputUsersCards(&card)
		sndr.Strg.EditUsersCard(v, &card)
		fmt.Println("Данные успешно изменены")
	case usersTexts:
		val := sndr.Strg.StringUsersText(v)
		if val == nil {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		var text = storage.Text{}
		inputUsersTexts(&text)
		sndr.Strg.EditUsersText(v, &text)
		fmt.Println("Данные успешно изменены")
	case usersBinaries:
		val := sndr.Strg.StringUsersBinary(v)
		if val == nil {
			fmt.Println("В базе нет строки с таким номером!")
			return
		}
		var binary = storage.Binary{}
		inputUsersBinaries(&binary)
		sndr.Strg.EditUsersBinary(v, &binary)
		fmt.Println("Данные успешно изменены")
	}
}

// inputUsersPasswords метод взаимодействует с пользователем для ввода данных в записи паролей.
func inputUsersPasswords(pass *storage.Password) {
	fmt.Print("Введите имя записи: ")
	fmt.Scanln(&pass.Name)
	fmt.Print("Введите логин: ")
	fmt.Scanln(&pass.Login)
	fmt.Print("Введите пароль: ")
	fmt.Scanln(&pass.Pass)
	fmt.Print("Введите примечание: ")
	fmt.Scanln(&pass.Comment)
}

// inputUsersCards метод взаимодействует с пользователем для ввода данных в записи карт.
func inputUsersCards(card *storage.Card) {
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
}

// inputUsersTexts метод взаимодействует с пользователем для ввода данных в записи текстов.
func inputUsersTexts(text *storage.Text) {
	fmt.Print("Введите имя записи: ")
	fmt.Scanln(&text.Name)
	fmt.Print("Введите запись: ")
	fmt.Scanln(&text.Data)
	fmt.Print("Введите примечание: ")
	fmt.Scanln(&text.Comment)
}

// inputUsersBinaries метод взаимодействует с пользователем для ввода данных в записи двоичных данных.
func inputUsersBinaries(binary *storage.Binary) {
	var file string
	fmt.Print("Введите путь к файлу, размер файла не должен превышать 64kB: ")
	fmt.Scanln(&file)
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
