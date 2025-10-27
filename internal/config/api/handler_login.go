package api

import (
	"net/http"
	"time"

	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/config"
	"github.com/rahullpanditaa/http-server/internal/database"
	"github.com/rahullpanditaa/http-server/internal/helpers"
)

// HandlerLogin is the handler function for the endpoint POST /api/login.
// Receives email, password in request body.
// Retrieves the user from users table in chirpy database using email received.
// Check if password received in request matches the hashed password stored for the user in table.
// Create a JWT (access token) for user with expiration time of 1 hour.
// Create a new refresh token for the user, store it in refresh_tokens table along with the user ID.
// Send back a JSON response which includes the JWT and the refresh token.
func (handler *ApiConfigHandler) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	// accept password and email in request body
	userReceivedInRequest := helpers.ReadRequestJSON[config.User](w, r)
	emailReceived := userReceivedInRequest.Email
	passwordReceived := userReceivedInRequest.Password

	// get user by email
	user, err := handler.Cfg.DbQueries.GetUserByEmail(r.Context(), emailReceived)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "incorrect email or password")
		helpers.LogErrorWithRequest(err, r, "incorrect email or password")
		return
	}

	userHashedPswdInDb := user.HashedPassword

	// check if passwords match
	match, err := auth.CheckPasswordHash(passwordReceived, userHashedPswdInDb)
	if err != nil || !match {
		helpers.RespondWithError(w, http.StatusUnauthorized, "invalid email or password")
		helpers.LogErrorWithRequest(err, r, "invalid email or password")
		return
	}

	// create jwt access token
	// expiration time for jwt - 1 hour
	JWTExpirationTime := int(time.Hour.Seconds())
	token, err := auth.MakeJWT(user.ID, handler.Cfg.TokenSecret, time.Duration(JWTExpirationTime)*time.Second)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to create a JWT")
		helpers.LogErrorWithRequest(err, r, "unable to create a JWT")
		return
	}

	// create refresh token
	refreshTokenStr, _ := auth.MakeRefreshToken()
	refreshToken, err := handler.Cfg.DbQueries.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:  refreshTokenStr,
			UserID: user.ID,
		},
	)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "error creating refresh token")
		helpers.LogErrorWithRequest(err, r, "error creating refresh token")
		return
	}

	userToReturn := config.User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	}

	helpers.RespondWithJson(w, 200, userToReturn)
}
