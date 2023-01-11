// Package storage - реализует чтение запись в базу всех основных сущностей, Карты, Логины, и тд.
package storage

import (
	"context"
	"errors"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ncyellow/GophKeeper/internal/models"
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/rs/zerolog/log"
)

var ErrPgConnect = errors.New("cant connect to pgsql")

type PgStorage struct {
	conf *config.Config
	pool pgxpoolmock.PgxPool
}

// NewPgStorage конструктор хранилища на основе postgresql, явно не используется, только через фабрику
func NewPgStorage(conf *config.Config) (*PgStorage, error) {
	pool, err := pgxpool.Connect(context.Background(), conf.DatabaseConn)
	if err != nil {
		return nil, ErrPgConnect
	}

	store := PgStorage{
		conf: conf,
		pool: pool,
	}
	return &store, nil
}

func (p *PgStorage) Close() {
	p.pool.Close()
}

func (p *PgStorage) Register(ctx context.Context, user models.User) (int64, error) {
	var lastInsertID int64 = 0
	err := p.pool.QueryRow(ctx, `
	INSERT INTO "users"("login", "password")
	VALUES ($1, $2)
	returning "@users"`, user.Login, user.Password).Scan(&lastInsertID)
	// так как логин у нас уникален, то при попытке вставить второй одинаковый логин будет ошибка
	if err != nil {
		return lastInsertID, err
	}
	return lastInsertID, nil
}

func (p *PgStorage) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	row := p.pool.QueryRow(ctx, `
	SELECT "@users", "login" FROM "users" WHERE "login" = $1
	LIMIT 1
	`, login)

	err := row.Scan(&user.UserID, &user.Login)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PgStorage) User(ctx context.Context, login string, password string) (*models.User, error) {
	var user models.User
	row := p.pool.QueryRow(ctx, `
	SELECT "@users", "login", "password" FROM "users" WHERE "login" = $1 AND "password" = $2
	LIMIT 1
	`, login, password)

	err := row.Scan(&user.UserID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PgStorage) AddCard(ctx context.Context, userID int64, card models.Card) error {
	var lastInsertID int64 = 0
	err := p.pool.QueryRow(ctx, `
	INSERT INTO "cards"("id", "user", "fio", "number", "date", "cvv", "metainfo")
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	returning "@cards"
	`, card.ID, userID, card.FIO, card.Number, card.Date, card.CVV, card.MetaInfo).Scan(&lastInsertID)
	// так как логин у нас уникален, то при попытке вставить второй одинаковый логин будет ошибка
	if err != nil {
		return err
	}
	log.Info().Msgf("Новая карта с идентификатором - %s добавлена. @card - %d", card.ID, lastInsertID)
	return nil
}

func (p *PgStorage) Card(ctx context.Context, userID int64, cardID string) (*models.Card, error) {
	var card models.Card

	err := p.pool.QueryRow(ctx, `
	SELECT "id", "user", "fio", "number", "date", "cvv", "metainfo"
	FROM "cards"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, cardID).Scan(&card.ID, &card.UserID, &card.FIO, &card.Number, &card.Date, &card.CVV, &card.MetaInfo)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (p *PgStorage) DeleteCard(ctx context.Context, userID int64, cardID string) error {
	_, err := p.pool.Exec(ctx, `
	DELETE FROM "cards"
	WHERE "user" = $1 and "id" = $2
	`, userID, cardID)

	return err
}

func (p *PgStorage) AddLogin(ctx context.Context, userID int64, login models.Login) error {
	var lastInsertID int64 = 0
	err := p.pool.QueryRow(ctx, `
	INSERT INTO "logins"("id", "user", "login", "password", "metainfo")
	VALUES ($1, $2, $3, $4, $5)
	returning "@logins"
	`, login.ID, userID, login.Login, login.Password, login.MetaInfo).Scan(&lastInsertID)
	if err != nil {
		return err
	}
	log.Info().Msgf("Новые учетные данные с идентификатором - %s добавлены. @login - %d", login.ID, lastInsertID)
	return nil
}

func (p *PgStorage) Login(ctx context.Context, userID int64, loginID string) (*models.Login, error) {
	var login models.Login

	err := p.pool.QueryRow(ctx, `
	SELECT "id", "user", "login", "password", "metainfo"
	FROM "logins"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, loginID).Scan(&login.ID, &login.UserID, &login.Login, &login.Password, &login.MetaInfo)
	if err != nil {
		return nil, err
	}
	return &login, nil
}

func (p *PgStorage) DeleteLogin(ctx context.Context, userID int64, loginID string) error {
	_, err := p.pool.Exec(ctx, `
	DELETE FROM "logins"
	WHERE "user" = $1 and "id" = $2
	`, userID, loginID)

	return err
}

func (p *PgStorage) AddText(ctx context.Context, userID int64, text models.Text) error {
	var lastInsertID int64 = 0
	err := p.pool.QueryRow(ctx, `
	INSERT INTO "text_data"("id", "user", "content", "metainfo")
	VALUES ($1, $2, $3, $4)
	returning "@text"
	`, text.ID, userID, text.Content, text.MetaInfo).Scan(&lastInsertID)
	if err != nil {
		return err
	}
	log.Info().Msgf("Новые текстовые данные с идентификатором - %s добавлены. @login - %d", text.ID, lastInsertID)
	return nil
}

func (p *PgStorage) Text(ctx context.Context, userID int64, textID string) (*models.Text, error) {
	var text models.Text

	err := p.pool.QueryRow(ctx, `
	SELECT "id", "user", "content", "metainfo"
	FROM "text_data"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, textID).Scan(&text.ID, &text.UserID, &text.Content, &text.MetaInfo)
	if err != nil {
		return nil, err
	}
	return &text, nil
}

func (p *PgStorage) DeleteText(ctx context.Context, userID int64, textID string) error {
	_, err := p.pool.Exec(ctx, `
	DELETE FROM "text_data"
	WHERE "user" = $1 and "id" = $2
	`, userID, textID)

	return err
}

func (p *PgStorage) AddBinary(ctx context.Context, userID int64, binData models.Binary) error {
	var lastInsertID int64 = 0
	err := p.pool.QueryRow(ctx, `
	INSERT INTO "bin_data"("id", "user", "content", "metainfo")
	VALUES ($1, $2, $3, $4)
	returning "@bin"
	`, binData.ID, userID, binData.Data, binData.MetaInfo).Scan(&lastInsertID)
	if err != nil {
		return err
	}
	log.Info().Msgf("Новые бинарные данные с идентификатором - %s добавлены. @bin - %d",
		binData.ID, lastInsertID)
	return nil
}

func (p *PgStorage) Binary(ctx context.Context, userID int64, binID string) (*models.Binary, error) {
	var binary models.Binary

	err := p.pool.QueryRow(ctx, `
	SELECT "id", "user", "content", "metainfo"
	FROM "bin_data"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, binID).Scan(&binary.ID, &binary.UserID, &binary.Data, &binary.MetaInfo)
	if err != nil {
		return nil, err
	}
	return &binary, nil
}

func (p *PgStorage) DeleteBinary(ctx context.Context, userID int64, binID string) error {
	_, err := p.pool.Exec(ctx, `
	DELETE FROM "bin_data"
	WHERE "user" = $1 and "id" = $2
	`, userID, binID)

	return err
}
