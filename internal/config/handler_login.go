package config

import (
	"log"
	"net/http"
	"time"

	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/handlers"
	"github.com/rahullpanditaa/http-server/internal/handlers/helpers"
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	// accept password and email in request body
	userReceivedInRequest := helpers.ReadRequestJSON[handlers.User](w, r)
	emailReceived := (*userReceivedInRequest).Email
	passwordReceived := (*userReceivedInRequest).Password
	expirationTimeReceived := (*userReceivedInRequest).ExpiresInSeconds

	defaultTokenExpirationTime := int(time.Hour.Seconds())

	// user did not send expires_in_seconds or received time > 1 hour
	if expirationTimeReceived == 0 || expirationTimeReceived > defaultTokenExpirationTime {
		expirationTimeReceived = defaultTokenExpirationTime
	}

	// get user by email
	user, err := cfg.DbQueries.GetUserByEmail(r.Context(), emailReceived)
	if err != nil {
		// w.WriteHeader(http.StatusUnauthorized)
		// w.Write([]byte("incorrect email or passwaord"))
		helpers.RespondWithError(w, 401, "incorrect email or password")
		log.Printf("Error: %v\n", err)
		return
	}

	userHashedPwdStore := user.HashedPassword

	// check if passwords match
	match, err := auth.CheckPasswordHash(passwordReceived, userHashedPwdStore)
	if err != nil || !match {
		helpers.RespondWithError(w, 401, "invalid email or password")
		log.Printf("Error: %v\n", err)
		return
	}

	// create token
	token, err := auth.MakeJWT(user.ID, cfg.JWTToken, time.Duration(expirationTimeReceived))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		return
	}

	userToReturn := handlers.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}

	w.WriteHeader(200)
	helpers.RespondWithJson(w, 200, userToReturn)

}
