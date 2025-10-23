package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rahullpanditaa/http-server/internal/config"
	"github.com/rahullpanditaa/http-server/internal/database"
	"github.com/rahullpanditaa/http-server/internal/handlers"
)

var apiCfg *config.ApiConfig

func main() {
	dbQueries := connectToDb()

	apiCfg.DbQueries = dbQueries

	mux := http.NewServeMux()

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(fileServerHandler))

	// api endpoint
	mux.HandleFunc("GET /api/healthz", handlers.ApiHandlerHealth)
	mux.HandleFunc("POST /api/validate_chirp", handlers.ApiValidateChirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.CreateUserHandler)
	// admin endpoint
	mux.HandleFunc("GET /admin/metrics", apiCfg.NumberOfRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetHits)

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
