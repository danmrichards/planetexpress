package response

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()

	assert.NoError(t, WriteJSON(w, map[string]interface{}{
		"foo": "bar",
	}))

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.NoError(t, err)

	assert.JSONEq(t, `{"foo": "bar"}`, string(body))
}
