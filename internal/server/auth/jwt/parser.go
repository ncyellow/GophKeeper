package jwt

type Parser interface {
	ParseToken(accessToken string, signingKey []byte) (string, error)
}
