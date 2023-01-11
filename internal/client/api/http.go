package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ncyellow/GophKeeper/internal/client/config"
	"github.com/ncyellow/GophKeeper/internal/models"
)

// HTTPSender структура http клиента. Реализует интерфейс Sender. Все комментарий по соотв. методам см там.
type HTTPSender struct {
	Client    *http.Client
	Conf      *config.Config
	AuthToken *string
}

// NewHTTPSender конструктор http клиента
func NewHTTPSender(conf *config.Config) *HTTPSender {
	// так как нам важна работа через tls, то все проблемы с tls вызывают фаталити
	clientCertFile := conf.CryptoCrt
	clientKeyFile := conf.CryptoKey
	caCertFile := conf.CACertFile

	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", clientCertFile, clientKeyFile)
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Error opening cert file %s, Error: %s", caCertFile, err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}

	return &HTTPSender{
		Client: &http.Client{
			Transport: t,
		},
		Conf:      conf,
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

	req, err := http.NewRequest("POST", s.Conf.Address+"/api/register", bytes.NewBuffer(result))
	if err != nil {
		return ErrRequestPrepare
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
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

	resp, err := s.Client.Do(req)
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
	return s.add(data, "/api/card")
}

func (s *HTTPSender) Card(cardID string) (*models.Card, error) {
	data, err := s.read(cardID, "api/card")
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
	return s.del(cardID, "api/card")
}

func (s *HTTPSender) AddLogin(login *models.Login) error {
	data, ok := json.Marshal(login)
	if ok != nil {
		return ErrSerialization
	}
	return s.add(data, "/api/login")
}

func (s *HTTPSender) Login(loginID string) (*models.Login, error) {
	data, err := s.read(loginID, "api/login")
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
	return s.del(loginID, "api/login")
}

func (s *HTTPSender) AddText(text *models.Text) error {
	data, ok := json.Marshal(text)
	if ok != nil {
		return ErrSerialization
	}
	return s.add(data, "api/txt")
}

func (s *HTTPSender) Text(textID string) (*models.Text, error) {
	data, err := s.read(textID, "api/txt")
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
	return s.del(textID, "api/txt")
}

func (s *HTTPSender) AddBin(binary *models.Binary) error {
	data, ok := json.Marshal(binary)
	if ok != nil {
		return ErrSerialization
	}
	return s.add(data, "api/bin")
}

func (s *HTTPSender) Bin(binID string) (*models.Binary, error) {
	data, err := s.read(binID, "api/bin")
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
	return s.del(binID, "api/bin")
}

// add общий метод по добавлению на сервер. Содержит общую часть для любого типа данных
func (s *HTTPSender) add(data []byte, urlSuffix string) error {
	if s.AuthToken == nil {
		return ErrAuthRequire
	}

	req, err := http.NewRequest("POST", s.Conf.Address+urlSuffix, bytes.NewBuffer(data))
	if err != nil {
		return ErrRequestPrepare
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", *s.AuthToken)

	resp, err := s.Client.Do(req)
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

// read общий метод по чтение с сервера. Содержит общую часть для любого типа данных
func (s *HTTPSender) read(textID string, urlSuffix string) ([]byte, error) {
	if s.AuthToken == nil {
		return nil, ErrAuthRequire
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", s.Conf.Address, urlSuffix, textID), nil)
	if err != nil {
		return nil, ErrRequestPrepare
	}
	req.Header.Set("Authorization", *s.AuthToken)

	resp, err := s.Client.Do(req)
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

// del общий метод по удалению с сервера. Содержит общую часть для любого типа данных
func (s *HTTPSender) del(binID string, urlSuffix string) error {
	if s.AuthToken == nil {
		return ErrAuthRequire
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", s.Conf.Address, urlSuffix, binID), nil)
	if err != nil {
		return ErrRequestPrepare
	}

	req.Header.Set("Authorization", *s.AuthToken)

	resp, err := s.Client.Do(req)
	if err != nil {
		return ErrServerTimout
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrInternalServer
	}
	return nil
}
