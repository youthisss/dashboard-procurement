package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application.
type Config struct {
	DatabaseURL       string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBLogLevel        string
	ServerPort        string
	UploadDir         string
	FrontendOrigins   string
	JWTSecret         string
	JWTExpiryHours    string
	CookieName        string
	CookieDomain      string
	CookieSameSite    string
	CookieSecure      bool
	AdminName         string
	AdminUsername     string
	AdminPassword     string
	AdminSyncPassword bool
	ImportKeepFiles   bool
	GinMode           string
	TrustedProxies    string
	MaxBodyBytes      int64
	ImportMaxBytes    int64
	ImportJSONBytes   int64
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "rygell_dashboard"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBLogLevel:        getEnv("DB_LOG_LEVEL", "warn"),
		ServerPort:        getEnv("SERVER_PORT", getEnv("PORT", "8080")),
		UploadDir:         getEnv("UPLOAD_DIR", "./uploads"),
		FrontendOrigins:   getEnv("FRONTEND_ORIGINS", "http://localhost:3000,http://localhost:3001"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		JWTExpiryHours:    getEnv("JWT_EXPIRY_HOURS", "24"),
		CookieName:        getEnv("COOKIE_NAME", "auth_token"),
		CookieDomain:      getEnv("COOKIE_DOMAIN", ""),
		CookieSameSite:    getEnv("COOKIE_SAMESITE", "Lax"),
		CookieSecure:      getEnvBool("COOKIE_SECURE", false),
		AdminName:         getEnv("ADMIN_NAME", "Administrator"),
		AdminUsername:     getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:     getEnv("ADMIN_PASSWORD", ""),
		AdminSyncPassword: getEnvBool("ADMIN_SYNC_PASSWORD", false),
		ImportKeepFiles:   getEnvBool("IMPORT_KEEP_FILES", false),
		GinMode:           getEnv("GIN_MODE", "debug"),
		TrustedProxies:    getEnv("TRUSTED_PROXIES", ""),
		MaxBodyBytes:      getEnvInt64("MAX_BODY_BYTES", 10*1024*1024),
		ImportMaxBytes:    getEnvInt64("IMPORT_MAX_BYTES", 55*1024*1024),
		ImportJSONBytes:   getEnvInt64("IMPORT_JSON_BYTES", 20*1024*1024),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
		return fallback
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(trimmed)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvInt64(key string, fallback int64) int64 {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
