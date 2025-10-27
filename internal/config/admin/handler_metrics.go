package admin

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/rahullpanditaa/http-server/internal"
	"github.com/rahullpanditaa/http-server/internal/config"
	"github.com/rahullpanditaa/http-server/internal/helpers"
)

type ApiConfigHandler struct {
	Cfg *config.ApiConfig
}

// HandlerNumberOfRequests is the handler function for the endpoint GET /admin/metrics.
func (handler *ApiConfigHandler) HandlerNumberOfRequests(w http.ResponseWriter, r *http.Request) {
	hits := int(handler.Cfg.FileServerHits.Load())
	w.Header().Set("Hits", strconv.Itoa((hits)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(internal.MetricsTemplate, hits)))
}

// post admin/reset
// HandlerResetHits is the handler function for the endpoint POST /admin/reset.
// Deletes all users in db.
func (handler *ApiConfigHandler) HandlerResetHits(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	handler.Cfg.Platform = platform
	if handler.Cfg.Platform != "dev" {
		helpers.RespondWithError(w, http.StatusForbidden, "endpoint access only available in a local dev environment")
		return
	}
	err := handler.Cfg.DbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to delete all users")
		helpers.LogErrorWithRequest(err, r, "unable to delete all users")
		return
	}
}
