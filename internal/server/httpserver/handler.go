// Package httpserver implements HTTP handlers using chi.NewRouter
package httpserver

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
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

// @Title GophKeeper API
// @Description Service for storing confidential data
// @Version 1.0

// @Contact.email ncyellow@yandex.ru

// @Tag.name Add
// @Tag.description "Group of requests for adding new data"

// @Tag.name Read
// @Tag.description "Group of requests for reading data"

// @Tag.name Delete
// @Tag.description "Group of requests for deleting data"

// Handler structure implements chi.Mux for routing functionality
type Handler struct {
	*chi.Mux
	store      storage.Storage
	authorizer *jwt.Authorizer
}

// NewRouter constructor of our routing object
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
		// Here will be the handlers ^_^

		// API for working with bank cards
		r.Get("/api/card/{id}", handler.Card())
		r.Post("/api/card", handler.AddCard())
		r.Delete("/api/card/{id}", handler.DeleteCard())

		// API for working with logins
		r.Get("/api/login/{id}", handler.Login())
		r.Post("/api/login", handler.AddLogin())
		r.Delete("/api/login/{id}", handler.DeleteLogin())

		// API for working with text data
		r.Get("/api/txt/{id}", handler.Text())
		r.Post("/api/txt", handler.AddText())
		r.Delete("/api/txt/{id}", handler.DeleteText())

		// API for working with binary data
		r.Get("/api/bin/{id}", handler.Binary())
		r.Post("/api/bin", handler.AddBinary())
		r.Delete("/api/bin/{id}", handler.DeleteBinary())
	})
	return handler
}

// Register register user
func (h *Handler) Register() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// check Content-Type
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

		// parse message
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

		// Attempting registration
		_, err = h.store.Register(r.Context(), user)
		if err != nil {
			rw.WriteHeader(http.StatusConflict)
			rw.Write([]byte("already have"))
			return
		}

		user.Password = originalPassword
		// Generating token if registration is successful
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

// SignIn authentication
func (h *Handler) SignIn() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// check Content-Type
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

		// parse message
		var user models.User
		err = json.Unmarshal(reqBody, &user)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		// Attempting authentication, if successful - generate token
		jwtToken, err := h.authorizer.SignIn(r.Context(), &user)
		// Either 200 or 401
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

// Card return specific card data
// @Tags Read
// @Summary Returns user card data
// @Description input rest URL, output JSON value
// @ID readCard
// @Produce json
// @Param id path string true "Card ID"
// @Success 200 {object} Card
// @Failure 404 {string} string ""
// @Failure 500 {string} string ""
// @Router /api/card/{id} [get]
func (h *Handler) Card() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		cardID := chi.URLParam(r, "id")
		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		// Requesting card information
		targetCard, err := h.store.Card(r.Context(), user.UserID, cardID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
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

// AddCard register a new card
// @Tags Add
// @Summary Registering a new card
// @Description Registration is performed using a unique pair of User ID + Card ID.
// @ID addCard
// @Accept json
// @Produce plain
// @Param card_data body Card true "Card object"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid deserialization"
// @Failure 409 {string} string ""
// @Failure 500 {string} string "read data problem"
// @Router /api/card [post]
func (h *Handler) AddCard() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("read data problem"))
			return
		}

		// check message
		var cardData models.Card
		err = json.Unmarshal(reqBody, &cardData)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

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

// DeleteCard delete a card by user and ID
// @Tags Delete
// @Summary Deleting a card
// @Description Deletion is performed using a unique pair of User ID + Card ID.
// @ID delCard
// @Produce plain
// @Param id path string true "Card ID"
// @Success 200 {string} string "ok"
// @Failure 500 {string} string ""
// @Router /api/card [delete]
func (h *Handler) DeleteCard() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		cardID := chi.URLParam(r, "id")
		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		// delete login
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

