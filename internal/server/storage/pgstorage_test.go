package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/ncyellow/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// PgStorageSuite - тесты работы с базой
// Используем pgxpoolmock.MockPgxPool для mock
type PgStorageSuite struct {
	suite.Suite
	mockPool *pgxpoolmock.MockPgxPool
	store    Storage
}

// TestPgStorageSuite запуск всех тестов PgStorageSuite
func TestPgStorageSuite(t *testing.T) {
	suite.Run(t, new(PgStorageSuite))
}

// SetupTest инициализация. Создаем. Repository, Storage и мокаем базу
func (suite *PgStorageSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	suite.mockPool = mockPool
	suite.store = &PgStorage{pool: mockPool}
}

func (suite *PgStorageSuite) TestRegister() {
	user := models.User{
		UserID:   1,
		Login:    "login",
		Password: "pwd",
	}

	//! Тест на корректную вставку
	columns := []string{"id"}
	pgxRows := pgxpoolmock.NewRows(columns).AddRow(user.UserID).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "users"("login", "password")
	VALUES ($1, $2)
	returning "@users"`, user.Login, user.Password).Return(pgxRows)

	userID, err := suite.store.Register(context.Background(), user)
	assert.Equal(suite.T(), userID, user.UserID)
	assert.NoError(suite.T(), err)

	//! Тест на вставку с конфликтами
	columns = []string{"id"}
	pgxRows = pgxpoolmock.
		NewRows(columns).
		AddRow(nil).
		RowError(0, pgx.ErrNoRows).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "users"("login", "password")
	VALUES ($1, $2)
	returning "@users"`, user.Login, user.Password).Return(pgxRows)

	userID, err = suite.store.Register(context.Background(), user)
	assert.Equal(suite.T(), userID, int64(0))
	assert.Error(suite.T(), err, pgx.ErrNoRows)
}

func (suite *PgStorageSuite) TestUser() {
	newID := int64(1)
	user := models.User{
		UserID:   newID,
		Login:    "login",
		Password: "pwd",
	}

	//! Тест на корректную вставку
	columns := []string{"@users", "login", "password"}
	pgxRows := pgxpoolmock.NewRows(columns).AddRow(newID, user.Login, user.Password).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "@users", "login", "password" FROM "users" WHERE "login" = $1 AND "password" = $2
	LIMIT 1
	`, user.Login, user.Password).Return(pgxRows)

	resultUser, err := suite.store.User(context.Background(), user.Login, user.Password)
	assert.Equal(suite.T(), *resultUser, user)
	assert.NoError(suite.T(), err)

	//! Тест на вставку с конфликтами
	pgxRows = pgxpoolmock.
		NewRows(columns).
		AddRow(newID, user.Login, user.Password).
		RowError(0, pgx.ErrNoRows).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "@users", "login", "password" FROM "users" WHERE "login" = $1 AND "password" = $2
	LIMIT 1
	`, user.Login, user.Password).Return(pgxRows)

	resultUser, err = suite.store.User(context.Background(), user.Login, user.Password)
	assert.Nil(suite.T(), resultUser)
	assert.Error(suite.T(), err, pgx.ErrNoRows)
}

func (suite *PgStorageSuite) TestUserByLogin() {
	userID := int64(1)
	user := models.User{
		UserID: userID,
		Login:  "login",
	}

	//! Тест на корректную вставку
	columns := []string{"@users", "login"}
	pgxRows := pgxpoolmock.NewRows(columns).AddRow(userID, user.Login).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "@users", "login" FROM "users" WHERE "login" = $1
	LIMIT 1
	`, user.Login).Return(pgxRows)

	resultUser, err := suite.store.UserByLogin(context.Background(), user.Login)
	assert.Equal(suite.T(), *resultUser, user)
	assert.NoError(suite.T(), err)

	//! Тест когда не найден пользователь
	pgxRows = pgxpoolmock.
		NewRows(columns).
		AddRow(userID, user.Login).
		RowError(0, pgx.ErrNoRows).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "@users", "login" FROM "users" WHERE "login" = $1
	LIMIT 1
	`, user.Login).Return(pgxRows)

	resultUser, err = suite.store.UserByLogin(context.Background(), user.Login)
	assert.Nil(suite.T(), resultUser)
	assert.Error(suite.T(), err, pgx.ErrNoRows)
}

