package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	Platform       string
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

// Delete all users in DB
func (cfg *ApiConfig) ResetHits(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	cfg.Platform = platform
	if cfg.Platform != "dev" {
		w.WriteHeader(403)
		log.Fatal("can only access this endpoint in a local dev environment")
	}
	err := cfg.DbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Error: %v\n", err)
	}
}

func (cfg *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// read from request body into user defined struct
	email := helpers.ReadRequestJSON[handlers.User](w, r).Email

	user, err := cfg.DbQueries.CreateUser(r.Context(), email)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Error: %v\n", err)
	}

	userToReturn := handlers.User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.CreatedAt, Email: user.Email}

	helpers.RespondWithJson(w, http.StatusCreated, userToReturn)
}
