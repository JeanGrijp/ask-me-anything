// JSON is a helper function to send JSON responses
// It sets the Content-Type header to application/json and encodes the response data as JSON.
package responses

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
