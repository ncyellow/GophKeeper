package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/ncyellow/GophKeeper/internal/models"
	"github.com/ncyellow/GophKeeper/internal/server/config"
	mock_jwt "github.com/ncyellow/GophKeeper/internal/server/mocks/auth/jwt"
	mock_storage "github.com/ncyellow/GophKeeper/internal/server/mocks/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type want struct {
	statusCode int
	body       string
}
type tests struct {
	name         string
	request      string
	requestType  string
	contentType  string
	body         []byte
	mockExpected func()
	want         want
}

type HandlersSuite struct {
	suite.Suite
	store  *mock_storage.MockStorage
	parser *mock_jwt.MockParser
	ts     *httptest.Server
}

// SetupSuite перед началом теста стартуем новый сервер httptest.Server делаем так, чтобы тестировать каждый
// handler отдельно и не сливать все тесты в один
func (suite *HandlersSuite) SetupTest() {

	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	conf := config.Config{}
	store := mock_storage.NewMockStorage(ctrl)
	suite.store = store

	parser := mock_jwt.NewMockParser(ctrl)
	suite.parser = parser

	r := NewRouter(&conf, store, parser)
	suite.ts = httptest.NewServer(r)
}

// TearDownSuite после теста отключаем сервер
func (suite *HandlersSuite) TearDownTest() {
	suite.store.EXPECT().Close()
	suite.store.Close()
	suite.ts.Close()
}

// TestHandlersSuite старт нашего HandlersSuite
func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlersSuite))
}

// runTestRequest вспомогательная функция для выполнения http запроса
func runTestRequest(t *testing.T, ts *httptest.Server, method, path string, contentType string, reqBody []byte) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

// runTableTests хелпер на запуск группы списочных тестов
func (suite *HandlersSuite) runTableTests(testList []tests) {
	for _, tt := range testList {
		if tt.mockExpected != nil {
			tt.mockExpected()
		}
		resp, body := runTestRequest(suite.T(), suite.ts, tt.requestType, tt.request, tt.contentType, tt.body)
		assert.Equal(suite.T(), tt.want.statusCode, resp.StatusCode, tt.name)
		assert.Equal(suite.T(), tt.want.body, body, tt.name)
		resp.Body.Close()
	}
}

// TestRegisterHandler основные тесты по регистрации
func (suite *HandlersSuite) TestRegisterHandler() {
	testData := []tests{
		{
			name:         "register with wrong content-type",
			request:      "/api/register",
			requestType:  "POST",
			contentType:  "plain/text",
			body:         nil,
			mockExpected: nil,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "content type not support",
			},
		},
		{
			name:         "register incorrect user data",
			request:      "/api/register",
			requestType:  "POST",
			contentType:  "application/json",
			body:         []byte(`{"login": "login", no password}`),
			mockExpected: nil,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid deserialization",
			},
		},
		{
			name:        "register with success",
			request:     "/api/register",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"login": "login", "password": "password"}`),

			mockExpected: func() {
				// Регистрация успешная вернулся id без ошибок
				suite.store.EXPECT().Register(gomock.Any(), models.User{
					Login:    "login",
					Password: "5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8", //sha1
				}).Return(int64(1), nil)

				// После регистрации аутентификация
				suite.store.EXPECT().User(gomock.Any(), "login",
					"5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8").Return(&models.User{
					Login:    "login",
					Password: "5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8",
				}, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "",
			},
		},
		{
			name:        "register with conflict",
			request:     "/api/register",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"login": "login", "password": "password"}`),

			mockExpected: func() {
				// Регистрация успешная вернулся id без ошибок
				suite.store.EXPECT().Register(gomock.Any(), models.User{
					Login:    "login",
					Password: "5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8", //sha1
				}).Return(int64(0), errors.New("SomeError"))
			},
			want: want{
				statusCode: http.StatusConflict,
				body:       "already have",
			},
		},
		//! тут еще нужен тест на проблемы базы, но не конфликт вставки в просто проблемы с базой
	}
	suite.runTableTests(testData)
}

