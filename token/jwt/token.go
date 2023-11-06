package jwt

import (
	"github.com/deividaspetraitis/wallet-screener/errors"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenNotValid = errors.New("Token is not valid")
)

// Token is representation of JWT token.
type Token struct {
	*jwt.Token
	Claims *Claims
}

// Claims is JWT claims.
type Claims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

// Parse parses given JWT token string, validates and return a Token.
// Parse does NOT check for signature validity.
func Parse(tkn string) (*Token, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tkn, &Claims{})
	if err != nil {
		return nil, err
	}

	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &Token{
		Token:  token,
		Claims: claims,
	}, nil
}
