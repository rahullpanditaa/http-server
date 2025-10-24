package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// A JWT essentially has 3 parts
// header - metadata about the token (algo used, type of token)
// Claims / payload - data to store in the token
// signature - ensures token isn't tampered with
// token secret - key used to sign the token and match against
// during validation of jwt

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// create a json web token
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   userID.String(),
		},
	)
	// key type for HS256 -> []byte
	jwtString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return jwtString, nil
}
