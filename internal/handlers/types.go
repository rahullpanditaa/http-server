package handlers

import (
	"time"

	"github.com/google/uuid"
)

type RequestParams struct {
	RequestBody string    `json:"body"`
	UserID      uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type User struct {
	ID               uuid.UUID `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	ExpiresInSeconds int       `json:"expires_in_seconds"`
	Token            string    `json:"token"`
	// HashedPassword string    `json:"hashed_password"`
}
