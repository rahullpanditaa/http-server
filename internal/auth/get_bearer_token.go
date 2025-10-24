package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	// in request -> header: Authorization, Bearer TOKEN_STRING
	authHeaderInfo := headers.Get("Authorization")
	if authHeaderInfo == "" {
		return "", fmt.Errorf("authorization header does not exist")
	}
	tokenString := strings.TrimSpace(strings.Replace(authHeaderInfo, "Bearer", "", 1))
	return tokenString, nil
}
