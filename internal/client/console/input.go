// Package console. This part of the module implements requests from the console for main entities: Card, Login, Text data, etc.
package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/ncyellow/GophKeeper/internal/models"
)

// credentials - reads login and password from the console. If everything is ok, error will be nil
func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

// readCard - reads card data from the console. If everything is ok, error will be nil
func readCard() (*models.Card, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter ID: ")
	cardID, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter FIO: ")
	fio, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter Number: ")
	number, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter Date: ")
	cardDate, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter CVV: ")
	cvv, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter MetaInfo: ")
	metaInfo, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return &models.Card{
		ID:       strings.TrimSpace(cardID),
		FIO:      strings.TrimSpace(fio),
		Number:   strings.TrimSpace(number),
		Date:     strings.TrimSpace(cardDate),
		CVV:      strings.TrimSpace(cvv),
		MetaInfo: strings.TrimSpace(metaInfo),
	}, nil
}

// readLogin - reads login data from the console. If everything is ok, error will be nil
func readLogin() (*models.Login, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter ID: ")
	cardID, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter Login: ")
	login, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	password := string(bytePassword)

	fmt.Print("Enter MetaInfo: ")
	metaInfo, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return &models.Login{
		ID:       strings.TrimSpace(cardID),
		Login:    strings.TrimSpace(login),
		Password: strings.TrimSpace(password),
		MetaInfo: strings.TrimSpace(metaInfo),
	}, nil
}

// readText - reads text data from the console. If everything is ok, error will be nil
func readText() (*models.Text, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter ID: ")
	cardID, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter Content: ")
	content, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter MetaInfo: ")
	metaInfo, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return &models.Text{
		ID:       strings.TrimSpace(cardID),
		Content:  strings.TrimSpace(content),
		MetaInfo: strings.TrimSpace(metaInfo),
	}, nil
}

// readText - reads binary data from the console. If everything is ok, error will be nil
func readBinary() (*models.Binary, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter ID: ")
	cardID, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter Filename: ")
	content, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	content = strings.TrimSpace(content)
	data, err := os.ReadFile(content)
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter MetaInfo: ")
	metaInfo, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return &models.Binary{
		ID:       strings.TrimSpace(cardID),
		Data:     data,
		MetaInfo: strings.TrimSpace(metaInfo),
	}, nil
}
