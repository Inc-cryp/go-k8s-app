package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Inc-cryp/go-k8s-app/internal/config"
)

func testConfig() config.Config {
	return config.Config{
		AppName:     "go-k8s-app",
		AppVersion:  "v1.0.0",
		Environment: "production",
		APIKey:      "test-key",
	}
}

// serve routes one request through the real server handler.
func serve(t *testing.T, method, path string) *httptest.ResponseRecorder {
	t.Helper()

	rec := httptest.NewRecorder()
	newServer(testConfig()).Handler.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) map[string]string {
	t.Helper()

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not valid JSON: %v (body=%q)", err, rec.Body.String())
	}
	return body
}

func TestRoutesServeTheirDocumentedEndpoints(t *testing.T) {
	tests := []struct {
		path       string
		wantStatus int
		wantKey    string
	}{
		{"/", http.StatusOK, "message"},
		{"/health", http.StatusOK, "status"},
		{"/version", http.StatusOK, "version"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			rec := serve(t, http.MethodGet, tt.path)

			if rec.Code != tt.wantStatus {
				t.Errorf("GET %s status = %d, want %d", tt.path, rec.Code, tt.wantStatus)
			}
			if _, ok := decode(t, rec)[tt.wantKey]; !ok {
				t.Errorf("GET %s body has no %q key: %s", tt.path, tt.wantKey, rec.Body.String())
			}
		})
	}
}

// TestUnknownPathsDoNotFallThroughToTheHomeHandler guards the bug where the
// root handler was registered as a bare "/" pattern. A bare "/" pattern in
// net/http matches every path that no other pattern claims, so /does-not-exist
// was answered with the home response and status 200 instead of a 404.
func TestUnknownPathsDoNotFallThroughToTheHomeHandler(t *testing.T) {
	for _, path := range []string{"/does-not-exist", "/deep/nested/path", "/healthz"} {
		t.Run(path, func(t *testing.T) {
			rec := serve(t, http.MethodGet, path)

			if rec.Code != http.StatusNotFound {
				t.Errorf("GET %s status = %d, want %d", path, rec.Code, http.StatusNotFound)
			}

			body := decode(t, rec)
			if _, ok := body["message"]; ok {
				t.Errorf("GET %s fell through to HomeHandler: %s", path, rec.Body.String())
			}
			if body["error"] != "not found" {
				t.Errorf("GET %s error = %q, want %q", path, body["error"], "not found")
			}
		})
	}
}

func TestAddrPrefersEnvOverride(t *testing.T) {
	t.Setenv("ADDR", "127.0.0.1:9999")
	if got := addr(); got != "127.0.0.1:9999" {
		t.Errorf("addr() = %q, want %q", got, "127.0.0.1:9999")
	}

	t.Setenv("ADDR", "")
	if got := addr(); got != defaultAddr {
		t.Errorf("addr() = %q, want default %q", got, defaultAddr)
	}
}

func TestServerHasTimeoutHardening(t *testing.T) {
	srv := newServer(testConfig())

	// Without ReadHeaderTimeout a client can hold a connection open
	// indefinitely by sending headers slowly.
	if srv.ReadHeaderTimeout <= 0 {
		t.Error("ReadHeaderTimeout is not set")
	}
	if srv.IdleTimeout <= 0 {
		t.Error("IdleTimeout is not set")
	}
	if srv.Handler == nil {
		t.Fatal("Handler is nil, the server would serve nothing")
	}
}
