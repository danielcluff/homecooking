package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

type MockAuthService struct {
	validToken  string
	user        *models.User
	validateErr error
}

func (m *MockAuthService) ValidateToken(token string) (*models.User, error) {
	if token == m.validToken {
		return m.user, m.validateErr
	}
	return nil, m.validateErr
}

func TestAuth_ValidToken(t *testing.T) {
	mockUser := &models.User{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Email: "test@example.com",
		Role:  "user",
	}
	mockAuth := &MockAuthService{
		validToken:  "valid-token",
		user:        mockUser,
		validateErr: nil,
	}

	authMiddleware := NewAuthMiddleware(mockAuth)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		user := r.Context().Value(UserKey).(*models.User)
		assert.Equal(t, mockUser.ID, user.ID)
		assert.Equal(t, mockUser.Email, user.Email)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	authMiddleware.Auth(next).ServeHTTP(w, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuth_MissingHeader(t *testing.T) {
	mockAuth := &MockAuthService{}
	middleware := NewAuthMiddleware(mockAuth)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.Auth(next).ServeHTTP(w, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestAuth_InvalidHeaderFormat(t *testing.T) {
	mockAuth := &MockAuthService{}
	middleware := NewAuthMiddleware(mockAuth)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "invalid-token")
	w := httptest.NewRecorder()

	middleware.Auth(next).ServeHTTP(w, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization header format")
}

func TestAuth_InvalidToken(t *testing.T) {
	mockAuth := &MockAuthService{
		validateErr: assert.AnError,
	}
	middleware := NewAuthMiddleware(mockAuth)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	middleware.Auth(next).ServeHTTP(w, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestRequireRole_Admin(t *testing.T) {
	mockUser := &models.User{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Email: "admin@example.com",
		Role:  "admin",
	}
	mockAuth := &MockAuthService{
		user: mockUser,
	}
	authMiddleware := NewAuthMiddleware(mockAuth)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	ctx := context.WithValue(req.Context(), UserKey, mockUser)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	authMiddleware.RequireRole("admin")(next).ServeHTTP(w, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_Unauthorized(t *testing.T) {
	mockUser := &models.User{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Email: "user@example.com",
		Role:  "user",
	}

	authMiddleware := NewAuthMiddleware(&MockAuthService{})

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	ctx := context.WithValue(req.Context(), UserKey, mockUser)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	authMiddleware.RequireRole("admin")(next).ServeHTTP(w, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Forbidden")
}
