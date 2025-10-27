package api

import (
	"database/sql"
	"net/http"

	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/helpers"
)

// HandlerRevoke is the handler function for the endpoint POST /api/revoke.
// Retreive refresh token from request header.
// Check whether the refresh token exists in refresh_tokens table in db.
// Revoke refresh token.
// Send back a response with code 204, no body
func (handler *ApiConfigHandler) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	// no req body,
	// yes refresh token in header
	refreshTokenReceived, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Authorization header not found")
		helpers.LogErrorWithRequest(err, r, "Authorization header not found")
		return
	}

	refreshToken, err := handler.Cfg.DbQueries.GetRefreshToken(r.Context(), refreshTokenReceived)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		helpers.LogErrorWithRequest(err, r, "cannot get refresh token from db")
		return
	}

	err = handler.Cfg.DbQueries.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to revoke token")
		helpers.LogErrorWithRequest(err, r, "unable to revoke token")
		return
	}

	helpers.RespondWithJson(w, http.StatusNoContent, "")
}
