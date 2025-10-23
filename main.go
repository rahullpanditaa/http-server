package main

import (
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
