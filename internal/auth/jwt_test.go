package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	t.Run("jwt generation -> validation cycle", func(t *testing.T) {
		want := uuid.New()
		tokenString, err := MakeJWT(want, "secret69", 2*time.Hour)
		assertError(t, err, nil)

		got, err := ValidateJWT(tokenString, "secret69")
		assertError(t, err, nil)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("wrong secret used for validation", func(t *testing.T) {
		userId := uuid.New()
		tokenString, err := MakeJWT(userId, "ssshhh-secret!!", time.Hour)
		assertError(t, err, nil)

		_, err = ValidateJWT(tokenString, "secret-key")
		assertError(t, err, jwt.ErrTokenSignatureInvalid)

	})
	t.Run("token expired", func(t *testing.T) {
		userID := uuid.New()
		tokenString, err := MakeJWT(userID, "sssshhhh", time.Millisecond)
		assertError(t, err, nil)

		time.Sleep(2 * time.Millisecond)
		_, err = ValidateJWT(tokenString, "sssshhhh")
		assertError(t, err, jwt.ErrTokenExpired)
	})

}
