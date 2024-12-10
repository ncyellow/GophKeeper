// Package console implements the module for parsing console input
package console

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/ncyellow/GophKeeper/internal/client/api"
	"github.com/ncyellow/GophKeeper/internal/client/config"
)

// LivePrefixState auxiliary structure to create a nice prompt with the name of the authorized user
var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

// changeLivePrefix special function for working with the input prompt using the go-prompt library
func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}

// Console structure for handling keyboard input. Using go-prompt
type Console struct {
	Conf   *config.Config
	Client api.Sender
}

// CreateExecutor function for processing all commands entered from the keyboard
func CreateExecutor(sender api.Sender) func(string) {
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
					fmt.Printf("Card with ID - %s successfully added\n", card.ID)
				}
			}
		case "card":
			if len(commands) != 2 {
				fmt.Println("Enter card number!")
			} else {
				card, err := sender.Card(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Card data read successfully - %#v!\n", card)
				}
			}
		case "card-del":
			if len(commands) != 2 {
				fmt.Println("Enter card identifier!")
			} else {
				err := sender.DelCard(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Card deleted!")
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
					fmt.Printf("Login with ID - %s successfully added\n", login.ID)
				}
			}
		case "login":
			if len(commands) != 2 {
				fmt.Println("Enter login identifier!")
			} else {
				login, err := sender.Login(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Login data read successfully - %#v!\n", login)
				}
			}
		case "login-del":
			if len(commands) != 2 {
				fmt.Println("Enter login identifier!")
			} else {
				err := sender.DelLogin(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Login deleted!")
				}
			}
		case "text-add":
			text, err := readText()
			if err == nil {
				err := sender.AddText(text)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Text with ID - %s successfully added\n", text.ID)
				}
			}
		case "text":
			if len(commands) != 2 {
				fmt.Println("Enter text identifier!")
			} else {
				text, err := sender.Text(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Text data read successfully - %#v!\n", text)
				}
			}
		case "text-del":
			if len(commands) != 2 {
				fmt.Println("Enter text identifier!")
			} else {
				err := sender.DelText(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Text deleted!")
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
					fmt.Printf("Bin with ID - %s successfully added\n", bin.ID)
				}
			}
		case "bin":
			if len(commands) != 2 {
				fmt.Println("Enter file identifier!")
			} else {
				bin, err := sender.Bin(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Bin data read successfully - %#v!\n", bin)
				}
			}
		case "bin-del":
			if len(commands) != 2 {
				fmt.Println("Enter file identifier!")
			} else {
				err := sender.DelBin(commands[1])
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Bin deleted!")
				}
			}
		}
	}
}

// completer - implementation of autocompletion
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

// Run starts our own console with a prompt
func (p *Console) Run() {
	executor := CreateExecutor(p.Client)
	prom := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>>"),
		prompt.OptionLivePrefix(changeLivePrefix),
	)
	prom.Run()
}
