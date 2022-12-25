// Package console реализует модуль для парсинга ввода консоли
package console

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/ncyellow/GophKeeper/internal/client/api"
	"github.com/ncyellow/GophKeeper/internal/client/config"
)

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}

// Console структура по обработке ввода с клавиатуры. С использованием go-prompt
type Console struct {
	Conf *config.Config
}

// CreateExecutor функция обработки всех введенных с клавиатуры команд
func CreateExecutor(conf *config.Config) func(string) {
	sender := api.NewHTTPSender(conf)
	return func(t string) {
		s := strings.TrimSpace(t)
		commands := strings.Split(s, " ")
		switch commands[0] {
		case "exit":
			fmt.Println("bye!")
			os.Exit(0)
		case "version":
			fmt.Printf("Build version: %s\n", config.BuildVersion)
			fmt.Printf("Build date: %s\n", config.BuildDate)
		case "register":
			username, password, err := credentials()
			if err == nil {
				err := sender.Register(username, password)
				if err != nil {
					fmt.Println("")
					fmt.Println(err.Error())
				} else {
					fmt.Println("")
					LivePrefixState.LivePrefix = fmt.Sprintf("%s >>>", username)
					LivePrefixState.IsEnable = true
				}
			}
		case "signin":
			username, password, err := credentials()
			if err == nil {
				err := sender.SignIn(username, password)
				if err != nil {
					fmt.Println("")
					fmt.Println(err.Error())
				} else {
					fmt.Println("")
					LivePrefixState.LivePrefix = fmt.Sprintf("%s >>>", username)
					LivePrefixState.IsEnable = true
				}
			}
		case "card-add":
			card, err := readCard()
			if err == nil {
				err := sender.AddCard(card)
				if err != nil {
					fmt.Println("")
					fmt.Println(err.Error())
				} else {
					fmt.Println("")
					fmt.Printf("Карта с ID - %s успешно добавлена\n", card.ID)
				}
			}
		case "card":
			if len(commands) != 2 {
				fmt.Println("Введите номер карты!")
			} else {
				card, err := sender.Card(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Данные карты прочитаны успешно - %#v!\n", card)
				}
			}
		case "card-del":
			if len(commands) != 2 {
				fmt.Println("Введите идентификатор карты!")
			} else {

				err := sender.DelCard(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Карта удалена!")
				}
			}
		case "login-add":
			login, err := readLogin()
			if err == nil {
				err := sender.AddLogin(login)
				if err != nil {
					fmt.Println("")
					fmt.Println(err.Error())
				} else {
					fmt.Println("")
					fmt.Printf("Логин с ID - %s успешно добавлена\n", login.ID)
				}
			}
		case "login":
			if len(commands) != 2 {
				fmt.Println("Введите идентификатор логина!")
			} else {
				login, err := sender.Login(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Данные логина прочитаны успешно - %#v!\n", login)
				}
			}
		case "login-del":
			if len(commands) != 2 {
				fmt.Println("Введите идентификатор логина!")
			} else {
				err := sender.DelLogin(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Логин удален!")
				}
			}
		case "text-add":
			text, err := readText()
			if err == nil {
				err := sender.AddText(text)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Text с ID - %s успешно добавлена\n", text.ID)
				}
			}
		case "text":
			if len(commands) != 2 {
				fmt.Println("Введите идентификатор text!")
			} else {
				text, err := sender.Text(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Данные text прочитаны успешно - %#v!\n", text)
				}
			}
		case "text-del":
			if len(commands) != 2 {
				fmt.Println("Введите идентификатор text!")
			} else {
				err := sender.DelText(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Text удален!")
				}
			}
		case "bin-add":
			bin, err := readBinary()
			if err == nil {
				err := sender.AddBin(bin)
				if err != nil {
					fmt.Println("")
					fmt.Println(err.Error())
				} else {
					fmt.Println("")
					fmt.Printf("Bin с ID - %s успешно добавлена\n", bin.ID)
				}
			}
		case "bin":
			if len(commands) != 2 {
				fmt.Println("Введите идентификатор файла!")
			} else {
				bin, err := sender.Bin(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Данные bin прочитаны успешно - %#v!\n", bin)
				}
			}
		case "bin-del":
			if len(commands) != 2 {
				fmt.Println("Введите идентификатор файла!")
			} else {
				err := sender.DelBin(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Bin удален!")
				}
			}
		}
	}
}

// completer - реализация автодополнения
func completer(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest
	if d.Text != "" {

		words := strings.Split(d.Text, " ")
		if len(words) == 1 {

			s = []prompt.Suggest{
				{Text: "register", Description: "Create new user"},
				{Text: "signin", Description: "SignIn user"},

				{Text: "card-add", Description: "Add new card"},
				{Text: "card", Description: "Get card data"},
				{Text: "card-del", Description: "Delete card"},

				{Text: "login-add", Description: "Add new login"},
				{Text: "login", Description: "Get login data"},
				{Text: "login-del", Description: "Delete login"},

				{Text: "text-add", Description: "Add new text"},
				{Text: "text", Description: "Get text data"},
				{Text: "text-del", Description: "Delete text"},

				{Text: "bin-add", Description: "Add new binary"},
				{Text: "bin", Description: "Get binary data"},
				{Text: "bin-del", Description: "Delete binary"},

				{Text: "help", Description: "List all available commands"},
				{Text: "version", Description: "Client version"},
				{Text: "exit", Description: "Exit"},
			}
		}
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// Run запуск нашей собственной консоли с приглашением
func (p *Console) Run() {
	executor := CreateExecutor(p.Conf)
	prom := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>>"),
		prompt.OptionLivePrefix(changeLivePrefix),
	)
	prom.Run()
}
