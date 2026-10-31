package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVercelHandlerDemoLogin(t *testing.T) {
	t.Setenv("POSTGRES_DSN", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("JWT_SECRET", "test-secret")

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"user@example.com","password":"password123","role":"user"}`))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	Handler(response, req)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d; got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var payload struct {
		Token string `json:"token"`
		Role  string `json:"role"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if payload.Token == "" || payload.Role != "user" || payload.Name != "Demo User" {
		t.Fatalf("unexpected login response: %+v", payload)
	}

}
