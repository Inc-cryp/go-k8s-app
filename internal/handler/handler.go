package handler

import (
	"net/http"

	"github.com/Inc-cryp/go-k8s-app/internal/config"
	"github.com/Inc-cryp/go-k8s-app/internal/response"
)

func HomeHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"message": "Hello from Production Kubernetes 🚀",
			"app":     cfg.AppName,
			"env":     cfg.Environment,
		})
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func VersionHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"app_name":    cfg.AppName,
			"version":     cfg.AppVersion,
			"environment": cfg.Environment,
		})
	}
}

// NotFoundHandler answers any path the service does not serve. It is
// registered last on the root pattern so that unmatched requests get a JSON
// body in the same shape as the rest of the API instead of Go's plain-text
// default 404 page.
func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusNotFound, map[string]string{
			"error": "not found",
			"path":  r.URL.Path,
		})
	}
}
