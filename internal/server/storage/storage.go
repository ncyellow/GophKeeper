package storage

import (
	"context"

	"github.com/ncyellow/GophKeeper/internal/models"
)

// Storage описываем интерфейс для записи и чтения в бд, аналогичен HTTP API
type Storage interface {
	Register(ctx context.Context, user models.User) (int64, error)
	UserByLogin(ctx context.Context, login string) (*models.User, error)
	User(ctx context.Context, login string, password string) (*models.User, error)

	AddCard(ctx context.Context, userID int64, card models.Card) error
	Card(ctx context.Context, userID int64, cardID string) (*models.Card, error)
	DeleteCard(ctx context.Context, userID int64, cardID string) error

	AddLogin(ctx context.Context, userID int64, login models.Login) error
	Login(ctx context.Context, userID int64, loginID string) (*models.Login, error)
	DeleteLogin(ctx context.Context, userID int64, loginID string) error

	AddText(ctx context.Context, userID int64, text models.Text) error
	Text(ctx context.Context, userID int64, textID string) (*models.Text, error)
	DeleteText(ctx context.Context, userID int64, textID string) error

	AddBinary(ctx context.Context, userID int64, binData models.Binary) error
	Binary(ctx context.Context, userID int64, binID string) (*models.Binary, error)
	DeleteBinary(ctx context.Context, userID int64, binID string) error

	Close()
}
