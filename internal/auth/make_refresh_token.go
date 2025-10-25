package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	data := make([]byte, 32)
	rand.Read(data)
	return hex.EncodeToString(data), nil
}
