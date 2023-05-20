// Модуль предназначен для формирования CLI-меню на клиенте для взаимодействия с пользователем.
package menu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gophkeeper/internal/client/sender"
)

// EnteringMenu функция запрашивает сессию перед началом аутентификации пользователя
func EnteringMenu(sndr sender.GophKeeperClient) bool {
	err := sndr.ReqSessionID()
	if err != nil {
		log.Error().Err(err).Msg("Get sessionID error")
		fmt.Println("Соединение с сервером не установлено, попробовать еще раз?")
		for {
			var act string
			fmt.Print("Введите команду Yes или No: ")
			fmt.Scanln(&act)
			switch act {
			case "Y", "y", "Yes", "yes":
				err = sndr.ReqSessionID()
				if err == nil {
					return true
				}
				log.Error().Err(err).Msg("Get sessionID error")
				fmt.Println("Соединение с сервером не установлено, попробовать еще раз?")
			case "N", "n", "Q", "q", "No", "no":
				return false
			default:
				fmt.Println("Команда не распознана")
			}
		}
	}
	return true
}

// AuthMenu функция раздела начального меню для выбора авторизации или регистрации пользователя
func AuthMenu(sndr sender.GophKeeperClient) bool {
	fmt.Println("Соединение с сервером установлено")
	for {
		fmt.Println("Для регистрации нового пользователя нажмите R, для входа существующего нажмите L, для завершения работы нажмите Q")
		var act string
		fmt.Scanln(&act)
		switch act {
		case "R", "r":
			return register(sndr)
		case "L", "l":
			return login(sndr)
		case "Q", "q":
			return false
		default:
			fmt.Println("Команда не распознана")
		}
	}
	// return true
}

// register функция меню для создания нового пользователя
func register(sndr sender.GophKeeperClient) bool {
	for {
		login, pass := enterLogin()
		err := sndr.RegisterUser(login, pass)
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.AlreadyExists {
				log.Info().Err(err).Msg("RegisterUser  error")
				fmt.Println("Такой логин уже занят")
				continue
			}
			if st.Code() == codes.Unauthenticated {
				log.Info().Err(err).Msg("RegisterUser sign error")
				fmt.Println("Ошибка проверки подписи. Перезапустить сессию?")
				for {
					var act string
					fmt.Print("Введите команду Yes или No: ")
					fmt.Scanln(&act)
					switch act {
					case "Y", "y", "Yes", "yes":
						err = sndr.RefreshToken()
						if err != nil {
							log.Error().Err(err).Msg("RefreshToken error")
							fmt.Println("Ошибка обновления сессии")
							return false
						}
						err = sndr.ReqSessionID()
						if err == nil {
							break
						}
						log.Error().Err(err).Msg("Get sessionID error")
						fmt.Println("Соединение с сервером не установлено, попробовать еще раз?")
					case "N", "n", "Q", "q", "No", "no":
						return false
					default:
						fmt.Println("Команда не распознана")
					}
				}
			}
		}
		if err != nil {
			log.Error().Err(err).Msg("RegisterUser error")
			fmt.Println("Неизвестная ошибка при регистрации нового пользователя")
			return false
		}
		break
	}
	return mainMenu(sndr)
}

// login функция меню для авторизации существующего пользователя
func login(sndr sender.GophKeeperClient) bool {
	login, pass := enterLogin()
	err := sndr.UserLogin(login, pass)
	st, ok := status.FromError(err)
	if ok {
		if st.Code() == codes.NotFound {
			fmt.Println("Пользователя с таким логином не существует. Попробуйте ввести данные еще раз")
			return true
		}
		if st.Code() == codes.PermissionDenied {
			fmt.Println("Неверный пароль. Попробуйте ввести данные еще раз")
			return true
		}
		if st.Code() == codes.Unauthenticated {
			log.Info().Err(err).Msg("UserLogin sign error")
			fmt.Println("Ошибка проверки подписи. Перезапустить сессию?")
		loop:
			for {
				var act string
				fmt.Print("Введите команду Yes или No: ")
				fmt.Scanln(&act)
				switch act {
				case "Y", "y", "Yes", "yes":
					err = sndr.RefreshToken()
					if err != nil {
						log.Error().Err(err).Msg("RefreshToken error")
						fmt.Println("Ошибка обновления сессии")
						return false
					}
					err = sndr.ReqSessionID()
					if err == nil {
						break loop
					}
					log.Error().Err(err).Msg("Get sessionID error")
					fmt.Println("Соединение с сервером не установлено, попробовать еще раз?")
				case "N", "n", "Q", "q", "No", "no":
					return false
				default:
					fmt.Println("Команда не распознана")
				}
			}
		}
	}

	if err != nil {
		log.Error().Err(err).Msg("UserLogin error")
		fmt.Println("Неизвестная ошибка при авторизации пользователя")
		return false
	}
	return mainMenu(sndr)
}

