package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"rygell-dashboard/internal/config"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware_AllowedOriginUsesExplicitOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(CORSMiddleware(&config.Config{FrontendOrigins: "https://dashboard-procurement-liard.vercel.app"}))
	r.OPTIONS("/api/v1/auth/login", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	req.Header.Set("Origin", "https://dashboard-procurement-liard.vercel.app")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://dashboard-procurement-liard.vercel.app" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want explicit allowed origin", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q, want true", got)
	}
}

func TestCORSMiddleware_WildcardConfigNeverReturnsWildcardOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(CORSMiddleware(&config.Config{FrontendOrigins: "*"}))
	r.OPTIONS("/api/v1/auth/login", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	req.Header.Set("Origin", "https://dashboard-procurement-liard.vercel.app")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got == "*" {
		t.Fatalf("Access-Control-Allow-Origin must never be wildcard for credentialed auth")
	}
}
