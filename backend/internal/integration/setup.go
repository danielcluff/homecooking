package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/config"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/handlers"
	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
	"github.com/homecooking/backend/internal/services"
	testutil "github.com/homecooking/backend/internal/testing"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type TestServer struct {
	DB          *sql.DB
	Queries     *sqlc.Queries
	Config      *config.Config
	UserRepo    *repository.UserRepository
	AuthHandler *handlers.AuthHandler
	AuthService *services.AuthService
}

// SetupTestServer creates a test server with minimal setup for integration tests
func SetupTestServer(t *testing.T) *TestServer {
	// Create in-memory database
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)

	// Initialize config
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key-for-integration-tests",
			RefreshSecret:    "test-refresh-secret-for-integration-tests",
			TokenExpiryHours: 1,
		},
	}

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(db, q)
	authService := services.NewAuthService(cfg, userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	return &TestServer{
		DB:          db,
		Queries:     q,
		Config:      cfg,
		UserRepo:    userRepo,
		AuthHandler: authHandler,
		AuthService: authService,
	}
}

// TeardownTestServer closes the test server resources
func TeardownTestServer(server *TestServer) {
	if server.DB != nil {
		server.DB.Close()
	}
}

// createTestUserDirect creates a test user using raw SQL (bypasses repository issues)
func createTestUserDirect(t *testing.T, server *TestServer, email, password string) string {
	userID := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	_, err = server.DB.Exec(`
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`, userID, email, string(hashedPassword), "user")

	require.NoError(t, err)
	return userID
}

// MakeRequestWithUser creates an HTTP request with user in context
func MakeRequestWithUser(method, url string, body interface{}, user *models.User) *http.Request {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req := httptest.NewRequest(method, url, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	ctx := context.WithValue(req.Context(), "user", user)
	req = req.WithContext(ctx)

	return req
}

// GetTestUser retrieves or creates a test user for integration tests
func GetTestUser(t *testing.T, server *TestServer, email, password string) *models.User {
	user, err := server.UserRepo.GetByEmail(email)
	if err == nil {
		return user
	}

	userID := createTestUserDirect(t, server, email, password)
	user, err = server.UserRepo.GetByID(userID)
	require.NoError(t, err)
	return user
}

// CreateTestUser creates a test user and returns an auth token
func CreateTestUser(t *testing.T, server *TestServer, email, password string) string {
	// Register user
	registerBody := map[string]string{
		"email":    email,
		"password": password,
	}
	registerBodyBytes, _ := json.Marshal(registerBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(registerBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.AuthHandler.Register(w, req)

	require.Equal(t, http.StatusCreated, w.Code, "Failed to register user: %s", w.Body.String())

	// Login to get token
	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}
	loginBodyBytes, _ := json.Marshal(loginBody)

	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	server.AuthHandler.Login(w, req)

	require.Equal(t, http.StatusOK, w.Code, "Failed to login user: %s", w.Body.String())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	token, ok := response["access_token"].(string)
	require.True(t, ok, "No access_token in response")
	require.NotEmpty(t, token, "Token is empty")

	return token
}

// MakeAuthenticatedRequest creates an HTTP request with auth token
func MakeAuthenticatedRequest(method, url string, body interface{}, token string) *http.Request {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req := httptest.NewRequest(method, url, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req
}
