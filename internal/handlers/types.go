package handlers

import (
	"time"

	"github.com/google/uuid"
)

type RequestParams struct {
	RequestBody string `json:"body"`
}

type ProfanityLessChirp struct {
	CleanedBody string `json:"cleaned_body"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
