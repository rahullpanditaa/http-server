package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/handlers/helpers"
	"github.com/rahullpanditaa/http-server/internal/helpers_temp"
)

// HandlerRefresh is the handler function for the endpoint POST /api/refresh.
// It receives no body in the request, only headers.
// Retreive refresh token from headers.
// Check whether received refresh token exists in db.
// Create a new JWT.
// Send the JWT as a JSON response.
func (handler *ApiConfigHandler) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenReceived, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Authorization header not found")
		helpers_temp.LogErrorWithRequest(err, r, "Authorization header not found")
		return
	}

	refreshToken, err := handler.cfg.DbQueries.GetRefreshToken(r.Context(), refreshTokenReceived)
	if err != nil {
		// check if err is because no rows returned
		if err == sql.ErrNoRows {
			helpers.RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		helpers_temp.LogErrorWithRequest(err, r, "cannot get refresh token from db")
		return
	}

	userID := refreshToken.UserID

	// create new jwt
	jwt, err := auth.MakeJWT(userID, handler.cfg.TokenSecret, time.Hour)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "cannot create a JWT")
		helpers_temp.LogErrorWithRequest(err, r, "cannot create a JWT")
		return
	}

	responseToSend := struct {
		Token string `json:"token"`
	}{
		Token: jwt,
	}

	helpers.RespondWithJson(w, http.StatusOK, responseToSend)
}
