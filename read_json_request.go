package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type requestParams struct {
	RequestBody string `json:"body"`
}

func readJSONRequest(w http.ResponseWriter, r *http.Request) *requestParams {
	params := requestParams{}
	body, err := io.ReadAll(r.Body)
	assertError(err, nil, w)
	defer r.Body.Close()

	err = json.Unmarshal(body, &params)
	assertError(err, nil, w)

	return &params
}
