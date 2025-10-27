package api

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/helpers"
)

type polkaRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (handler *ApiConfigHandler) HandlerUpdateUserChirpyRed(w http.ResponseWriter, r *http.Request) {
	apiKeyReceived, err := auth.GetAPIKey(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid api key received")
		return
	}
	if apiKeyReceived != handler.Cfg.ApiPolkaKey {
		helpers.RespondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}
	webhookRequestBody := helpers.ReadRequestJSON[polkaRequest](w, r)
	eventReceived := webhookRequestBody.Event
	if eventReceived != "user.upgraded" {
		helpers.RespondWithError(w, http.StatusNoContent, "")
		return
	}

	user, err := handler.Cfg.DbQueries.GetUserByID(r.Context(), webhookRequestBody.Data.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			// user does not exist
			helpers.RespondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to retreive user from db")
		helpers.LogError(err, "unable to retreieve user from db")
		return
	}

	err = handler.Cfg.DbQueries.UpgradeUserToChirpyRed(r.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to upgrade user to chirpy red")
		helpers.LogError(err, "unable to upgrade user to chirpy red")
		return
	}

	helpers.RespondWithJson(w, http.StatusNoContent, "")
}
