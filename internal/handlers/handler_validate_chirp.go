package handlers

import (
	"net/http"
	"strings"

	"github.com/rahullpanditaa/http-server/internal/handlers/helpers"
)

// ApiValidateChirpHandler is the handler function
// for the endpoint "GET /api/validate_chirp"
func ApiValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	requestBody := helpers.ReadRequestJSON[RequestParams](w, r).RequestBody

	if len(requestBody) > 140 {
		helpers.RespondWithError(w, http.StatusBadRequest, "chirp length greater than 140 chars")
	} else {
		req_words := strings.Split(requestBody, " ")
		cleaned := checkForProfanity(req_words)
		resp := ProfanityLessChirp{CleanedBody: cleaned}
		helpers.RespondWithJson(w, http.StatusOK, resp)

	}
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
