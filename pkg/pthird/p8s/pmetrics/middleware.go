package pmetrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// Middleware is the http middleware for go-restful framework
func (s *Service) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/metrics" || path == "/metrics/" {
			s.ServeHTTP(c.Writer, c.Request)
			return
		}

		before := time.Now()
		c.Next()

		s.requestDuration.With(prometheus.Labels{LabelHandler: path}).
			Observe(float64(time.Since(before) / time.Millisecond))
		s.requestTotal.With(
			prometheus.Labels{LabelHandler: path, LabelHTTPStatus: strconv.Itoa(c.Writer.Status())},
		).Inc()
	}
}
