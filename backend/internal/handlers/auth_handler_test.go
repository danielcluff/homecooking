package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	handler := NewAuthHandler(nil)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(`{"email":"test@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	handler := NewAuthHandler(nil)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"test@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Refresh_InvalidJSON(t *testing.T) {
	handler := NewAuthHandler(nil)

	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Refresh(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_ResponseHeaders(t *testing.T) {
	handler := NewAuthHandler(nil)

	tests := []struct {
		name       string
		method     string
		url        string
		body       string
		wantStatus int
	}{
		{"Register missing fields", "POST", "/api/v1/auth/register", `{"email":"test@example.com"}`, http.StatusBadRequest},
		{"Login missing password", "POST", "/api/v1/auth/login", `{"email":"test@example.com"}`, http.StatusBadRequest},
		{"Refresh missing token", "POST", "/api/v1/auth/refresh", `{}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			switch tt.url {
			case "/api/v1/auth/register":
				handler.Register(w, req)
			case "/api/v1/auth/login":
				handler.Login(w, req)
			case "/api/v1/auth/refresh":
				handler.Refresh(w, req)
			}

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")
		})
	}
}
