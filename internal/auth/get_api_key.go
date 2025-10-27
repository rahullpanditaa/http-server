package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	v := headers.Get("Authorization")
	if v == "" {
		return "", fmt.Errorf("authorization header not found")
	}
	parts := strings.Fields(v)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "ApiKey") {
		return "", fmt.Errorf("invalid authorization header format")
	}
	apiKey := strings.TrimSpace(parts[1])
	if apiKey == "" {
		return "", fmt.Errorf("empty api key")
	}
	return apiKey, nil
}
