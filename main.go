package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	metricsTemplate = `
	<html>
  		<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
)

func main() {
	mux := http.NewServeMux()
	var apiCfg apiConfig

	readinessHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))

	// api endpoint
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	// admin endpoint
	mux.HandleFunc("GET /admin/metrics", apiCfg.numberOfRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

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

func assertError(got, want error, w http.ResponseWriter) {
	if got != want {
		log.Printf("Error: %s\n", got)
		w.WriteHeader(500)
		return
	}
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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	requestBody := readJSONRequest(w, r).RequestBody

	if len(requestBody) > 140 {
		respondWithError(w, http.StatusBadRequest, "chirp length greater than 140 chars")
	} else {
		req_words := strings.Split(requestBody, " ")
		cleaned := checkForProfanity(req_words)
		resp := responseStruct{CleanedBody: cleaned}
		respondWithJSON(w, http.StatusOK, resp)

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
