package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// token string -> jwtString from make_jwt
	// token secret

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}
