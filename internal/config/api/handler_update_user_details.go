package api

import (
	"net/http"
	"time"

	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/config"
	"github.com/rahullpanditaa/http-server/internal/database"
	"github.com/rahullpanditaa/http-server/internal/helpers"
)

// HandlerUpdateUserDetails is the handler function for the endpoint PUT /api/users.
// Via this endpoint, users can update their own email and password.
// Requires an access token (jwt) in header.
// new email, password in request body
func (handler *ApiConfigHandler) HandlerUpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	accessTokenReceivedInHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Authorization header wrong")
		return
	}

	requestBody := helpers.ReadRequestJSON[config.User](w, r)
	newEmailReceievedFromUser := requestBody.Email
	newPasswordReceivedFromUser := requestBody.Password

	userID, err := auth.ValidateJWT(accessTokenReceivedInHeader, handler.Cfg.TokenSecret)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "invalid access token")
		return
	}

	user, err := handler.Cfg.DbQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to retrieve user from db")
		helpers.LogError(err, "unable to retreive user from db")
		return
	}

	hash, err := auth.HashPassword(newPasswordReceivedFromUser)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to hash password")
		helpers.LogError(err, "unable to hash password")
		return
	}

	err = handler.Cfg.DbQueries.UpdateUserDetails(
		r.Context(),
		database.UpdateUserDetailsParams{
			Email:          newEmailReceievedFromUser,
			HashedPassword: hash,
			ID:             user.ID,
		},
	)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "cannot updae user details")
		helpers.LogError(err, "cannot update user details")
		return
	}

	userToReturn := config.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now().UTC(),
		Email:     newEmailReceievedFromUser,
	}

	helpers.RespondWithJson(w, http.StatusOK, userToReturn)

}
