package middleware

import (
	"net/http"
	"strings"

	"rygell-dashboard/internal/config"

	"github.com/gin-gonic/gin"
)

func RequestBodyLimit(limit int64, exemptPaths ...string) gin.HandlerFunc {
	exempt := make(map[string]bool, len(exemptPaths))
	for _, path := range exemptPaths {
		exempt[path] = true
	}
	return func(c *gin.Context) {
		if exempt[c.FullPath()] || exempt[c.Request.URL.Path] {
			c.Next()
			return
		}
		if limit > 0 && c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		}
		c.Next()
	}
}

func AdminOnly(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		if strings.TrimSpace(asString(username)) != adminUsername(cfg) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
			return
		}
		c.Next()
	}
}

func StrictOriginForUnsafeMethods(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := exactFrontendOrigins(cfg)
	return func(c *gin.Context) {
		if !isUnsafeMethod(c.Request.Method) {
			c.Next()
			return
		}

		origin := strings.TrimSpace(c.GetHeader("Origin"))
		if origin == "" {
			origin = strings.TrimSpace(c.GetHeader("Referer"))
			if origin != "" {
				origin = originFromReferer(origin)
			}
		}
		if origin == "" || !allowedOrigins[origin] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden origin"})
			return
		}
		c.Next()
	}
}

func exactFrontendOrigins(cfg *config.Config) map[string]bool {
	origins := map[string]bool{
		"http://localhost:3000": true,
		"http://localhost:3001": true,
	}
	if cfg == nil || strings.TrimSpace(cfg.FrontendOrigins) == "" {
		return origins
	}
	parsed := make(map[string]bool)
	for _, raw := range strings.Split(cfg.FrontendOrigins, ",") {
		origin := strings.TrimSpace(raw)
		if origin == "" || strings.Contains(origin, "*") {
			continue
		}
		parsed[origin] = true
	}
	if len(parsed) == 0 {
		return origins
	}
	return parsed
}

func originFromReferer(referer string) string {
	req, err := http.NewRequest(http.MethodGet, referer, nil)
	if err != nil || req.URL == nil || req.URL.Scheme == "" || req.URL.Host == "" {
		return ""
	}
	return req.URL.Scheme + "://" + req.URL.Host
}

func isUnsafeMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func adminUsername(cfg *config.Config) string {
	if cfg == nil || strings.TrimSpace(cfg.AdminUsername) == "" {
		return "admin"
	}
	return strings.TrimSpace(cfg.AdminUsername)
}

func asString(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}