// mainMenu функция основного меню клиента
func mainMenu(sndr sender.GophKeeperClient) bool {
	fmt.Println("Добро пожаловать в основное меню")
	for {
		var act string
		fmt.Println(`Введите команду:
		C - Проверить актуальный статус данных
		D - скачать пользовательские данные;
		S - сохранить данные на сервер;
		V - посмотреть пользовательские данные
		E - отредактировать или добавить новые данные;
		U - изменить пароль;
		L - разлогиниться;
		Q - завершить работу`)
		fmt.Scanln(&act)
		switch act {
		case "C", "c":
			err := sndr.CheckTimeStamp()
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unauthenticated {
				fmt.Println("Ошибка проверки подписи или время сессии истекло. Попробуйте перелогиниться")
				continue
			}
			if err != nil {
				fmt.Println("Произошла ошибка в процессе запроса")
				log.Error().Err(err).Msg("CheckTimeStamp error")
				continue
			}
		case "D", "d":
			err := sndr.Download()
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unauthenticated {
				fmt.Println("Ошибка проверки подписи или время сессии истекло. Попробуйте перелогиниться")
				continue
			}
			if err != nil {
				fmt.Println("Произошла ошибка в процессе получения данных")
				log.Error().Err(err).Msg("Download error")
				continue
			}
		case "S", "s":
			err := sndr.SaveData()
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unauthenticated {
				fmt.Println("Ошибка проверки подписи или время сессии истекло. Попробуйте перелогиниться")
				continue
			}
			if err != nil {
				fmt.Println("Произошла ошибка в процессе сохранения данных")
				log.Error().Err(err).Msg("SaveData error")
				continue
			}
			fmt.Println("Данные успешно сохранены на сервере")
		case "V", "v":
			viewData(sndr)
		case "E", "e":
			editData(sndr)
		case "U", "u":
			editPassword(sndr)
		case "L", "l":
			sndr.UserLogOut()
			return true
		case "Q", "q":
			sndr.UserLogOut()
			return false
		default:
			fmt.Println("Команда не распознана")
		}
	}
}

// enterLogin функция меню ввода и проверки имени пользователя и пароля
func enterLogin() (string, string) {
	var login, pass string
	for {
		fmt.Println("Введите имя пользователя. Имя не должно содержать символов \",\", \";\"")
		fmt.Scanln(&login)
		if login == "" {
			fmt.Println("Имя пользователя не может быть пустым")
			continue
		}
		if strings.ContainsAny(login, ",;") {
			fmt.Println("Имя пользователя содержит недопустимый символ")
			continue
		}
		break
	}
	for {
		fmt.Println("Введите пароль")
		fmt.Scanln(&pass)
		if pass == "" {
			fmt.Println("Пароль не может быть пустым")
			continue
		}
		break
	}
	return login, pass
}

// viewData функция меню для просмотра существующих данных клиента
func viewData(sndr sender.GophKeeperClient) {
	sndr.Strg.ListUserData()
	for {
		var act string
		fmt.Println(`Введите команду для просмотра детальной информации:
			P - пароли;
			C - карты;
			T - тексты;
			B - бинарные данные;
			R - вернуться в предыдущее меню.`)
		fmt.Scanln(&act)
		switch act {
		case "P", "p":
			sndr.Strg.SliceUserData(1)
		loop_P:
			for {
				fmt.Println("Для просмотра детальной информации введите номер записи, или введите 0 для возврата в предыдущее меню")
				var act string
				fmt.Scanln(&act)
				switch act {
				case "0":
					break loop_P
				default:
					i, err := strconv.Atoi(act)
					if err != nil {
						fmt.Println("Команда не распознана")
						break
					}
					sndr.Strg.StringUserData(1, i-1)
				}
			}
		case "C", "c":
			sndr.Strg.SliceUserData(2)
		loop_C:
			for {
				fmt.Println("Для просмотра детальной информации введите номер записи, или введите 0 для возврата в предыдущее меню")
				var act string
				fmt.Scanln(&act)
				switch act {
				case "0":
					break loop_C
				default:
					i, err := strconv.Atoi(act)
					if err != nil {
						fmt.Println("Команда не распознана")
						break
					}
					sndr.Strg.StringUserData(2, i-1)
				}
			}
		case "T", "t":
			sndr.Strg.SliceUserData(3)
		loop_T:
			for {
				fmt.Println("Для просмотра детальной информации введите номер записи, или введите 0 для возврата в предыдущее меню")
				var act string
				fmt.Scanln(&act)
				switch act {
				case "0":
					break loop_T
				default:
					i, err := strconv.Atoi(act)
					if err != nil {
						fmt.Println("Команда не распознана")
						break
					}
					sndr.Strg.StringUserData(3, i-1)
				}
			}
		case "B", "b":
			sndr.Strg.SliceUserData(4)
		loop_B:
			for {
				fmt.Println("Для просмотра детальной информации введите номер записи, или введите 0 для возврата в предыдущее меню")
				var act string
				fmt.Scanln(&act)
				switch act {
				case "0":
					break loop_B
				default:
					i, err := strconv.Atoi(act)
					if err != nil {
						fmt.Println("Команда не распознана")
						break
					}
					sndr.Strg.StringUserData(4, i-1)
				}
			}
		case "R", "r":
			return
		default:
			fmt.Println("Команда не распознана")
		}
	}
}

