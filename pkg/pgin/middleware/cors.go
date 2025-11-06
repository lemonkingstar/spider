package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/lemonkingstar/spider/pkg/plog"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// If "api-cors-header" is not given, but "api-enable-cors" is true, we set cors to "*"
			// otherwise, all head values will be passed to HTTP handler
			corsHeaders := "*"

			// cors handle
			plog.Debugf("CORS header is enabled and set to: %s, Origin: %s", corsHeaders, origin)
			c.Writer.Header().Set("Access-Control-Allow-Origin", corsHeaders)
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept,Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, DELETE, PUT, OPTIONS")

			if c.Request.Method == http.MethodOptions {
				// c.Writer.WriteHeader(http.StatusOK)
				c.AbortWithStatus(http.StatusOK)
				return
			}
		}
		c.Next()
	}
}

func Cors(domain string) gin.HandlerFunc {
	var domains = []string{"*"}
	if domain != "" {
		domains = strings.Split(domain, ",")
	}

	return cors.New(cors.Config{
		AllowOrigins:     domains,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
