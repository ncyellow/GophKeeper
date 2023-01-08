// Package api модуль реализует client для взаимодействия с сервером. Представлены для варианта grpc и https.
package api

import (
	"github.com/ncyellow/GophKeeper/internal/models"
)

// Sender всякий клиент для работы с сервером должен реализовывать такой интерфейс
type Sender interface {
	// Register запрос на сервер по регистрации клиента
	Register(login string, pwd string) error
	// SignIn запрос на сервер по авторизации клиента
	SignIn(login string, pwd string) error

	// AddCard запрос добавления новой карты
	AddCard(card *models.Card) error
	// Card запрос на чтение уже существующей по ид карты
	Card(cardID string) (*models.Card, error)
	// DelCard запрос на удаление уже существующей по ид карты
	DelCard(cardID string) error

	// AddLogin запрос добавления нового логина
	AddLogin(login *models.Login) error
	// Login запрос на чтение уже существующего логина по ид
	Login(loginID string) (*models.Login, error)
	// DelLogin запрос на удаление уже существующего логина по ид
	DelLogin(loginID string) error

	// AddText запрос на чтение уже существующего текста по ид
	AddText(text *models.Text) error
	// Text запрос на чтение уже существующего текста по ид
	Text(textID string) (*models.Text, error)
	// DelText запрос на удаление уже существующего текста по ид
	DelText(textID string) error

	// AddBin запрос на добавление уже существующих бинарных данных по ид
	AddBin(binary *models.Binary) error
	// Bin запрос на чтение уже существующего бинарных данных по ид
	Bin(binID string) (*models.Binary, error)
	// DelBin запрос на удаление уже существующих бинарных данных по ид
	DelBin(binID string) error
}
