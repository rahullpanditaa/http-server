package api

import (
	"net/http"
)

// ApiHandlerHealth is the handler function for the
// pattern "GET /api/healthz"
func (handler *ApiConfigHandler) HandlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "plain/text; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
