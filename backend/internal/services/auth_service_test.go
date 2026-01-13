package services

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/config"
	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
	testutil "github.com/homecooking/backend/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func createTestUserForAuth(t *testing.T, db *sql.DB, email string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	userID := uuid.New().String()

	_, err := db.Exec(`
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`, userID, email, string(hashedPassword), "user")

	require.NoError(t, err)

	return userID
}

func TestAuthService_Login(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db, q)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key",
			RefreshSecret:    "test-refresh-secret",
			TokenExpiryHours: 1,
		},
	}
	service := NewAuthService(cfg, userRepo)

	email := "loginuser@example.com"
	createTestUserForAuth(t, db, email)

	loginReq := &models.LoginRequest{
		Email:    email,
		Password: "password",
	}

	resp, err := service.Login(loginReq)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Greater(t, resp.ExpiresIn, 0)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db, q)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key",
			RefreshSecret:    "test-refresh-secret",
			TokenExpiryHours: 1,
		},
	}
	service := NewAuthService(cfg, userRepo)

	email := "wrongpassword@example.com"
	createTestUserForAuth(t, db, email)

	loginReq := &models.LoginRequest{
		Email:    email,
		Password: "wrongpassword",
	}

	_, err = service.Login(loginReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db, q)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key",
			RefreshSecret:    "test-refresh-secret",
			TokenExpiryHours: 1,
		},
	}
	service := NewAuthService(cfg, userRepo)

	loginReq := &models.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err = service.Login(loginReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_ValidateToken(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db, q)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key",
			RefreshSecret:    "test-refresh-secret",
			TokenExpiryHours: 1,
		},
	}
	service := NewAuthService(cfg, userRepo)

	email := "validate@example.com"
	userID := createTestUserForAuth(t, db, email)

	loginReq := &models.LoginRequest{
		Email:    email,
		Password: "password",
	}
	loginResp, err := service.Login(loginReq)
	require.NoError(t, err)

	user, err := service.ValidateToken(loginResp.AccessToken)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID.String())
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db, q)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key",
			RefreshSecret:    "test-refresh-secret",
			TokenExpiryHours: 1,
		},
	}
	service := NewAuthService(cfg, userRepo)

	_, err = service.ValidateToken("invalid.token.string")
	assert.Error(t, err)
}

func TestAuthService_RefreshToken(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db, q)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key",
			RefreshSecret:    "test-refresh-secret",
			TokenExpiryHours: 1,
		},
	}
	service := NewAuthService(cfg, userRepo)

	email := "refresh@example.com"
	createTestUserForAuth(t, db, email)

	loginReq := &models.LoginRequest{
		Email:    email,
		Password: "password",
	}
	loginResp, err := service.Login(loginReq)
	require.NoError(t, err)

	refreshResp, err := service.RefreshToken(loginResp.RefreshToken)
	require.NoError(t, err)
	assert.NotNil(t, refreshResp)
	assert.NotEmpty(t, refreshResp.AccessToken)
	assert.NotEmpty(t, refreshResp.RefreshToken)
	assert.Greater(t, refreshResp.ExpiresIn, 0)
}

func TestAuthService_RefreshToken_Invalid(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db, q)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key",
			RefreshSecret:    "test-refresh-secret",
			TokenExpiryHours: 1,
		},
	}
	service := NewAuthService(cfg, userRepo)

	_, err = service.RefreshToken("invalid.token.string")
	assert.Error(t, err)
}
