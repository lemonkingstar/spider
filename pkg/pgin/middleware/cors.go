package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/lemonkingstar/spider/pkg/plog"
)

func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowedOriginMap := map[string]bool{}
	for _, o := range allowedOrigins {
		allowedOriginMap[o] = true
	}
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// If "api-cors-header" is not given, but "api-enable-cors" is true, we set cors to "*"
			// otherwise, all head values will be passed to HTTP handler
			corsHeaders := ""
			if len(allowedOriginMap) == 0 {
				corsHeaders = "*"
			} else if allowedOriginMap[origin] {
				corsHeaders = origin
			}
			// cors handle
			plog.Debugf("CORS header is enabled and set to: %s, Origin: %s", corsHeaders, origin)
			if corsHeaders != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", corsHeaders)
			}
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, DELETE, PUT, OPTIONS")

			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusOK)
				return
			}
		}
		c.Next()
	}
}

func GinCors(allowedOrigins []string) gin.HandlerFunc {
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
