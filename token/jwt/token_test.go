package jwt

import (
	"testing"

	"github.com/deividaspetraitis/wallet-screener/errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	var testcases = []struct {
		tokenstring string
		claims      *Claims
		err         error
	}{
		// should fail: token is expired
		{
			tokenstring: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJoaUBkZWl2aWRhc3BldHJhaXRpcy5sdCIsImV4cCI6MTY5NjI1MTY2OH0.n4lEyGjTBJ5pmnSlX7lT3a-OljjxzRwg_3s2q7U1Vuk",
			err:         jwt.ErrTokenExpired,
		},
		// should pass: token is not expired
		{
			tokenstring: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJoaUBkZWl2aWRhc3BldHJhaXRpcy5sdCIsImV4cGlyZXMiOjQwOTQ2ODY5Mjd9.nDThjOeZYEYeiqD2x0cXtclEVohiJcOjEacwaur3j-c",
			claims: &Claims{
				Sub: "hi@deividaspetraitis.lt",
			},
			err: nil,
		},
	}

	for _, tt := range testcases {
		token, err := Parse(tt.tokenstring)
		if !errors.Is(err, tt.err) {
			t.Fatalf("token %s got %v, want %v", tt.tokenstring, err, tt.err)
		}

		if tt.claims == nil {
			continue
		}

		if !cmp.Equal(token.Claims, tt.claims) {
			t.Errorf("token %s claims got %v, want %v", tt.tokenstring, token.Claims, tt.claims)
		}
	}
}
