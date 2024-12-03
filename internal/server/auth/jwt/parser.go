package jwt

// Parser interface, which we use to verify the correctness of the jwt token.
// It is needed for testing authorization through gomock
type Parser interface {
	ParseToken(accessToken string, signingKey []byte) (string, error)
}
