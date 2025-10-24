package auth

import (
	"testing"

	"github.com/alexedwards/argon2id"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	got, err := HashPassword(password)
	assertError(t, err, nil)

	match, err := argon2id.ComparePasswordAndHash(password, got)
	assertError(t, err, nil)

	if !match {
		t.Errorf("expected input password %v to match hash %v", password, got)
	}
}

func TestCheckHashPassword(t *testing.T) {
	pwd := "testpassword456789"
	hash, err := argon2id.CreateHash(pwd, argon2id.DefaultParams)
	assertError(t, err, nil)

	match, err := CheckPasswordHash(pwd, hash)
	assertError(t, err, nil)

	if !match {
		t.Errorf("match fail: input password %v didn't match hash %v", pwd, hash)
	}
}

func assertError(t testing.TB, err_got, err_want error) {
	t.Helper()
	if err_got != err_want {
		t.Errorf("got error %v, want %v", err_got, err_want)
	}
}
