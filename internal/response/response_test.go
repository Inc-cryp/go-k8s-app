package response

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONWritesStatusContentTypeAndBody(t *testing.T) {
	rec := httptest.NewRecorder()
	JSON(rec, http.StatusCreated, map[string]string{"hello": "world"})

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not valid JSON: %v (body=%q)", err, rec.Body.String())
	}
	if body["hello"] != "world" {
		t.Errorf("body[hello] = %q, want %q", body["hello"], "world")
	}
}

// TestJSONKeepsTheStatusWhenThePayloadCannotBeEncoded pins the ordering
// guarantee: the header and status are written before encoding starts, so a
// payload json cannot represent still yields the requested status code rather
// than a 200 with a truncated body.
func TestJSONKeepsTheStatusWhenThePayloadCannotBeEncoded(t *testing.T) {
	rec := httptest.NewRecorder()
	// A NaN float is not representable in JSON, so Encode fails midway.
	JSON(rec, http.StatusInternalServerError, map[string]float64{"n": math.NaN()})

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
}

func TestJSONEncodesNilPayloadWithoutPanicking(t *testing.T) {
	rec := httptest.NewRecorder()
	JSON(rec, http.StatusOK, nil)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != "null\n" {
		t.Errorf("body = %q, want %q", got, "null\n")
	}
}
