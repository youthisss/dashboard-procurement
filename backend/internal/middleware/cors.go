package middleware

import (
	"time"

	"rygell-dashboard/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a configured CORS middleware.
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	origins := frontendOriginList(cfg)

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowWildcard:    false,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
