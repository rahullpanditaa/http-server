package config

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/handlers/helpers"
)

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	// does not accept a request body
	// does require a refresh token to be present in headers
	refreshTokenReceived, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		log.Printf("Error: %v\n", err)
		return
	}

	// check refresh token in db
	refreshToken, err := cfg.DbQueries.GetRefreshToken(r.Context(), refreshTokenReceived)
	if err != nil {
		// check if err is because no rows returned
		if err == sql.ErrNoRows {
			helpers.RespondWithError(w, 401, "invalid refresh token")
		}
		log.Printf("Error: %v\n", err)
		return
	}

	userID := refreshToken.UserID

	// create new jwt
	jwt, err := auth.MakeJWT(userID, cfg.TokenSecret, time.Hour)
	if err != nil {
		helpers.RespondWithError(w, 500, err.Error())
		log.Printf("Error: %v\n", err)
		return
	}

	// do i really need to create an entire user struct?
	// let's find out

	responseToSend := struct {
		token string
	}{
		token: jwt,
	}

	helpers.RespondWithJson(w, 200, responseToSend)
}
