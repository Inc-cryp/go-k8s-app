package handler

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

// decodeBody asserts the response advertises JSON and returns the decoded
// object, so each test does not repeat the header and unmarshal boilerplate.
func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]string {
	t.Helper()

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v (body=%q)", err, rec.Body.String())
	}
	return body
}

func TestHomeHandlerReportsConfigValues(t *testing.T) {
	rec := httptest.NewRecorder()
	HomeHandler(testConfig())(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := decodeBody(t, rec)
	if body["app"] != "go-k8s-app" {
		t.Errorf("app = %q, want %q", body["app"], "go-k8s-app")
	}
	if body["env"] != "production" {
		t.Errorf("env = %q, want %q", body["env"], "production")
	}
	if body["message"] == "" {
		t.Error("message is empty, want the welcome message")
	}
}

func TestHealthHandlerAlwaysReportsOK(t *testing.T) {
	rec := httptest.NewRecorder()
	HealthHandler(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body := decodeBody(t, rec); body["status"] != "ok" {
		t.Errorf("status field = %q, want %q", body["status"], "ok")
	}
}

func TestVersionHandlerReportsVersion(t *testing.T) {
	rec := httptest.NewRecorder()
	VersionHandler(testConfig())(rec, httptest.NewRequest(http.MethodGet, "/version", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := decodeBody(t, rec)
	if body["app_name"] != "go-k8s-app" {
		t.Errorf("app_name = %q, want %q", body["app_name"], "go-k8s-app")
	}
	if body["version"] != "v1.0.0" {
		t.Errorf("version = %q, want %q", body["version"], "v1.0.0")
	}
	if body["environment"] != "production" {
		t.Errorf("environment = %q, want %q", body["environment"], "production")
	}
}

func TestNotFoundHandlerReports404ForUnknownPath(t *testing.T) {
	rec := httptest.NewRecorder()
	NotFoundHandler()(rec, httptest.NewRequest(http.MethodGet, "/does-not-exist", nil))

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	body := decodeBody(t, rec)
	if body["error"] != "not found" {
		t.Errorf("error = %q, want %q", body["error"], "not found")
	}
	// The path echoes back what was requested, which makes debugging a
	// misrouted request from the logs straightforward.
	if body["path"] != "/does-not-exist" {
		t.Errorf("path = %q, want %q", body["path"], "/does-not-exist")
	}
}