func (suite *PgStorageSuite) TestAddCard() {
	userID := int64(1)
	card := models.Card{
		ID:       "testID",
		FIO:      "fio",
		Number:   "number",
		Date:     "date",
		CVV:      "cvv",
		MetaInfo: "metainfo",
	}

	//! Тест на корректную вставку
	columns := []string{"id"}
	pgxRows := pgxpoolmock.NewRows(columns).
		AddRow(int64(1)).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "cards"("id", "user", "fio", "number", "date", "cvv", "metainfo")
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	returning "@cards"
	`, card.ID, userID, card.FIO, card.Number, card.Date, card.CVV, card.MetaInfo).Return(pgxRows)

	err := suite.store.AddCard(context.Background(), userID, card)
	assert.NoError(suite.T(), err)

	//! Тест на ошибки sql
	targetErr := errors.New("some error")
	pgxRows = pgxpoolmock.NewRows(columns).
		AddRow(int64(1)).
		RowError(0, targetErr).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "cards"("id", "user", "fio", "number", "date", "cvv", "metainfo")
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	returning "@cards"
	`, card.ID, userID, card.FIO, card.Number, card.Date, card.CVV, card.MetaInfo).Return(pgxRows)

	err = suite.store.AddCard(context.Background(), userID, card)
	assert.Error(suite.T(), err, targetErr)
}

func (suite *PgStorageSuite) TestCard() {
	userID := int64(1)
	card := models.Card{
		ID:       "testID",
		FIO:      "fio",
		Number:   "number",
		Date:     "date",
		CVV:      "cvv",
		MetaInfo: "metainfo",
		UserID:   userID,
	}

	columns := []string{"id", "user", "fio", "number", "date", "cvv", "metainfo"}
	pgxRows := pgxpoolmock.NewRows(columns).
		AddRow(card.ID, userID, card.FIO, card.Number, card.Date, card.CVV, card.MetaInfo).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "id", "user", "fio", "number", "date", "cvv", "metainfo"
	FROM "cards"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, card.ID).Return(pgxRows)

	targetCard, err := suite.store.Card(context.Background(), userID, card.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), *targetCard, card)

	//! Тест на ошибки sql
	targetErr := errors.New("some error")
	pgxRows = pgxpoolmock.NewRows(columns).
		AddRow(card.ID, userID, card.FIO, card.Number, card.Date, card.CVV, card.MetaInfo).
		RowError(0, targetErr).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "id", "user", "fio", "number", "date", "cvv", "metainfo"
	FROM "cards"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, card.ID).Return(pgxRows)

	targetCard, err = suite.store.Card(context.Background(), userID, card.ID)
	assert.Error(suite.T(), err, targetErr)
	assert.Nil(suite.T(), targetCard)
}

func (suite *PgStorageSuite) TestDeleteCard() {
	userID := int64(1)
	cardID := "testID"

	suite.mockPool.EXPECT().Exec(gomock.Any(), `
	DELETE FROM "cards"
	WHERE "user" = $1 and "id" = $2
	`, userID, cardID).Return([]byte("DELETE"), nil)

	err := suite.store.DeleteCard(context.Background(), userID, cardID)
	assert.NoError(suite.T(), err)
}

func (suite *PgStorageSuite) TestAddLogin() {
	userID := int64(1)
	login := models.Login{
		UserID:   userID,
		ID:       "testID",
		Login:    "login",
		Password: "password",
		MetaInfo: "metainfo",
	}

	//! Тест на корректную вставку
	columns := []string{"id"}
	pgxRows := pgxpoolmock.NewRows(columns).
		AddRow(int64(1)).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "logins"("id", "user", "login", "password", "metainfo")
	VALUES ($1, $2, $3, $4, $5)
	returning "@logins"
	`, login.ID, userID, login.Login, login.Password, login.MetaInfo).Return(pgxRows)

	err := suite.store.AddLogin(context.Background(), userID, login)
	assert.NoError(suite.T(), err)

	//! Тест на ошибки sql
	targetErr := errors.New("some error")
	pgxRows = pgxpoolmock.NewRows(columns).
		AddRow(int64(1)).
		RowError(0, targetErr).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "logins"("id", "user", "login", "password", "metainfo")
	VALUES ($1, $2, $3, $4, $5)
	returning "@logins"
	`, login.ID, userID, login.Login, login.Password, login.MetaInfo).Return(pgxRows)

	err = suite.store.AddLogin(context.Background(), userID, login)
	assert.Error(suite.T(), err, targetErr)
}

