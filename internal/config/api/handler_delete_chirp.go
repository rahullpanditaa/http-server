package api

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/helpers"
)

func (handler *ApiConfigHandler) HandlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "could not get chirp ID from path")
		return
	}

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "cannot parse given chirpID into a uuid")
		helpers.LogError(err, "cannot parse given chirpID into a uuid")
		return
	}

	accessTokenReceivedInHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Authorization header wrong or missing")
		return
	}

	userID, err := auth.ValidateJWT(accessTokenReceivedInHeader, handler.Cfg.TokenSecret)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "JWT invalid")
		return
	}

	chirp, err := handler.Cfg.DbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondWithError(w, http.StatusNotFound, "invalid chirp ID, chirp doesn't exist")
			return
		}
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to retreive chirp from db")
		helpers.LogError(err, "unable to retreive chirp from db")
		return
	}

	if chirp.UserID != userID {
		helpers.RespondWithError(w, http.StatusForbidden, "you are not the author of this chirp")
		return
	}

	err = handler.Cfg.DbQueries.DeleteChirpByID(r.Context(), chirp.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to delete chrip")
		helpers.LogError(err, "unable to delete chirp")
		return
	}

	helpers.RespondWithJson(w, http.StatusNoContent, "")

}
