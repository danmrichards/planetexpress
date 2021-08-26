package response

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/danmrichards/planetexpress/internal/api"
)

// WriteError handles writing a proper error response.
func WriteError(w http.ResponseWriter, status int, title, detail string) {
	msg, err := errorResponse(status, title, detail)
	if err != nil {
		log.Printf("could not marshal error: %s\n", err)

		// Use a generic HTTP error in the event we cannot marshal the error.
		http.Error(w, title, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(msg)
}

// errorResponse return a general error response for all ErrorResponses
func errorResponse(status int, title, desc string) (json.RawMessage, error) {
	msg := api.ErrorResponse{
		Errors: []api.Error{
			{
				Description: desc,
				Status:      strconv.Itoa(status),
				Title:       title,
			},
		},
	}

	return json.Marshal(msg)
}
