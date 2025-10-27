package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	errResp := ErrorResponse{Error: msg}

	// marshal ((decode)) error struct into a slice of bytes
	body, err := json.Marshal(&errResp)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error: %v\n", err)
		return
	}

	// send a http response of json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	// payload - the type (most likely struct) to be marshalled into a slice of bytes
	// and sent as http response
	body, err := json.Marshal(&payload)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error: %v\n", err)
		return
	}

	// send http response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}
