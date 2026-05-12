package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"rygell-dashboard/internal/config"

	"github.com/gin-gonic/gin"
)

func TestStrictOriginForUnsafeMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(StrictOriginForUnsafeMethods(&config.Config{FrontendOrigins: "https://app.example.com"}))
	r.POST("/mutate", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	r.GET("/read", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	blocked := httptest.NewRecorder()
	blockedReq := httptest.NewRequest(http.MethodPost, "/mutate", nil)
	blockedReq.Header.Set("Origin", "https://evil.example.com")
	r.ServeHTTP(blocked, blockedReq)
	if blocked.Code != http.StatusForbidden {
		t.Fatalf("blocked status = %d, want %d", blocked.Code, http.StatusForbidden)
	}

	allowed := httptest.NewRecorder()
	allowedReq := httptest.NewRequest(http.MethodPost, "/mutate", nil)
	allowedReq.Header.Set("Origin", "https://app.example.com")
	r.ServeHTTP(allowed, allowedReq)
	if allowed.Code != http.StatusNoContent {
		t.Fatalf("allowed status = %d, want %d", allowed.Code, http.StatusNoContent)
	}

	read := httptest.NewRecorder()
	r.ServeHTTP(read, httptest.NewRequest(http.MethodGet, "/read", nil))
	if read.Code != http.StatusNoContent {
		t.Fatalf("read status = %d, want %d", read.Code, http.StatusNoContent)
	}
}

func TestAdminOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.POST("/mutate", func(c *gin.Context) {
		c.Set("username", "staff")
	}, AdminOnly(&config.Config{AdminUsername: "owner"}), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/mutate", nil))
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func TestRequestBodyLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestBodyLimit(4))
	r.POST("/body", func(c *gin.Context) {
		if _, err := c.GetRawData(); err != nil {
			c.Status(http.StatusRequestEntityTooLarge)
			return
		}
		c.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/body", strings.NewReader("12345")))
	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusRequestEntityTooLarge)
	}
}
