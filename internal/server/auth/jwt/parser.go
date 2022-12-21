package jwt

// Parser интерфейс, который мы используем для проверки jwt токена на корректность.
// Нужен нам для того, чтобы тестировать авторизацию через gomock
type Parser interface {
	ParseToken(accessToken string, signingKey []byte) (string, error)
}
