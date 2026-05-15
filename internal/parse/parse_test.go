package parse_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stpotter16/hours/internal/parse"
)

func TestParseLoginPost(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{"valid", `{"passphrase":"secret"}`, false},
		{"missing passphrase", `{}`, true},
		{"empty passphrase", `{"passphrase":""}`, true},
		{"invalid json", `not json`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
			got, err := parse.ParseLoginPost(r)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Passphrase == "" {
				t.Error("expected non-empty passphrase")
			}
		})
	}
}

func TestParseProjectCreateRequest(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{"valid", `{"name":"my-project"}`, false},
		{"missing name", `{}`, true},
		{"empty name", `{"name":""}`, true},
		{"invalid json", `not json`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(tt.body))
			got, err := parse.ParseProjectCreateRequest(r)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name == "" {
				t.Error("expected non-empty name")
			}
		})
	}
}
