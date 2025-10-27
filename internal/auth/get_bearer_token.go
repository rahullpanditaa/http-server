package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(h http.Header) (string, error) {
	v := h.Get("Authorization")
	if v == "" {
		return "", fmt.Errorf("authorization header not found")
	}
	parts := strings.Fields(v)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid authorization header format")
	}
	tok := strings.TrimSpace(parts[1])
	if tok == "" {
		return "", fmt.Errorf("empty bearer token")
	}
	return tok, nil
}
