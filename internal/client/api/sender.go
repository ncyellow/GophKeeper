package api

import (
	"github.com/ncyellow/GophKeeper/internal/models"
)

type Sender interface {
	Register(login string, pwd string) error
	SignIn(login string, pwd string) error

	AddCard(card *models.Card) error
	Card(cardID string) (*models.Card, error)
	DelCard(cardID string) error

	AddLogin(login *models.Login) error
	Login(loginID string) (*models.Login, error)
	DelLogin(loginID string) error

	AddText(text *models.Text) error
	Text(textID string) (*models.Text, error)
	DelText(textID string) error

	AddBin(binary *models.Binary) error
	Bin(binID string) (*models.Binary, error)
	DelBin(binID string) error
}