func (suite *PgStorageSuite) TestLogin() {
	userID := int64(1)
	login := models.Login{
		UserID:   userID,
		ID:       "testID",
		Login:    "login",
		Password: "password",
		MetaInfo: "metainfo",
	}

	columns := []string{"id", "user", "login", "password", "metainfo"}
	pgxRows := pgxpoolmock.NewRows(columns).
		AddRow(login.ID, userID, login.Login, login.Password, login.MetaInfo).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "id", "user", "login", "password", "metainfo"
	FROM "logins"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, login.ID).Return(pgxRows)

	targetLogin, err := suite.store.Login(context.Background(), userID, login.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), *targetLogin, login)

	//! Тест на ошибки sql
	targetErr := errors.New("some error")
	pgxRows = pgxpoolmock.NewRows(columns).
		AddRow(login.ID, userID, login.Login, login.Password, login.MetaInfo).
		RowError(0, targetErr).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "id", "user", "login", "password", "metainfo"
	FROM "logins"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, login.ID).Return(pgxRows)

	targetLogin, err = suite.store.Login(context.Background(), userID, login.ID)
	assert.Error(suite.T(), err, targetErr)
	assert.Nil(suite.T(), targetLogin)
}

func (suite *PgStorageSuite) TestDeleteLogin() {
	userID := int64(1)
	loginID := "testID"

	suite.mockPool.EXPECT().Exec(gomock.Any(), `
	DELETE FROM "logins"
	WHERE "user" = $1 and "id" = $2
	`, userID, loginID).Return([]byte("DELETE"), nil)

	err := suite.store.DeleteLogin(context.Background(), userID, loginID)
	assert.NoError(suite.T(), err)
}

func (suite *PgStorageSuite) TestAddText() {
	userID := int64(1)
	text := models.Text{
		UserID:   userID,
		ID:       "testID",
		Content:  "content",
		MetaInfo: "metainfo",
	}

	//! Тест на корректную вставку
	columns := []string{"id"}
	pgxRows := pgxpoolmock.NewRows(columns).
		AddRow(int64(1)).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "text_data"("id", "user", "content", "metainfo")
	VALUES ($1, $2, $3, $4)
	returning "@text"
	`, text.ID, userID, text.Content, text.MetaInfo).Return(pgxRows)

	err := suite.store.AddText(context.Background(), userID, text)
	assert.NoError(suite.T(), err)

	//! Тест на ошибки sql
	targetErr := errors.New("some error")
	pgxRows = pgxpoolmock.NewRows(columns).
		AddRow(int64(1)).
		RowError(0, targetErr).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "text_data"("id", "user", "content", "metainfo")
	VALUES ($1, $2, $3, $4)
	returning "@text"
	`, text.ID, userID, text.Content, text.MetaInfo).Return(pgxRows)

	err = suite.store.AddText(context.Background(), userID, text)
	assert.Error(suite.T(), err, targetErr)
}

func (suite *PgStorageSuite) TestText() {
	userID := int64(1)
	text := models.Text{
		UserID:   userID,
		ID:       "testID",
		Content:  "content",
		MetaInfo: "metainfo",
	}

	//! Тест на корректную вставку
	columns := []string{"id", "user", "content", "metainfo"}
	pgxRows := pgxpoolmock.NewRows(columns).
		AddRow(text.ID, userID, text.Content, text.MetaInfo).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "id", "user", "content", "metainfo"
	FROM "text_data"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, text.ID).Return(pgxRows)

	targetBin, err := suite.store.Text(context.Background(), userID, text.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), *targetBin, text)

	//! Тест на ошибки sql
	targetErr := errors.New("some error")
	pgxRows = pgxpoolmock.NewRows(columns).
		AddRow(text.ID, userID, text.Content, text.MetaInfo).
		RowError(0, targetErr).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "id", "user", "content", "metainfo"
	FROM "text_data"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, text.ID).Return(pgxRows)

	targetBin, err = suite.store.Text(context.Background(), userID, text.ID)
	assert.Error(suite.T(), err, targetErr)
	assert.Nil(suite.T(), targetBin)
}