// Login return specific login data
// @Tags Read
// @Summary Returns user login data
// @Description input rest URL, output JSON value
// @ID readLogin
// @Produce json
// @Param id path string true "Login ID"
// @Success 200 {object} Login
// @Failure 404 {string} string ""
// @Failure 500 {string} string ""
// @Router /api/login/{id} [get]
func (h *Handler) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		loginID := chi.URLParam(r, "id")
		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		// Requesting login information
		targetLogin, err := h.store.Login(r.Context(), user.UserID, loginID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
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

// AddLogin register a new login
// @Tags Add
// @Summary Registering a new login
// @Description Registration is performed using a unique pair of User ID + Card ID.
// @ID addLogin
// @Accept json
// @Produce plain
// @Param login_data body Login true "Login object"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid deserialization"
// @Failure 409 {string} string ""
// @Failure 500 {string} string "read data problem"
// @Router /api/login [post]
func (h *Handler) AddLogin() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("read data problem"))
			return
		}

		// parse message
		var loginData models.Login
		err = json.Unmarshal(reqBody, &loginData)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		err = h.store.AddLogin(r.Context(), user.UserID, loginData)
		if err != nil {
			rw.WriteHeader(http.StatusConflict)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// DeleteLogin delete a login by user and ID
// @Tags Delete
// @Summary Deleting a login
// @Description Deletion is performed using a unique pair of User ID + Card ID.
// @ID delLogin
// @Produce plain
// @Param id path string true "Login ID"
// @Success 200 {string} string "ok"
// @Failure 500 {string} string ""
// @Router /api/login [delete]
func (h *Handler) DeleteLogin() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		loginID := chi.URLParam(r, "id")
		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		// delete login
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

// Text return specific text data
// @Tags Read
// @Summary Returns user text data
// @Description input rest URL, output JSON value
// @ID readText
// @Produce json
// @Param id path string true "Text ID"
// @Success 200 {object} Text
// @Failure 404 {string} string ""
// @Failure 500 {string} string ""
// @Router /api/text/{id} [get]
func (h *Handler) Text() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		textID := chi.URLParam(r, "id")
		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		// Requesting text data information
		targetText, err := h.store.Text(r.Context(), user.UserID, textID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
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

// AddText register a new text by user and ID
// @Tags Add
// @Summary Registering a new text
// @Description Registration is performed using a unique pair of User ID + Card ID.
// @ID addText
// @Accept json
// @Produce plain
// @Param text_data body Text true "Text object"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid deserialization"
// @Failure 409 {string} string ""
// @Failure 500 {string} string "read data problem"
// @Router /api/text [post]
func (h *Handler) AddText() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}

		// parse message
		var textData models.Text
		err = json.Unmarshal(reqBody, &textData)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		err = h.store.AddText(r.Context(), user.UserID, textData)
		if err != nil {
			rw.WriteHeader(http.StatusConflict)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// DeleteText delete text by user and ID
// @Tags Delete
// @Summary Deleting text
// @Description Deletion is performed using a unique pair of User ID + Card ID.
// @ID delText
// @Produce plain
// @Param id path string true "Text ID"
// @Success 200 {string} string "ok"
// @Failure 500 {string} string ""
// @Router /api/text [delete]
func (h *Handler) DeleteText() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		textID := chi.URLParam(r, "id")
		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		// delete text
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

// Binary return specific binary data
// @Tags Read
// @Summary Returns user binary data
// @Description input rest URL, output JSON value
// @ID readBinary
// @Produce json
// @Param id path string true "Binary ID"
// @Success 200 {object} Binary
// @Failure 404 {string} string ""
// @Failure 500 {string} string ""
// @Router /api/text/{id} [get]
func (h *Handler) Binary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		binID := chi.URLParam(r, "id")
		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		// Requesting binary data information
		targetBin, err := h.store.Binary(r.Context(), user.UserID, binID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
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

// AddBinary register binary data by user and ID
// @Tags Add
// @Summary Registering a new set of binary data
// @Description Registration is performed using a unique pair of User ID + Card ID.
// @ID addBinary
// @Accept json
// @Produce plain
// @Param binary_data body Binary true "Binary object"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid deserialization"
// @Failure 409 {string} string ""
// @Failure 500 {string} string "read data problem"
// @Router /api/text [post]
func (h *Handler) AddBinary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("read data problem"))
			return
		}

		// parse message
		var binData models.Binary
		err = json.Unmarshal(reqBody, &binData)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		// Requesting binary data information
		err = h.store.AddBinary(r.Context(), user.UserID, binData)
		if err != nil {
			rw.WriteHeader(http.StatusConflict)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// DeleteBinary delete binary data by user and ID
// @Tags Delete
// @Summary Deleting binary data
// @Description Deletion is performed using a unique pair of User ID + Card ID.
// @ID delBinary
// @Produce plain
// @Param id path string true "Binary ID"
// @Success 200 {string} string "ok"
// @Failure 500 {string} string ""
// @Router /api/text [delete]
func (h *Handler) DeleteBinary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		binID := chi.URLParam(r, "id")
		user := r.Context().Value(auth.UserContextKey{}).(*models.User)

		// delete binary data
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
