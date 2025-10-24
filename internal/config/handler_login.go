package config

import (
	"log"
	"net/http"

	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/handlers"
	"github.com/rahullpanditaa/http-server/internal/handlers/helpers"
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	// accept password and email in request body
	userSentInRequest := helpers.ReadRequestJSON[handlers.User](w, r)
	emailSent := (*userSentInRequest).Email
	passwordSent := (*userSentInRequest).Password

	// get user by email
	user, err := cfg.DbQueries.GetUserByEmail(r.Context(), emailSent)
	if err != nil {
		// w.WriteHeader(http.StatusUnauthorized)
		// w.Write([]byte("incorrect email or passwaord"))
		helpers.RespondWithError(w, 401, "incorrect email or password")
		log.Printf("Error: %v\n", err)
		return
	}

	userHashedPwdStore := user.HashedPassword

	// check if passwords match
	match, err := auth.CheckPasswordHash(passwordSent, userHashedPwdStore)
	if err != nil || !match {
		helpers.RespondWithError(w, 401, "invalid email or password")
		log.Printf("Error: %v\n", err)
		return
	}

	userToReturn := handlers.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	w.WriteHeader(200)
	helpers.RespondWithJson(w, 200, userToReturn)

}
