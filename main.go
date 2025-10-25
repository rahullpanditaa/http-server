package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rahullpanditaa/http-server/internal/config"
	"github.com/rahullpanditaa/http-server/internal/config/admin"
	"github.com/rahullpanditaa/http-server/internal/config/api"
	"github.com/rahullpanditaa/http-server/internal/database"
)

func main() {
	// var apiCfg config.ApiConfig
	// apiCfg = make(config.ApiConfig)

	dbQueries := connectToDb()
	jwtSecretToken := os.Getenv("TOKEN_SECRET")

	cfg := &config.ApiConfig{
		DbQueries:   dbQueries,
		TokenSecret: jwtSecretToken,
	}

	apiCfgHandler := &api.ApiConfigHandler{
		Cfg: cfg,
	}

	adminCfgHandler := &admin.ApiConfigHandler{
		Cfg: cfg,
	}

	mux := http.NewServeMux()

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.MiddlewareMetricsInc(fileServerHandler))

	// api endpoint
	mux.HandleFunc("GET /api/healthz", apiCfgHandler.HandlerHealth)
	mux.HandleFunc("POST /api/users", apiCfgHandler.HandlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfgHandler.HandlerLogin)
	mux.HandleFunc("POST /api/chirps", apiCfgHandler.HandlerValidateChirps)
	mux.HandleFunc("POST /api/refresh", apiCfgHandler.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfgHandler.HandlerRevoke)

	mux.HandleFunc("GET /api/chirps", apiCfgHandler.HandlerReturnAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfgHandler.HandlerReturnChirpByID)
	// admin endpoint
	mux.HandleFunc("GET /admin/metrics", adminCfgHandler.HandlerNumberOfRequests)
	mux.HandleFunc("POST /admin/reset", adminCfgHandler.HandlerResetHits)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func connectToDb() *database.Queries {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	return database.New(db)
}
