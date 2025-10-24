package auth

import (
	"testing"

	"github.com/alexedwards/argon2id"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	got, err := HashPassword(password)
	assertError(t, err, nil)

	want, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	assertError(t, err, nil)

	if got != want {
		t.Errorf("got hashed password: %v, want %v\n", got, want)
	}
}

func assertError(t testing.TB, err_got, err_want error) {
	t.Helper()
	if err_got != err_want {
		t.Errorf("got error %v, want %v", err_got, err_want)
	}
}
