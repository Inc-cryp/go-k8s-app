// Command server runs the go-k8s-app HTTP service.
//
// It serves the JSON API documented in the README and shuts down gracefully
// when it receives SIGINT or SIGTERM. The latter matters under Kubernetes:
// the kubelet sends SIGTERM and then waits for the pod to exit before it
// escalates to SIGKILL, so the process has to stop accepting new connections
// and finish the in-flight requests inside that window.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Inc-cryp/go-k8s-app/internal/config"
	"github.com/Inc-cryp/go-k8s-app/internal/handler"
)

const (
	// defaultAddr is the listen address used when ADDR is unset. The
	// container exposes 8080 and the Kubernetes probes target it.
	defaultAddr = ":8080"

	// shutdownTimeout bounds how long in-flight requests get to finish.
	// Keep it comfortably below the pod's terminationGracePeriodSeconds.
	shutdownTimeout = 10 * time.Second

	// readHeaderTimeout bounds how long a client may take to send request
	// headers, which protects the server from slowloris-style connections.
	readHeaderTimeout = 5 * time.Second

	// idleTimeout reclaims keep-alive connections that have gone quiet.
	idleTimeout = 60 * time.Second
)

// newServer builds the HTTP server with every route the service exposes.
func newServer(cfg config.Config) *http.Server {
	mux := http.NewServeMux()
	// "{$}" matches only "/" exactly. A bare "/" pattern would act as a
	// catch-all and answer every unknown path with the home response.
	mux.Handle("/{$}", handler.HomeHandler(cfg))
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.Handle("/version", handler.VersionHandler(cfg))
	mux.Handle("/", handler.NotFoundHandler())

	return &http.Server{
		Addr:              addr(),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
	}
}

// addr returns the listen address, honouring the ADDR override.
func addr() string {
	if v := os.Getenv("ADDR"); v != "" {
		return v
	}
	return defaultAddr
}

func main() {
	cfg := config.Load()
	srv := newServer(cfg)

	// Surface unexpected listener failures (for example a port already in
	// use) on a channel so the select below can tell them apart from a
	// clean shutdown.
	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s (app=%s env=%s version=%s)",
			srv.Addr, cfg.AppName, cfg.Environment, cfg.AppVersion)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
	case <-ctx.Done():
		stop()
		log.Println("shutdown signal received, draining connections")
	}

	// Use a fresh context: ctx is already cancelled by the signal, so
	// deriving from it would abort the drain immediately.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown incomplete: %v", err)
		// A drain that ran out of time still has to release the listener.
		_ = srv.Close()
		os.Exit(1)
	}

	log.Println("shutdown complete")
}
