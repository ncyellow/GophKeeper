// Package jwt provides functionality for JWT authentication
package jwt

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ncyellow/GophKeeper/internal/models"
	"github.com/ncyellow/GophKeeper/internal/server/storage"
)

var ErrInvalidToken = errors.New(`invalid token`)

// Claims used for working with golang-jwt to sign the token and verify it
type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

// Authorizer structure for verifying the user and issuing the authorization token
type Authorizer struct {
	Store      storage.Storage
	SigningKey []byte
}

// DefaultParser default parser for jwt tokens
type DefaultParser struct{}

// SignIn - checks if such a user exists in the database and generates a token
func (a *Authorizer) SignIn(ctx context.Context, user *models.User) (string, error) {
	expirationDuration := time.Hour * 24
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	repoUser, err := a.Store.User(ctx, user.Login, user.Password)
	if err != nil {
		return "", fmt.Errorf("invalid login or password: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Username: repoUser.Login,
	})
	return token.SignedString(a.SigningKey)
}

// ParseToken - checks if the token is valid, if yes returns the login of this user
func (p *DefaultParser) ParseToken(accessToken string, signingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing error")
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, nil
	}

	return "", ErrInvalidToken
}
