// Package auth Предоставляет возможность функционал JWT аутентификации + middleware для его работы
package auth

import (
	"context"
	"net/http"

	"github.com/ncyellow/GophKeeper/internal/server/auth/jwt"
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/storage"
)

// UserContextKey для доступа в контексте данных авторизованного пользователя
type UserContextKey struct{}

// Auth - middleware проверка токена и если все ок проверяем наличие в базе.
func Auth(store storage.Storage, conf *config.Config, parser jwt.Parser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			login, err := parser.ParseToken(authHeader, []byte(conf.SigningKey))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err := store.UserByLogin(r.Context(), login)
			// Если все ок, но по каким-то причинам пользователя в базе нет - тоже не авторизован. Я бы воткнул редис,
			// но чем богаты тем и рады
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r = r.Clone(context.WithValue(r.Context(), UserContextKey{}, user))
			next.ServeHTTP(w, r)
		})
	}
}
