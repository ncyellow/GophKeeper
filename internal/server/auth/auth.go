// Package auth Предоставляет возможность функционал JWT аутентификации + middleware для его работы
package auth

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

var (
	ErrInvalidLoginPwd = errors.New(`invalid login or password`)
	ErrInvalidToken    = errors.New(`invalid token`)
)

// Claims используем для работы с golang-jwt для подписи токена и его проверки
type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

// Authorizer структура для проверки пользователя и выдачи токена авторизации
type Authorizer struct {
	Store      storage.Storage
	SigningKey []byte
}

// SignIn - проверяет есть ли такой пользователь в базе и генерируем токен
func (a *Authorizer) SignIn(ctx context.Context, user *models.User) (string, error) {
	expirationDuration := time.Hour * 24
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	repoUser, err := a.Store.User(ctx, user.Login, user.Password)
	if err != nil {
		return "", ErrInvalidLoginPwd
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

// ParseToken - проверяем является ли токен корректным, если да возвращаем логин данного пользователя
func ParseToken(accessToken string, signingKey []byte) (string, error) {
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

// HashPWD - хеширования пароля
func HashPWD(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
