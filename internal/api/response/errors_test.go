package response

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestWriteDefaultStatusError(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestWriteError(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		status int
		title  string
		detail string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_errorResponse(t *testing.T) {
	type args struct {
		status int
		title  string
		desc   string
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := errorResponse(tt.args.status, tt.args.title, tt.args.desc)
			if (err != nil) != tt.wantErr {
				t.Errorf("errorResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("errorResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
