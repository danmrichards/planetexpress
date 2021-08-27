package response

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()

	WriteError(w, http.StatusInternalServerError, "title", "description")

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.NoError(t, err)

	assert.JSONEq(
		t,
		`{"errors": [{"description": "description", "status": "500", "title": "title"}]}`,
		string(body),
	)
}

func TestWriteDefaultStatusError(t *testing.T) {
	w := httptest.NewRecorder()

	WriteDefaultStatusError(w, http.StatusInternalServerError)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.NoError(t, err)

	assert.JSONEq(
		t,
		`{"errors": [{"description": "Internal Server Error", "status": "500", "title": "Internal Server Error"}]}`,
		string(body),
	)
}
