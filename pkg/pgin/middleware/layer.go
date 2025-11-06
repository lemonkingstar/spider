package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lemonkingstar/spider/pkg/pbase"
	"github.com/lemonkingstar/spider/pkg/pgin/pginutil"
	"github.com/lemonkingstar/spider/pkg/plog"
)

// LayerMiddleware ~
// detail 是否展示详细数据
// user   是否显示用户信息
func LayerMiddleware(detail bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		method := c.Request.Method
		body := ""
		if method != http.MethodGet && detail {
			b, _ := pginutil.PeekRequest(c.Request)
			body = string(b)
		}
		c.Next()
		plog.Debugf("[ %s ]->[ %s|%d ] : %dms, {\"Client-Ip\": \"%s\", \"User-Agent\": \"%s\", \"User\": \"%s\", \"Body\": \"%s\"}",
			c.Request.RequestURI, method, c.Writer.Status(),
			time.Since(t)/time.Millisecond,
			pginutil.GetRemoteIP(c),
			c.Request.UserAgent(),
			pbase.GetUser(c.Request.Header),
			body,
		)
	}
}
