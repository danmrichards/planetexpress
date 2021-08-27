package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteJSON writes the given data to w as JSON with a default 200 OK status.
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	return WriteStatusJSON(w, http.StatusOK, data)
}

// WriteStatusJSON writes the given data to w as JSON with the given HTTP
// status code.
func WriteStatusJSON(w http.ResponseWriter, status int, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %w", err)
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(b)
	return err
}
