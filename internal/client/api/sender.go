// Package api implements a client for interacting with the server. It provides both grpc and https variants.
package api

import (
	"github.com/ncyellow/GophKeeper/internal/models"
)

// Sender any client for working with the server should implement this interface
type Sender interface {
	// Register request to the server for client registration
	Register(login string, pwd string) error
	// SignIn request to the server for client authorization
	SignIn(login string, pwd string) error

	// AddCard request to add a new card
	AddCard(card *models.Card) error
	// Card request to read an existing card by id
	Card(cardID string) (*models.Card, error)
	// DelCard request to delete an existing card by id
	DelCard(cardID string) error

	// AddLogin request to add a new login
	AddLogin(login *models.Login) error
	// Login request to read an existing login by id
	Login(loginID string) (*models.Login, error)
	// DelLogin request to delete an existing login by id
	DelLogin(loginID string) error

	// AddText request to add new text content
	AddText(text *models.Text) error
	// Text request to read existing text content by id
	Text(textID string) (*models.Text, error)
	// DelText request to delete existing text content by id
	DelText(textID string) error

	// AddBin request to add new binary data
	AddBin(binary *models.Binary) error
	// Bin request to read existing binary data by id
	Bin(binID string) (*models.Binary, error)
	// DelBin request to delete existing binary data by id
	DelBin(binID string) error
}
