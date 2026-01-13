package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompleteAuthFlow(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	email := "authflow@example.com"
	password := "password123"

	userID := createTestUserDirect(t, server, email, password)

	assert.NotNil(t, userID)
	assert.NotEmpty(t, userID)

	user, err := server.UserRepo.GetByEmail(email)
	require.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "user", user.Role)

	loginReq := map[string]string{
		"email":    email,
		"password": password,
	}
	loginBodyBytes, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.AuthHandler.Login(w, req)

	require.Equal(t, http.StatusOK, w.Code, "Login failed: %s", w.Body.String())

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	token, ok := response["access_token"].(string)
	require.True(t, ok, "No access_token in response")
	require.NotEmpty(t, token, "Token is empty")

	user, err = server.AuthService.ValidateToken(token)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
}

func TestLoginWithWrongPassword(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	email := "wrongpass@example.com"
	password := "correctpassword"

	createTestUserDirect(t, server, email, password)

	loginReq := map[string]string{
		"email":    email,
		"password": "wrongpassword",
	}
	loginBodyBytes, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.AuthHandler.Login(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "invalid")
}

func TestLoginNonExistentUser(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	loginReq := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "password123",
	}
	loginBodyBytes, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.AuthHandler.Login(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "invalid")
}

func TestValidateValidToken(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	email := "validtoken@example.com"
	password := "password123"

	createTestUserDirect(t, server, email, password)

	loginReq := map[string]string{
		"email":    email,
		"password": password,
	}
	loginBodyBytes, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.AuthHandler.Login(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	token, ok := response["access_token"].(string)
	require.True(t, ok)

	user, err := server.AuthService.ValidateToken(token)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
}

func TestValidateInvalidToken(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	_, err := server.AuthService.ValidateToken("invalid.token.string")
	assert.Error(t, err)
}

func TestRefreshTokenFlow(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	email := "refreshflow@example.com"
	password := "password123"

	createTestUserDirect(t, server, email, password)

	loginReq := map[string]string{
		"email":    email,
		"password": password,
	}
	loginBodyBytes, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.AuthHandler.Login(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	refreshToken, ok := response["refresh_token"].(string)
	require.True(t, ok, "No refresh_token in response")

	refreshReq := map[string]string{
		"refresh_token": refreshToken,
	}
	refreshBodyBytes, _ := json.Marshal(refreshReq)

	req = httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader(refreshBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	server.AuthHandler.Refresh(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Refresh failed: %s", w.Body.String())

	var refreshResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &refreshResponse)
	require.NoError(t, err)

	newToken, ok := refreshResponse["access_token"].(string)
	require.True(t, ok)
	assert.NotEmpty(t, newToken)

	user, err := server.AuthService.ValidateToken(newToken)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
}