// editData функция меню для редактирования и создания новых записей пользователя
func editData(sndr sender.GophKeeperClient) {
	err := sndr.LockUserData()
	if err != nil {
		fmt.Println("Ошибка блокировки данных. Возвращаю в предыдущее меню")
		return
	}
	sndr.Strg.ListUserData()
	for {
		var act string
		fmt.Println(`Введите команду для выбора раздела:
			P - пароли;
			C - карты;
			T - тексты;
			B - бинарные данные;
			R - вернуться в предыдущее меню.`)
		fmt.Scanln(&act)
		switch act {
		case "P", "p":
			sndr.Strg.SliceUserData(1)
		loop_P:
			for {
				fmt.Println("для добавления новой записи введите N\nДля редактирования введите номер записи, или введите 0 для возврата в предыдущее меню")
				var act string
				fmt.Scanln(&act)
				switch act {
				case "N", "n":
					sndr.Strg.AddUserData(1)
				case "0":
					break loop_P
				default:
					i, err := strconv.Atoi(act)
					if err != nil {
						fmt.Println("Команда не распознана")
						break
					}
					sndr.Strg.EditUserData(1, i-1)
				}
			}
		case "C", "c":
			sndr.Strg.SliceUserData(2)
		loop_N:
			for {
				fmt.Println("для добавления новой записи введите N\nДля редактирования введите номер записи, или введите 0 для возврата в предыдущее меню")
				var act string
				fmt.Scanln(&act)
				switch act {
				case "N", "n":
					sndr.Strg.AddUserData(2)
				case "0":
					break loop_N
				default:
					i, err := strconv.Atoi(act)
					if err != nil {
						fmt.Println("Команда не распознана")
						break
					}
					sndr.Strg.EditUserData(2, i-1)
				}
			}
		case "T", "t":
			sndr.Strg.SliceUserData(3)
		loop_T:
			for {
				fmt.Println("для добавления новой записи введите N\nДля редактирования введите номер записи, или введите 0 для возврата в предыдущее меню")
				var act string
				fmt.Scanln(&act)
				switch act {
				case "N", "n":
					sndr.Strg.AddUserData(3)
				case "0":
					break loop_T
				default:
					i, err := strconv.Atoi(act)
					if err != nil {
						fmt.Println("Команда не распознана")
						break
					}
					sndr.Strg.EditUserData(3, i-1)
				}
			}
		case "B", "b":
			sndr.Strg.SliceUserData(4)
		loop_B:
			for {
				fmt.Println("для добавления новой записи введите N\nДля редактирования введите номер записи, или введите 0 для возврата в предыдущее меню")
				var act string
				fmt.Scanln(&act)
				switch act {
				case "N", "n":
					sndr.Strg.AddUserData(4)
				case "0":
					break loop_B
				default:
					i, err := strconv.Atoi(act)
					if err != nil {
						fmt.Println("Команда не распознана")
						break
					}
					sndr.Strg.EditUserData(4, i-1)
				}
			}
		case "R", "r":
			return
		default:
			fmt.Println("Команда не распознана")
		}
	}
}

func editPassword(sndr sender.GophKeeperClient) {
	var oldPassword, newPassword, checkPassword string
	for {
		fmt.Println("Введите текущий пароль")
		fmt.Scanln(&oldPassword)
		if oldPassword == "" {
			fmt.Println("Пароль не может быть пустым")
			continue
		}
		break
	}
	for {
		fmt.Println("Введите новый пароль")
		fmt.Scanln(&newPassword)
		if newPassword == "" {
			fmt.Println("Пароль не может быть пустым")
			continue
		}
		fmt.Println("Повторите новый пароль")
		fmt.Scanln(&checkPassword)
		if checkPassword == "" {
			fmt.Println("Пароль не может быть пустым")
			continue
		}
		if newPassword != checkPassword {
			fmt.Println("Новый пароль и его повторение не совпадают")
			continue
		}
		break
	}
	true, err := sndr.ChangePassword(oldPassword, newPassword)
	st, ok := status.FromError(err)
	if ok {
		if st.Code() == codes.InvalidArgument {
			fmt.Println("Неверный пароль")
			return
		}
		if st.Code() == codes.Unauthenticated {
			fmt.Println("Ошибка проверки подписи или время сессии истекло. Попробуйте перелогиниться")
			return
		}
	}
	if err != nil {
		log.Error().Err(err).Msg("ChangePassword error")
		fmt.Println("Произошла ошибка при попытке изменения пароля")
		return
	}
	if true {
		fmt.Println("Пароль успешно изменен")
		return
	}
	fmt.Println("Не удалось измененить пароль")
}
