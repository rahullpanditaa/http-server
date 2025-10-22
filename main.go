package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
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

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}

func (cfg *apiConfig) numberOfRequests(w http.ResponseWriter, r *http.Request) {
	hits := int(cfg.fileServerHits.Load())
	w.Header().Set("Hits", strconv.Itoa((hits)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(metricsTemplate, hits)))
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits reset\n"))
}

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
	type requestParams struct {
		RequestBody string `json:"body"`
	}

	params := requestParams{}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Printf("Error reading request body: %s", err)
		w.WriteHeader(500)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		log.Printf("Error unmarshalling parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.RequestBody) > 140 {
		type respBody struct {
			Error string `json:"error"`
		}

		resp := respBody{
			Error: "chirp length is greater than 140 chars",
		}
		body, err := json.Marshal(&resp)
		if err != nil {
			log.Printf("Error marshalling json response: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
	} else {
		type respBody struct {
			Valid bool `json:"valid"`
		}
		resp := respBody{
			Valid: true,
		}
		body, err := json.Marshal(&resp)
		if err != nil {
			log.Printf("Error marshalling json response: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	}
}
