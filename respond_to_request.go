package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type respError struct {
		Error string `json:"error"`
	}

	resp := respError{Error: msg}
	body, err := json.Marshal(&resp)
	assertError(err, nil, w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

type responseStruct struct {
	CleanedBody string `json:"cleaned_body"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	body, err := json.Marshal(&payload)
	assertError(err, nil, w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
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

func assertError(got, want error, w http.ResponseWriter) {
	if got != want {
		log.Printf("Error: %s\n", got)
		w.WriteHeader(500)
		return
	}
}
