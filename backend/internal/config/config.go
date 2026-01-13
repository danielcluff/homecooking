package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	AI       AIConfig
	Storage  StorageConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port        int
	Environment string
	BaseURL     string
}

type DatabaseConfig struct {
	Type     string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Path     string
}

type AuthConfig struct {
	JWTSecret        string
	RefreshSecret    string
	TokenExpiryHours int
}

type AIConfig struct {
	Enabled   bool
	Provider  string
	APIKey    string
	Model     string
	BaseURL   string
	AutoApply bool
}

type StorageConfig struct {
	Type        string
	LocalPath   string
	MaxFileSize int64
}

type EmailConfig struct {
	Enabled  bool
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	From     string
}

func Load() (*Config, error) {
	cfg := &Config{}

	cfg.Server = ServerConfig{
		Port:        getEnvInt("SERVER_PORT", 8080),
		Environment: getEnv("SERVER_ENV", "development"),
		BaseURL:     getEnv("SERVER_BASE_URL", "http://localhost:8080"),
	}

	cfg.Database = DatabaseConfig{
		Type:     getEnv("DATABASE_TYPE", "postgres"),
		Host:     getEnv("DATABASE_HOST", "localhost"),
		Port:     getEnvInt("DATABASE_PORT", 5432),
		Name:     getEnv("DATABASE_NAME", "homecooking"),
		User:     getEnv("DATABASE_USER", "postgres"),
		Password: getEnv("DATABASE_PASSWORD", "password"),
		Path:     getEnv("DATABASE_PATH", "./data.db"),
	}

	cfg.Auth = AuthConfig{
		JWTSecret:        getEnv("JWT_SECRET", "change-me-in-production"),
		RefreshSecret:    getEnv("REFRESH_SECRET", "change-me-in-production"),
		TokenExpiryHours: getEnvInt("TOKEN_EXPIRY_HOURS", 24),
	}

	cfg.AI = AIConfig{
		Enabled:   getEnvBool("AI_ENABLED", false),
		Provider:  getEnv("AI_PROVIDER", "openai"),
		APIKey:    getEnv("AI_API_KEY", ""),
		Model:     getEnv("AI_MODEL", "gpt-4o"),
		BaseURL:   getEnv("AI_BASE_URL", ""),
		AutoApply: getEnvBool("AI_AUTO_APPLY", false),
	}

	cfg.Storage = StorageConfig{
		Type:        getEnv("STORAGE_TYPE", "local"),
		LocalPath:   getEnv("STORAGE_LOCAL_PATH", "./uploads"),
		MaxFileSize: getEnvInt64("STORAGE_MAX_UPLOAD_SIZE", 10485760),
	}

	cfg.Email = EmailConfig{
		Enabled:  getEnvBool("EMAIL_ENABLED", false),
		SMTPHost: getEnv("EMAIL_SMTP_HOST", ""),
		SMTPPort: getEnvInt("EMAIL_SMTP_PORT", 587),
		SMTPUser: getEnv("EMAIL_SMTP_USER", ""),
		SMTPPass: getEnv("EMAIL_SMTP_PASS", ""),
		From:     getEnv("EMAIL_FROM", ""),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Auth.JWTSecret == "change-me-in-production" && c.Server.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}

	if c.Auth.RefreshSecret == "change-me-in-production" && c.Server.Environment == "production" {
		return fmt.Errorf("REFRESH_SECRET must be set in production")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
