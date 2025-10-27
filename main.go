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
	godotenv.Load()
	dbQueries := connectToDb()
	jwtSecretToken := os.Getenv("TOKEN_SECRET")
	platform := os.Getenv("PLATFORM")
	polkaApiKey := os.Getenv("POLKA_KEY")

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	cfg := &config.ApiConfig{
		DbQueries:   dbQueries,
		TokenSecret: jwtSecretToken,
		Platform:    platform,
		ApiPolkaKey: polkaApiKey,
	}
	apiCfgHandler := &api.ApiConfigHandler{
		Cfg: cfg,
	}
	adminCfgHandler := &admin.ApiConfigHandler{
		Cfg: cfg,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(fileServerHandler))

	apiCfgHandler.RegisterRoutes(mux)
	adminCfgHandler.RegisterRoutes(mux)

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
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	return database.New(db)
}
