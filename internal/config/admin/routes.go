package admin

import "net/http"

func (handler *ApiConfigHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /admin/metrics", handler.HandlerNumberOfRequests)
	mux.HandleFunc("POST /admin/reset", handler.HandlerResetHits)
}
