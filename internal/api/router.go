package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/adesokandijam/go-rest-api-boilerplate/internal/api/middleware"
	"github.com/adesokandijam/go-rest-api-boilerplate/internal/api/respond"
)

func NewHandler(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		respond.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/v1/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		respond.JSON(w, http.StatusOK, map[string]string{"message": "Hello, " + r.PathValue("name") + "!"})
	})
	return middleware.RequestID(middleware.Logging(logger, middleware.Recover(logger, middleware.Timeout(15*time.Second, mux))))
}
