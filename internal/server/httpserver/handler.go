// Package httpserver реализует http обработчики через chi.NewRouter
package httpserver

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	"github.com/ncyellow/GophKeeper/internal/models"
	"github.com/ncyellow/GophKeeper/internal/server/auth"
	"github.com/ncyellow/GophKeeper/internal/server/auth/jwt"
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/storage"
)

type Handler struct {
	*chi.Mux
	store      storage.Storage
	authorizer *jwt.Authorizer
}

func NewRouter(conf *config.Config, store storage.Storage, parser jwt.Parser) chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	authorizer := &jwt.Authorizer{
		Store:      store,
		SigningKey: []byte(conf.SigningKey),
	}

	handler := Handler{
		Mux:        r,
		store:      store,
		authorizer: authorizer,
	}

	r.Group(func(r chi.Router) {
		r.Post("/api/register", handler.Register())
		r.Post("/api/signin", handler.SignIn())
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.Auth(store, conf, parser))
		// Тут будут обработчики ^_^

		// API для работы с банковскими картами
		r.Get("/api/card/{id}", handler.Card())
		r.Post("/api/card", handler.AddCard())
		r.Delete("/api/card/{id}", handler.DeleteCard())

		// API для работы с логинами
		r.Get("/api/login/{id}", handler.Login())
		r.Post("/api/login", handler.AddLogin())
		r.Delete("/api/login/{id}", handler.DeleteLogin())

		// API для работы с текстовыми данными
		r.Get("/api/txt/{id}", handler.Text())
		r.Post("/api/txt", handler.AddText())
		r.Delete("/api/txt/{id}", handler.DeleteText())

		// API для работы с бинарными данными
		r.Get("/api/bin/{id}", handler.Binary())
		r.Post("/api/bin", handler.AddBinary())
		r.Delete("/api/bin/{id}", handler.DeleteBinary())
	})
	return handler
}

// Register регистрация пользователя
func (h *Handler) Register() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Проверяем Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("content type not support"))
			return
		}
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}

		// разбираем сообщение
		var user models.User
		err = json.Unmarshal(reqBody, &user)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		originalPassword := user.Password

		pwd := sha1.New()
		pwd.Write([]byte(user.Password))
		hashPwd := fmt.Sprintf("%x", pwd.Sum(nil))
		user.Password = hashPwd

		// Выполняем попытку регистрации
		_, err = h.store.Register(r.Context(), user)
		if err != nil {
			rw.WriteHeader(http.StatusConflict)
			rw.Write([]byte("already have"))
			return
		}

		user.Password = originalPassword
		// Генерация токена если регистрация успешна
		jwtToken, err := h.authorizer.SignIn(r.Context(), &user)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid login"))
			return
		}
		rw.Header().Set("Authorization", jwtToken)
		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
	}
}

// SignIn аутентификация
func (h *Handler) SignIn() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Проверяем Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("content type not support"))
			return
		}
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Read data problem"))
			return
		}

		// разбираем сообщение
		var user models.User
		err = json.Unmarshal(reqBody, &user)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		// Попытка аутентификации если проходит - генерируем токен
		jwtToken, err := h.authorizer.SignIn(r.Context(), &user)

		// либо 200, либо 401
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("invalid login or password"))
			return
		}
		rw.Header().Set("Authorization", jwtToken)
		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
	}
}

// Card вернуть данные конкретной карты
func (h *Handler) Card() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		cardID := chi.URLParam(r, "id")
		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Запрашиваем инфу по карте
		targetCard, err := h.store.Card(r.Context(), user.UserID, cardID)
		if err != nil {
			if err == pgx.ErrNoRows {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(targetCard)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid serialization"))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(result)
	}
}

// AddCard зарегистрировать новую карту
func (h *Handler) AddCard() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}

		// Разбираем сообщение
		var cardData models.Card
		err = json.Unmarshal(reqBody, &cardData)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = h.store.AddCard(r.Context(), user.UserID, cardData)
		if err != nil {
			rw.WriteHeader(http.StatusConflict)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// DeleteCard зарегистрировать новую карту
func (h *Handler) DeleteCard() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		cardID := chi.URLParam(r, "id")
		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Удаляем логин
		err := h.store.DeleteCard(r.Context(), user.UserID, cardID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))

	}
}

// Login вернуть данные конкретной карты
func (h *Handler) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		loginID := chi.URLParam(r, "id")
		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Запрашиваем инфу по логину
		targetLogin, err := h.store.Login(r.Context(), user.UserID, loginID)
		if err != nil {
			if err == pgx.ErrNoRows {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(targetLogin)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid serialization"))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(result)
	}
}

// AddLogin зарегистрировать новую карту
func (h *Handler) AddLogin() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}

		// Разбираем сообщение
		var loginData models.Login
		err = json.Unmarshal(reqBody, &loginData)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = h.store.AddLogin(r.Context(), user.UserID, loginData)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// DeleteLogin зарегистрировать новую карту
func (h *Handler) DeleteLogin() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		loginID := chi.URLParam(r, "id")
		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Удаляем логин
		err := h.store.DeleteLogin(r.Context(), user.UserID, loginID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// Text вернуть данные конкретной карты
func (h *Handler) Text() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		textID := chi.URLParam(r, "id")
		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Запрашиваем инфу по текстовым данным
		targetText, err := h.store.Text(r.Context(), user.UserID, textID)
		if err != nil {
			if err == pgx.ErrNoRows {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(targetText)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid serialization"))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(result)
	}
}

// AddText зарегистрировать новую карту
func (h *Handler) AddText() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}

		// Разбираем сообщение
		var textData models.Text
		err = json.Unmarshal(reqBody, &textData)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = h.store.AddText(r.Context(), user.UserID, textData)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// DeleteText зарегистрировать новую карту
func (h *Handler) DeleteText() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		textID := chi.URLParam(r, "id")
		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Удаляем текст
		err := h.store.DeleteText(r.Context(), user.UserID, textID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// Binary вернуть данные конкретной карты
func (h *Handler) Binary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		binID := chi.URLParam(r, "id")
		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Запрашиваем инфу по бинарным
		targetBin, err := h.store.Binary(r.Context(), user.UserID, binID)
		if err != nil {
			if err == pgx.ErrNoRows {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(targetBin)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid serialization"))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(result)
	}
}

// AddBinary зарегистрировать новую карту
func (h *Handler) AddBinary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}

		// Разбираем сообщение
		var binData models.Binary
		err = json.Unmarshal(reqBody, &binData)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Запрашиваем инфу по бинарным
		err = h.store.AddBinary(r.Context(), user.UserID, binData)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// DeleteBinary зарегистрировать новую карту
func (h *Handler) DeleteBinary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		binID := chi.URLParam(r, "id")
		user, ok := r.Context().Value(auth.UserContextKey{}).(*models.User)
		if !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Удаляем бинарь
		err := h.store.DeleteBinary(r.Context(), user.UserID, binID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}
