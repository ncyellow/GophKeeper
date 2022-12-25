// Package api модуль реализует http client для взаимодействия с сервером
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ncyellow/GophKeeper/internal/client/config"
	"github.com/ncyellow/GophKeeper/internal/models"
)

var (
	ErrInternalServer    = errors.New("cервер недоступен, попробуйте позднее")
	ErrServerTimout      = errors.New("cервер недоступен, попробуйте позднее")
	ErrSerialization     = errors.New("ошибка сериализации")
	ErrDeserialization   = errors.New("ошибка десериализации")
	ErrRequestPrepare    = errors.New("не удалось подготовить http запрос")
	ErrUserAlreadyExists = errors.New("уже зарегистрирован пользователь с таким логином")
	ErrUserNotFound      = errors.New("пользователь с таким логином не найден")
	ErrAuthRequire       = errors.New("необходим авторизоваться")

	ErrAlreadyExists = errors.New("ID с таким идентификатором уже зарегистрирован")
	ErrNotFound      = errors.New("не найдена запись с таким идентификатором")
)

// HTTPSender структура http клиента
type HTTPSender struct {
	Client    *http.Client
	Conf      *config.Config
	AuthToken *string
}

// NewHTTPSender конструктор http клиента
func NewHTTPSender(conf *config.Config) *HTTPSender {
	return &HTTPSender{
		Client:    &http.Client{},
		Conf:      conf,
		AuthToken: nil,
	}
}

// Register запрос на сервер по регистрации клиента
func (s *HTTPSender) Register(login string, pwd string) error {
	user := models.User{
		Login:    login,
		Password: pwd,
	}

	result, ok := json.Marshal(user)
	if ok != nil {
		return ErrSerialization
	}

	req, err := http.NewRequest("POST", s.Conf.Address+"/api/register", bytes.NewBuffer(result))
	if err != nil {
		return ErrRequestPrepare
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ErrServerTimout
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return ErrUserAlreadyExists
	} else if resp.StatusCode != http.StatusOK {
		return ErrInternalServer
	}

	authToken := resp.Header.Get("Authorization")
	s.AuthToken = &authToken
	return nil
}

// SignIn запрос на сервер по авторизации клиента
func (s *HTTPSender) SignIn(login string, pwd string) error {

	user := models.User{
		Login:    login,
		Password: pwd,
	}

	result, ok := json.Marshal(user)
	if ok != nil {
		return ErrSerialization
	}

	req, err := http.NewRequest("POST", s.Conf.Address+"/api/signin", bytes.NewBuffer(result))
	if err != nil {
		return ErrRequestPrepare
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ErrServerTimout
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrUserNotFound
	}

	authToken := resp.Header.Get("Authorization")
	s.AuthToken = &authToken
	return nil
}

// AddCard запрос добавления новой карты
func (s *HTTPSender) AddCard(card *models.Card) error {
	data, ok := json.Marshal(card)
	if ok != nil {
		return ErrSerialization
	}
	return s.Add(data, "/api/card")
}

// Card запрос на чтение уже существующей по ид карты
func (s *HTTPSender) Card(cardID string) (*models.Card, error) {
	data, err := s.Read(cardID, "api/card")
	if err != nil {
		return nil, err
	}
	// разбираем сообщение
	var card models.Card
	err = json.Unmarshal(data, &card)

	if err != nil {
		return nil, ErrDeserialization
	}
	return &card, nil
}

// DelCard запрос на удаление уже существующей по ид карты
func (s *HTTPSender) DelCard(cardID string) error {
	return s.Del(cardID, "api/card")
}

// AddLogin запрос добавления нового логина
func (s *HTTPSender) AddLogin(login *models.Login) error {
	data, ok := json.Marshal(login)
	if ok != nil {
		return ErrSerialization
	}
	return s.Add(data, "/api/login")
}

// Login запрос на чтение уже существующего логина по ид
func (s *HTTPSender) Login(loginID string) (*models.Login, error) {

	data, err := s.Read(loginID, "api/login")
	if err != nil {
		return nil, err
	}
	// разбираем сообщение
	var login models.Login
	err = json.Unmarshal(data, &login)

	if err != nil {
		return nil, ErrDeserialization
	}
	return &login, nil
}

// DelLogin запрос на удаление уже существующего логина по ид
func (s *HTTPSender) DelLogin(loginID string) error {
	return s.Del(loginID, "api/login")
}

// AddText запрос на чтение уже существующего текста по ид
func (s *HTTPSender) AddText(text *models.Text) error {

	data, ok := json.Marshal(text)
	if ok != nil {
		return ErrSerialization
	}
	return s.Add(data, "api/txt")
}

// Text запрос на чтение уже существующего текста по ид
func (s *HTTPSender) Text(textID string) (*models.Text, error) {
	data, err := s.Read(textID, "api/txt")
	if err != nil {
		return nil, err
	}

	// разбираем сообщение
	var text models.Text
	err = json.Unmarshal(data, &text)
	if err != nil {
		return nil, ErrDeserialization
	}
	return &text, nil
}

// DelText запрос на удаление уже существующего текста по ид
func (s *HTTPSender) DelText(textID string) error {
	return s.Del(textID, "api/txt")
}

// AddBin запрос на добавление уже существующих бинарных данных по ид
func (s *HTTPSender) AddBin(binary *models.Binary) error {
	data, ok := json.Marshal(binary)
	if ok != nil {
		return ErrSerialization
	}
	return s.Add(data, "api/bin")
}

// Bin запрос на чтение уже существующего бинарных данных по ид
func (s *HTTPSender) Bin(binID string) (*models.Binary, error) {
	data, err := s.Read(binID, "api/bin")
	if err != nil {
		return nil, err
	}
	// разбираем сообщение
	var binary models.Binary
	err = json.Unmarshal(data, &binary)

	if err != nil {
		return nil, ErrDeserialization
	}
	return &binary, nil
}

// DelBin запрос на удаление уже существующих бинарных данных по ид
func (s *HTTPSender) DelBin(binID string) error {
	return s.Del(binID, "api/bin")
}

// Add общий метод по добавлению на сервер. Содержит общую часть для любого типа данных
func (s *HTTPSender) Add(data []byte, urlSuffix string) error {
	if s.AuthToken == nil {
		return ErrAuthRequire
	}

	req, err := http.NewRequest("POST", s.Conf.Address+urlSuffix, bytes.NewBuffer(data))
	if err != nil {
		return ErrRequestPrepare
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", *s.AuthToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ErrServerTimout
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return ErrAlreadyExists
	} else if resp.StatusCode != http.StatusOK {
		return ErrInternalServer
	}
	return nil
}

// Read общий метод по чтение с сервера. Содержит общую часть для любого типа данных
func (s *HTTPSender) Read(textID string, urlSuffix string) ([]byte, error) {
	if s.AuthToken == nil {
		return nil, ErrAuthRequire
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", s.Conf.Address, urlSuffix, textID), nil)
	if err != nil {
		return nil, ErrRequestPrepare
	}
	req.Header.Set("Authorization", *s.AuthToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ErrServerTimout
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrInternalServer
	}

	reqBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, ErrSerialization
	}
	return reqBody, nil
}

// Del общий метод по удалению с сервера. Содержит общую часть для любого типа данных
func (s *HTTPSender) Del(binID string, urlSuffix string) error {
	if s.AuthToken == nil {
		return ErrAuthRequire
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", s.Conf.Address, urlSuffix, binID), nil)
	if err != nil {
		return ErrRequestPrepare
	}

	req.Header.Set("Authorization", *s.AuthToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ErrServerTimout
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrInternalServer
	}
	return nil
}
