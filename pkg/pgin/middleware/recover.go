package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/lemonkingstar/spider/pkg/pgin/response"
	"github.com/lemonkingstar/spider/pkg/plog"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(response.E); ok {
					plog.Errorf("[ App Error ] code: %d, msg: %s", e.Code(), e.Error())
					c.JSON(http.StatusOK, response.D{
						Status: 1, Code: e.Code(), Message: e.Error(),
					})
					return
				}
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				httpReq, _ := httputil.DumpRequest(c.Request, true)
				plog.Errorf("[ Recovery ] panic recovered:\n%s\n%s\n%s", string(httpReq), err, buf[:n])
				c.JSON(http.StatusInternalServerError, response.D{
					Status: 1, Code: 500, Message: fmt.Sprintf("%v", err),
				})
			}
		}()
		c.Next()
	}
}