// TestRegisterHandler основные тесты по регистрации
func (suite *HandlersSuite) TestSignIn() {
	testData := []tests{
		{
			name:         "signin incorrect user data",
			request:      "/api/signin",
			requestType:  "POST",
			contentType:  "application/json",
			body:         []byte(`{"login": "login", no password}`),
			mockExpected: nil,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid deserialization",
			},
		},
		{
			name:        "signin with success",
			request:     "/api/signin",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"login": "login", "password": "password"}`),

			mockExpected: func() {
				suite.store.EXPECT().User(gomock.Any(), "login",
					"5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8").Return(&models.User{
					Login:    "login",
					Password: "5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8",
				}, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "",
			},
		},
		{
			name:        "signin with any error",
			request:     "/api/signin",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"login": "login", "password": "password"}`),

			mockExpected: func() {
				suite.store.EXPECT().User(gomock.Any(), "login",
					"5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8").Return(nil, errors.New("any error"))
			},
			want: want{
				statusCode: http.StatusUnauthorized,
				body:       "invalid login or password",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestCard тесты по чтение карт
func (suite *HandlersSuite) TestCard() {

	userID := int64(1)
	cardID := "testID"
	url := fmt.Sprintf("/api/card/%s", cardID)
	defaultCard := &models.Card{
		UserID:   userID,
		ID:       cardID,
		FIO:      "fio",
		Number:   "number",
		Date:     "date",
		CVV:      "cvv",
		MetaInfo: "metainfo",
	}
	byteCard, _ := json.Marshal(defaultCard)

	testData := []tests{
		{
			name:        "read card successfully",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: int64(userID),
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Card(gomock.Any(), user.UserID, cardID).
					Return(defaultCard, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       string(byteCard),
			},
		},
		{
			name:        "read card not found",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: int64(userID),
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Card(gomock.Any(), user.UserID, cardID).
					Return(nil, pgx.ErrNoRows)
			},
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:        "read card internal error",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: int64(userID),
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Card(gomock.Any(), user.UserID, cardID).
					Return(nil, errors.New("some error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestLogin тесты по запросу логинов
func (suite *HandlersSuite) TestLogin() {

	userID := int64(1)
	loginID := "testID"
	url := fmt.Sprintf("/api/login/%s", loginID)

	defaultLogin := &models.Login{
		UserID:   userID,
		ID:       loginID,
		Login:    "login",
		Password: "password",
		MetaInfo: "metainfo",
	}
	byteLogin, _ := json.Marshal(defaultLogin)

	testData := []tests{
		{
			name:        "read login successfully",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Login(gomock.Any(), user.UserID, loginID).
					Return(defaultLogin, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       string(byteLogin),
			},
		},
		{
			name:        "read login not found",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Login(gomock.Any(), user.UserID, loginID).
					Return(nil, pgx.ErrNoRows)
			},
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:        "read login internal error",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Login(gomock.Any(), user.UserID, loginID).
					Return(nil, errors.New("some error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestText тесты по запросу логинов
func (suite *HandlersSuite) TestText() {

	userID := int64(1)
	textID := "testID"
	url := fmt.Sprintf("/api/txt/%s", textID)

	defaultText := &models.Text{
		UserID:   userID,
		ID:       textID,
		Content:  "content",
		MetaInfo: "metainfo",
	}
	byteText, _ := json.Marshal(defaultText)

	testData := []tests{
		{
			name:        "read text successfully",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Text(gomock.Any(), user.UserID, textID).
					Return(defaultText, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       string(byteText),
			},
		},
		{
			name:        "read text not found",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Text(gomock.Any(), user.UserID, textID).
					Return(nil, pgx.ErrNoRows)
			},
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:        "read text internal error",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Text(gomock.Any(), user.UserID, textID).
					Return(nil, errors.New("some error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestText тесты по запросу логинов
func (suite *HandlersSuite) TestBin() {

	userID := int64(1)
	binID := "testID"
	url := fmt.Sprintf("/api/bin/%s", binID)

	defaultBin := &models.Binary{
		UserID:   userID,
		ID:       binID,
		Data:     []byte("data"),
		MetaInfo: "metainfo",
	}
	byteBin, _ := json.Marshal(defaultBin)

	testData := []tests{
		{
			name:        "read text successfully",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Binary(gomock.Any(), user.UserID, binID).
					Return(defaultBin, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       string(byteBin),
			},
		},
		{
			name:        "read text not found",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Binary(gomock.Any(), user.UserID, binID).
					Return(nil, pgx.ErrNoRows)
			},
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:        "read text internal error",
			request:     url,
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().Binary(gomock.Any(), user.UserID, binID).
					Return(nil, errors.New("some error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestAddCard тесты по добавлению карт
func (suite *HandlersSuite) TestAddCard() {

	userID := int64(1)
	cardID := "testID"
	url := "/api/card"

	defaultCard := &models.Card{
		ID:       cardID,
		FIO:      "fio",
		Number:   "number",
		Date:     "date",
		CVV:      "cvv",
		MetaInfo: "metainfo",
	}

	byteCard, _ := json.Marshal(defaultCard)

	testData := []tests{
		{
			name:        "add card successfully",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        byteCard,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().AddCard(gomock.Any(), user.UserID, *defaultCard).
					Return(nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "add card with conflict",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        byteCard,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().AddCard(gomock.Any(), user.UserID, *defaultCard).
					Return(errors.New("some error"))
			},
			want: want{
				statusCode: http.StatusConflict,
				body:       "",
			},
		},
		{
			name:        "add card invalid json body",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        []byte(`"{"test":dd,ddla,la,dlww`),
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid deserialization",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestDelCard тесты по удалению карт
func (suite *HandlersSuite) TestDelCard() {

	userID := int64(1)
	cardID := "testID"
	url := fmt.Sprintf("/api/card/%s", cardID)

	testData := []tests{
		{
			name:        "del card successfully",
			request:     url,
			requestType: "DELETE",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().DeleteCard(gomock.Any(), user.UserID, cardID).
					Return(nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "del card with problem",
			request:     url,
			requestType: "DELETE",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().DeleteCard(gomock.Any(), user.UserID, cardID).
					Return(errors.New("some errors"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestAddLogin тесты по добавлению карт
func (suite *HandlersSuite) TestAddLogin() {

	userID := int64(1)
	loginID := "testID"
	url := "/api/login"

	defaultLogin := &models.Login{
		ID:       loginID,
		Login:    "login",
		Password: "password",
		MetaInfo: "metainfo",
	}
	byteLogin, _ := json.Marshal(defaultLogin)

	testData := []tests{
		{
			name:        "add login successfully",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        byteLogin,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().AddLogin(gomock.Any(), user.UserID, *defaultLogin).
					Return(nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "add login with conflict",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        byteLogin,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().AddLogin(gomock.Any(), user.UserID, *defaultLogin).
					Return(errors.New("some error"))
			},
			want: want{
				statusCode: http.StatusConflict,
				body:       "",
			},
		},
		{
			name:        "add login invalid json body",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        []byte(`"{"test":dd,ddla,la,dlww`),
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid deserialization",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestDelLogin тесты по удалению логинов
func (suite *HandlersSuite) TestDelLogin() {

	userID := int64(1)
	loginID := "testID"
	url := fmt.Sprintf("/api/login/%s", loginID)

	testData := []tests{
		{
			name:        "del login successfully",
			request:     url,
			requestType: "DELETE",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().DeleteLogin(gomock.Any(), user.UserID, loginID).
					Return(nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "del card with problem",
			request:     url,
			requestType: "DELETE",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().DeleteLogin(gomock.Any(), user.UserID, loginID).
					Return(errors.New("some errors"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestAddText тесты по добавлению текста
func (suite *HandlersSuite) TestAddText() {

	userID := int64(1)
	textID := "testID"
	url := "/api/txt"

	defaultText := &models.Text{
		ID:       textID,
		Content:  "content",
		MetaInfo: "metainfo",
	}
	byteText, _ := json.Marshal(defaultText)

	testData := []tests{
		{
			name:        "add text successfully",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        byteText,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().AddText(gomock.Any(), user.UserID, *defaultText).
					Return(nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "add text with conflict",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        byteText,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().AddText(gomock.Any(), user.UserID, *defaultText).
					Return(errors.New("some error"))
			},
			want: want{
				statusCode: http.StatusConflict,
				body:       "",
			},
		},
		{
			name:        "add text invalid json body",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        []byte(`"{"test":dd,ddla,la,dlww`),
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid deserialization",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestDelText тесты по удалению текстов
func (suite *HandlersSuite) TestDelText() {

	userID := int64(1)
	textID := "testID"
	url := fmt.Sprintf("/api/txt/%s", textID)

	testData := []tests{
		{
			name:        "del text successfully",
			request:     url,
			requestType: "DELETE",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().DeleteText(gomock.Any(), user.UserID, textID).
					Return(nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "del card with problem",
			request:     url,
			requestType: "DELETE",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().DeleteText(gomock.Any(), user.UserID, textID).
					Return(errors.New("some errors"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestAddBin тесты по добавлению binary
func (suite *HandlersSuite) TestAddBin() {

	userID := int64(1)
	binID := "testID"
	url := "/api/bin"

	defaultBin := &models.Binary{
		ID:       binID,
		Data:     []byte("data"),
		MetaInfo: "metainfo",
	}
	byteBin, _ := json.Marshal(defaultBin)

	testData := []tests{
		{
			name:        "add bin successfully",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        byteBin,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().AddBinary(gomock.Any(), user.UserID, *defaultBin).
					Return(nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "add bin with conflict",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        byteBin,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().AddBinary(gomock.Any(), user.UserID, *defaultBin).
					Return(errors.New("some error"))
			},
			want: want{
				statusCode: http.StatusConflict,
				body:       "",
			},
		},
		{
			name:        "add bin invalid json body",
			request:     url,
			requestType: "POST",
			contentType: "",
			body:        []byte(`"{"test":dd,ddla,la,dlww`),
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid deserialization",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestDelBin тесты по удалению binary
func (suite *HandlersSuite) TestDelBin() {

	userID := int64(1)
	binID := "testID"
	url := fmt.Sprintf("/api/bin/%s", binID)

	testData := []tests{
		{
			name:        "del bin successfully",
			request:     url,
			requestType: "DELETE",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().DeleteBinary(gomock.Any(), user.UserID, binID).
					Return(nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "del bin with problem",
			request:     url,
			requestType: "DELETE",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: userID,
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return(user.Login, nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), user.Login).Return(user, nil)
				suite.store.EXPECT().DeleteBinary(gomock.Any(), user.UserID, binID).
					Return(errors.New("some errors"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "",
			},
		},
	}
	suite.runTableTests(testData)
}
