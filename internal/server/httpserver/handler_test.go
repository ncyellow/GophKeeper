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

// TestRegisterHandler основные тесты по регистрации
func (suite *HandlersSuite) TestCard() {

	userID := 1
	defaultCard := &models.Card{
		UserID:   userID,
		ID:       "test_id",
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
			request:     fmt.Sprintf("/api/card/%d", userID),
			requestType: "GET",
			contentType: "",
			body:        nil,
			mockExpected: func() {
				user := &models.User{
					UserID: int64(userID),
					Login:  "login",
				}
				suite.parser.EXPECT().ParseToken(gomock.Any(), gomock.Any()).Return("login", nil)
				suite.store.EXPECT().UserByLogin(gomock.Any(), gomock.Any()).Return(user, nil)
				suite.store.EXPECT().Card(gomock.Any(), user.UserID, fmt.Sprintf("%d", userID)).
					Return(defaultCard, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       string(byteCard),
			},
		},
	}
	suite.runTableTests(testData)
}
