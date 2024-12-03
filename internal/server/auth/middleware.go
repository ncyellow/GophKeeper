// Package auth Provides functionality for JWT authentication + middleware to support its operation
package auth

import (
	"context"
	"net/http"

	"github.com/ncyellow/GophKeeper/internal/server/auth/jwt"
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/storage"
)

// UserContextKey for accessing authorized user data in context
type UserContextKey struct{}

// Auth - middleware checks the token and if all is well, verifies its presence in the database.
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
			// If all is well, but for some reason the user is not in the database - also not authorized.
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r = r.Clone(context.WithValue(r.Context(), UserContextKey{}, user))
			next.ServeHTTP(w, r)
		})
	}
}
