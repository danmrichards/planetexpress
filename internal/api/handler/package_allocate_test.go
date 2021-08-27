package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("oh dear oh dear")

func Test_handler_packageAllocate(t *testing.T) {
	tests := []struct {
		name        string
		req         json.RawMessage
		allocateErr error
		expCode     int
		expResponse json.RawMessage
	}{
		{
			name:        "golden path",
			req:         json.RawMessage(`{"size": 10}`),
			expCode:     http.StatusCreated,
			expResponse: json.RawMessage(`{"data": {"size": 10, "package_id": "test"}}`),
		},
		{
			name:        "no body",
			expCode:     http.StatusBadRequest,
			expResponse: json.RawMessage(`{"errors": [{"description": "Bad Request", "status": "400", "title": "Bad Request"}]}`),
		},
		{
			name:        "event error",
			req:         json.RawMessage(`{"size": 10}`),
			expCode:     http.StatusInternalServerError,
			allocateErr: errTest,
			expResponse: json.RawMessage(`{"errors": [{"description": "oh dear oh dear", "status": "500", "title": "could not allocate package"}]}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler{
				pkgIDGen: func() string {
					return "test"
				},
				evtSvc: &EventServiceMock{
					PackageAllocateFunc: func(ctx context.Context, id string, size int) error {
						return tt.allocateErr
					},
				},
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"http://example.com/v1/package/allocate",
				bytes.NewReader(tt.req),
			)
			w := httptest.NewRecorder()
			h.packageAllocate(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			assert.NoError(t, err)

			assert.JSONEq(t, string(tt.expResponse), string(body))
		})
	}
}
