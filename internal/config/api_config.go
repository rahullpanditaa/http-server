package config

import (
	"net/http"
	"sync/atomic"

	"github.com/rahullpanditaa/http-server/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
	TokenSecret    string
	ApiPolkaKey    string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}
