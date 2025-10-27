package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/config"
	"github.com/rahullpanditaa/http-server/internal/database"
	"github.com/rahullpanditaa/http-server/internal/helpers"
)

// HandlerValidateChirps is the handler function for the endpoint POST /api/chirps.
// Read chirp contents in request body, headers.
// Read the JWT sent by user in request headers.
// Validate the JWT sent by the user.
// Validate chirp - if length of text is greater than 140 chars, respond with an Error.
// Create a new chirp in the chirps table in chirpy database.
// Return the details of chirp created in a JSON response.
func (handler *ApiConfigHandler) HandlerValidateChirps(w http.ResponseWriter, r *http.Request) {
	// requestPayload := helpers.[handlers.RequestParams](w, r)
	requestPayload := helpers.ReadRequestJSON[config.RequestParams](w, r)
	requestBody := requestPayload.RequestBody
	requestHeaders := r.Header

	// jwt token sent by user in request
	userTokenStringReceived, err := auth.GetBearerToken(requestHeaders)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid header sent")
		helpers.LogErrorWithRequest(err, r, "invalid header in request")
		return
	}

	// validate jwt sent by user
	userIDFromJWT, err := auth.ValidateJWT(userTokenStringReceived, handler.Cfg.TokenSecret)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "JWT invalid")
		// helpers.LogErrorWithRequest(err, r, "JWT invalid")
		return
		// if errors.Is(err, auth.ErrInvalidToken) {
		// 	helpers.RespondWithError(w, http.StatusUnauthorized, "JWT invalid")
		// 	helpers.LogErrorWithRequest(err, r, "JWT invalid")
		// 	return
		// }
		// helpers.RespondWithError(w, 500, "")
		// helpers.LogErrorWithRequest(err, r, "error occurred while validating JWT")
		// return
	}

	if len(requestBody) > 140 {
		helpers.RespondWithError(w, http.StatusBadRequest, "chirp length greater than 140 chars")
		return
	}

	req_words := strings.Split(requestBody, " ")
	cleaned := checkForProfanity(req_words)

	chirp, err := handler.Cfg.DbQueries.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body:   cleaned,
			UserID: userIDFromJWT,
		},
	)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "cannot create a chir")
		helpers.LogErrorWithRequest(err, r, "cannot create a chirp")
		return
	}

	chirpResource := config.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    userIDFromJWT,
	}

	helpers.RespondWithJson(w, http.StatusCreated, chirpResource)

}

// HandlerReturnAllChirps is the handler function for the endpoint GET /api/chirps.
// Retreives all the chirps from chirps table in db.
// Sends back a JSON response with a slice of all the chirps in db.
func (handler *ApiConfigHandler) HandlerReturnAllChirps(w http.ResponseWriter, r *http.Request) {
	// get all chirps from table
	allChirps, err := handler.Cfg.DbQueries.GetChirps(r.Context())
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "cannot retreive chirps from db")
		helpers.LogErrorWithRequest(err, r, "cannot retreive chirps from db")
		return
	}

	var chirpsToReturn []config.Chirp
	for _, chirp := range allChirps {
		c := config.Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		chirpsToReturn = append(chirpsToReturn, c)
	}

	helpers.RespondWithJson(w, http.StatusOK, chirpsToReturn)
}

// get api/chirps/{chirpId}

// HandlerReturnChirpByID is the handler function for the endpoint GET /api/chirps/{chirpID}.
// Get chirpID from request query parameter.
// Get the chirp from the chirps table in db.
// Send back a JSON response with details of the chirp.
func (handler *ApiConfigHandler) HandlerReturnChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "could not get chirp ID from path")
		return
	}

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "cannot parse given chirpID into a uuid")
		helpers.LogErrorWithRequest(err, r, "cannot parse given chirpID into a uuid")
		return
	}

	chirp, err := handler.Cfg.DbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondWithError(w, http.StatusNotFound, "invalid chirp ID")
			return
		}
		helpers.RespondWithError(w, http.StatusInternalServerError, "cannot retreive chirp from db")
		helpers.LogErrorWithRequest(err, r, "cannot retreive chirp from db")
		return
	}

	chirpToReturn := config.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	helpers.RespondWithJson(w, http.StatusOK, chirpToReturn)
}

func checkForProfanity(sentence []string) string {
	var sanitizedSentence []string

	for _, word := range sentence {
		w := strings.ToLower(word)
		switch w {
		case "kerfuffle", "sharbert", "fornax":
			word = "****"
		}
		sanitizedSentence = append(sanitizedSentence, word)

	}

	return strings.Join(sanitizedSentence, " ")
}
