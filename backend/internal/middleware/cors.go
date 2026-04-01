package middleware

import (
	"strings"
	"time"

	"rygell-dashboard/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a configured CORS middleware.
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	origins := []string{"http://localhost:3000", "http://localhost:3001"}
	if cfg != nil && strings.TrimSpace(cfg.FrontendOrigins) != "" {
		rawOrigins := strings.Split(cfg.FrontendOrigins, ",")
		parsed := make([]string, 0, len(rawOrigins))
		for _, origin := range rawOrigins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				parsed = append(parsed, trimmed)
			}
		}
		if len(parsed) > 0 {
			origins = parsed
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowWildcard:    true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
