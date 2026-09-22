package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON writes payload as the response body with the given status code.
//
// The status line is sent before the body, so an encoding failure can no
// longer be reported to the client; it is logged instead of being dropped.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("response: failed to encode body: %v", err)
	}
}
