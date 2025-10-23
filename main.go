package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rahullpanditaa/http-server/internal/database"
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
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	// open connection to db
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	var apiCfg apiConfig

	dbQueries := database.New(db)
	apiCfg.dbQueries = dbQueries

	mux := http.NewServeMux()

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

	err = server.ListenAndServe()
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
