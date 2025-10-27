package api

import "net/http"

func (handler *ApiConfigHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/healthz", handler.HandlerHealth)
	mux.HandleFunc("POST /api/users", handler.HandlerCreateUser)
	mux.HandleFunc("POST /api/login", handler.HandlerLogin)
	mux.HandleFunc("POST /api/chirps", handler.HandlerValidateChirps)
	mux.HandleFunc("POST /api/refresh", handler.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", handler.HandlerRevoke)
	mux.HandleFunc("GET /api/chirps", handler.HandlerReturnAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", handler.HandlerReturnChirpByID)

	mux.HandleFunc("PUT /api/users", handler.HandlerUpdateUserDetails)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", handler.HandlerDeleteChirpByID)
}
