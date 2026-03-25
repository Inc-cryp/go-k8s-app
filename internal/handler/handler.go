package handler

import (
	"net/http"

	"github.com/Inc-cryp/go-k8s-app/internal/response"
	"github.com/Inc-cryp/go-k8s-app/internal/config"

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
