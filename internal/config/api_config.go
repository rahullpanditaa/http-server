package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/rahullpanditaa/http-server/internal"
	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/database"
	"github.com/rahullpanditaa/http-server/internal/handlers"
	"github.com/rahullpanditaa/http-server/internal/handlers/helpers"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
	TokenSecret    string
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
	// user will now also send a plain text password
	// in json request along with email
	userDetailsSent := helpers.ReadRequestJSON[handlers.User](w, r)

	emailInRequest := (*userDetailsSent).Email
	passwordInRequest := (*userDetailsSent).Password

	// hash the password
	hashedPassword, err := auth.HashPassword(passwordInRequest)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error: %v\n", err)
		return
	}

	user, err := cfg.DbQueries.CreateUser(r.Context(),
		database.CreateUserParams{
			Email:          emailInRequest,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Error: %v\n", err)
	}

	userToReturn := handlers.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email:     user.Email,
		Password:  passwordInRequest,
	}

	helpers.RespondWithJson(w, http.StatusCreated, userToReturn)
}

func (cfg *ApiConfig) ValidateChirpsHandler(w http.ResponseWriter, r *http.Request) {
	// read from request body
	requestHeaders := r.Header
	requestPayload := helpers.ReadRequestJSON[handlers.RequestParams](w, r)
	requestBody := requestPayload.RequestBody
	// userID := requestPayload.UserID

	// jwt token sent by user in request
	// this needs to be validated
	userTokenStringReceived, err := auth.GetBearerToken(requestHeaders)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error: %v\n", err)
		return
	}

	userIDFromJWT, err := auth.ValidateJWT(userTokenStringReceived, cfg.TokenSecret)
	if err != nil {
		if errors.Is(err, errors.New("invalid token")) {
			// w.WriteHeader(http.StatusUnauthorized)
			helpers.RespondWithError(w, http.StatusUnauthorized, "JWT invalid")
			return
		}
		w.WriteHeader(500)
		log.Printf("Error: %v\n", err)
		return
	}

	if len(requestBody) > 140 {
		helpers.RespondWithError(w, http.StatusBadRequest, "chirp length greater than 140 chars")
	} else {
		req_words := strings.Split(requestBody, " ")
		cleaned := checkForProfanity(req_words)

		chirp, err := cfg.DbQueries.CreateChirp(
			r.Context(),
			database.CreateChirpParams{
				Body:   cleaned,
				UserID: userIDFromJWT,
			},
		)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("Error: %v\n", err)
			return
		}

		chirpResource := handlers.Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    userIDFromJWT,
		}

		helpers.RespondWithJson(w, http.StatusCreated, chirpResource)

	}
}

func (cfg *ApiConfig) HandlerReturnAllChirps(w http.ResponseWriter, r *http.Request) {
	// get all chirps from table
	allChirps, err := cfg.DbQueries.GetChirps(r.Context())
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Error: %v\n", err)
	}

	var chirpsToReturn []handlers.Chirp
	for _, chirp := range allChirps {
		c := handlers.Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		chirpsToReturn = append(chirpsToReturn, c)
	}

	helpers.RespondWithJson(w, http.StatusOK, chirpsToReturn)
}

func (cfg *ApiConfig) HandlerReturnChirpByID(w http.ResponseWriter, r *http.Request) {
	// get chirp from table by id

	// get user id from path
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
		w.WriteHeader(500)
		log.Fatalln("Error: could not get user ID from path")
	}

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Error: %v\n", err)
	}

	chirp, err := cfg.DbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			log.Printf("Error: %v\n", err)
			return
			// log.Fatalf("Error: %v\n", err)
		}
		w.WriteHeader(500)
		log.Printf("Error: %v\n", err)
		return
	}
	chirpToReturn := handlers.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	helpers.RespondWithJson(w, http.StatusOK, chirpToReturn)

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