func (suite *PgStorageSuite) TestDeleteText() {
	userID := int64(1)
	textID := "testID"

	suite.mockPool.EXPECT().Exec(gomock.Any(), `
	DELETE FROM "text_data"
	WHERE "user" = $1 and "id" = $2
	`, userID, textID).Return([]byte("DELETE"), nil)

	err := suite.store.DeleteText(context.Background(), userID, textID)
	assert.NoError(suite.T(), err)
}

func (suite *PgStorageSuite) TestAddBinary() {
	userID := int64(1)
	bin := models.Binary{
		UserID:   userID,
		ID:       "testID",
		Data:     []byte("data"),
		MetaInfo: "metainfo",
	}

	//! Тест на корректную вставку
	columns := []string{"id"}
	pgxRows := pgxpoolmock.NewRows(columns).
		AddRow(int64(1)).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "bin_data"("id", "user", "content", "metainfo")
	VALUES ($1, $2, $3, $4)
	returning "@bin"
	`, bin.ID, userID, bin.Data, bin.MetaInfo).Return(pgxRows)

	err := suite.store.AddBinary(context.Background(), userID, bin)
	assert.NoError(suite.T(), err)

	//! Тест на ошибки sql
	targetErr := errors.New("some error")
	pgxRows = pgxpoolmock.NewRows(columns).
		AddRow(int64(1)).
		RowError(0, targetErr).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	INSERT INTO "bin_data"("id", "user", "content", "metainfo")
	VALUES ($1, $2, $3, $4)
	returning "@bin"
	`, bin.ID, userID, bin.Data, bin.MetaInfo).Return(pgxRows)

	err = suite.store.AddBinary(context.Background(), userID, bin)
	assert.Error(suite.T(), err, targetErr)
}

func (suite *PgStorageSuite) TestBinary() {
	userID := int64(1)
	bin := models.Binary{
		UserID:   userID,
		ID:       "testID",
		Data:     []byte("data"),
		MetaInfo: "metainfo",
	}

	//! Тест на корректную вставку
	columns := []string{"id", "user", "content", "metainfo"}
	pgxRows := pgxpoolmock.NewRows(columns).
		AddRow(bin.ID, userID, bin.Data, bin.MetaInfo).ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "id", "user", "content", "metainfo"
	FROM "bin_data"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, bin.ID).Return(pgxRows)

	targetBin, err := suite.store.Binary(context.Background(), userID, bin.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), *targetBin, bin)

	//! Тест на ошибки sql
	targetErr := errors.New("some error")
	pgxRows = pgxpoolmock.NewRows(columns).
		AddRow(bin.ID, userID, bin.Data, bin.MetaInfo).
		RowError(0, targetErr).
		ToPgxRows()
	pgxRows.Next()

	suite.mockPool.EXPECT().QueryRow(gomock.Any(), `
	SELECT "id", "user", "content", "metainfo"
	FROM "bin_data"
	WHERE "user" = $1 and "id" = $2
	LIMIT 1
	`, userID, bin.ID).Return(pgxRows)

	targetBin, err = suite.store.Binary(context.Background(), userID, bin.ID)
	assert.Error(suite.T(), err, targetErr)
	assert.Nil(suite.T(), targetBin)
}

func (suite *PgStorageSuite) TestDeleteBinary() {
	userID := int64(1)
	binID := "testID"

	suite.mockPool.EXPECT().Exec(gomock.Any(), `
	DELETE FROM "bin_data"
	WHERE "user" = $1 and "id" = $2
	`, userID, binID).Return([]byte("DELETE"), nil)

	err := suite.store.DeleteBinary(context.Background(), userID, binID)
	assert.NoError(suite.T(), err)
}
