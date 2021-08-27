package response

import (
	"net/http"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteJSON(tt.args.w, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("WriteJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
