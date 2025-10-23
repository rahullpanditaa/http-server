package config

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/rahullpanditaa/http-server/internal"
	"github.com/rahullpanditaa/http-server/internal/database"
	"github.com/rahullpanditaa/http-server/internal/handlers"
	"github.com/rahullpanditaa/http-server/internal/handlers/helpers"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	DbQueries      *database.Queries
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}

func (cfg *ApiConfig) NumberOfRequests(w http.ResponseWriter, r *http.Request) {
	hits := int(cfg.FileServerHits.Load())
	w.Header().Set("Hits", strconv.Itoa((hits)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(internal.MetricsTemplate, hits)))
}

func (cfg *ApiConfig) ResetHits(w http.ResponseWriter, r *http.Request) {
	cfg.FileServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits reset\n"))
}

func (cfg *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// read from request body into user defined struct
	email := helpers.ReadRequestJSON[handlers.User](w, r).Email

	user, err := cfg.DbQueries.CreateUser(r.Context(), email)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Error: %v\n", err)
	}

	helpers.RespondWithJson(w, http.StatusCreated, user)
}
