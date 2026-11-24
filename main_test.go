package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	mux := newMux("hello", "secret", "test-pod")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("body = %q, want %q", rec.Body.String(), "ok")
	}
}

func TestRootReturnsMessageAndVersion(t *testing.T) {
	mux := newMux("hello", "secret", "test-pod")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body["message"] != "hello" {
		t.Errorf("message = %q, want %q", body["message"], "hello")
	}
	if body["version"] != version {
		t.Errorf("version = %q, want %q", body["version"], version)
	}
}

func TestPrivateNeedsAPIKey(t *testing.T) {
	mux := newMux("hello", "secret", "test-pod")

	tests := []struct {
		name string
		key  string
		want int
	}{
		{"no key", "", http.StatusUnauthorized},
		{"wrong key", "nope", http.StatusUnauthorized},
		{"right key", "secret", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/private", nil)
			if tt.key != "" {
				req.Header.Set("X-API-KEY", tt.key)
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.want {
				t.Fatalf("status = %d, want %d", rec.Code, tt.want)
			}
		})
	}
}
