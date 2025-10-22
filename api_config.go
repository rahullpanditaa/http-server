package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
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
