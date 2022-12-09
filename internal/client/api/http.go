// Package api модуль реализует http client для взаимодействия с сервером
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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

type HTTPSender struct {
	Client    *http.Client
	URL       string
	AuthToken *string
}

func NewHTTPSender() *HTTPSender {
	return &HTTPSender{
		Client:    &http.Client{},
		URL:       "http://localhost:8085",
		AuthToken: nil,
	}
}

func (s *HTTPSender) Register(login string, pwd string) error {
	user := models.User{
		Login:    login,
		Password: pwd,
	}

	result, ok := json.Marshal(user)
	if ok != nil {
		return ErrSerialization
	}

	req, err := http.NewRequest("POST", s.URL+"/api/register", bytes.NewBuffer(result))
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
	//fmt.Printf("Пользователь с логином, %s, успешно зарегистрирован!\n", user.Login)
	return nil
}

func (s *HTTPSender) SignIn(login string, pwd string) error {

	user := models.User{
		Login:    login,
		Password: pwd,
	}

	result, ok := json.Marshal(user)
	if ok != nil {
		return ErrSerialization
	}

	req, err := http.NewRequest("POST", s.URL+"/api/signin", bytes.NewBuffer(result))
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

func (s *HTTPSender) AddCard(card *models.Card) error {
	data, ok := json.Marshal(card)
	if ok != nil {
		return ErrSerialization
	}
	return s.Add(data, "/api/card")
}

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

func (s *HTTPSender) DelCard(cardID string) error {
	return s.Del(cardID, "api/card")
}

func (s *HTTPSender) AddLogin(login *models.Login) error {
	data, ok := json.Marshal(login)
	if ok != nil {
		return ErrSerialization
	}
	return s.Add(data, "/api/login")
}

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

func (s *HTTPSender) DelLogin(loginID string) error {
	return s.Del(loginID, "api/login")
}

func (s *HTTPSender) AddText(text *models.Text) error {

	data, ok := json.Marshal(text)
	if ok != nil {
		return ErrSerialization
	}
	return s.Add(data, "api/txt")
}

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

func (s *HTTPSender) DelText(textID string) error {
	return s.Del(textID, "api/txt")
}

func (s *HTTPSender) AddBin(binary *models.Binary) error {
	data, ok := json.Marshal(binary)
	if ok != nil {
		return ErrSerialization
	}
	return s.Add(data, "api/bin")
}

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

func (s *HTTPSender) DelBin(binID string) error {
	return s.Del(binID, "api/bin")
}

func (s *HTTPSender) Add(data []byte, urlSuffix string) error {
	if s.AuthToken == nil {
		return ErrAuthRequire
	}

	req, err := http.NewRequest("POST", s.URL+urlSuffix, bytes.NewBuffer(data))
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

func (s *HTTPSender) Read(textID string, urlSuffix string) ([]byte, error) {
	if s.AuthToken == nil {
		return nil, ErrAuthRequire
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", s.URL, urlSuffix, textID), nil)
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

func (s *HTTPSender) Del(binID string, urlSuffix string) error {
	if s.AuthToken == nil {
		return ErrAuthRequire
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", s.URL, urlSuffix, binID), nil)
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
