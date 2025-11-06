package prequest

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ProxyHttp gin
// addr: target
// router.Any("/proxy/*name", proxyHandler)
// ProxyHttp(c, "http://127.0.0.1:8081")
func ProxyHttp(c *gin.Context, addr string) {

	u, err := url.Parse(addr)
	if err == nil {
		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		c.Writer.Write([]byte(err.Error()))
	}
}

// ProxyHttpChangePath gin
// router.Any("/proxy/*name", proxyHandler)
// 修改request path: http://127.0.0.1:8080/proxy/test -> http://127.0.0.1:8081/test
func ProxyHttpChangePath(c *gin.Context, addr string) {

	u, err := url.Parse(addr)
	if err == nil {
		c.Request.URL.Path = c.Param("name")
		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		c.Writer.Write([]byte(err.Error()))
	}
}

// ProxyHttpChangeRequest gin
// router.Any("/proxy/*name", proxyHandler)
func ProxyHttpChangeRequest(c *gin.Context, addr string) {

	u, err := url.Parse(addr)
	if err == nil {
		c.Request.URL.Path = c.Param("name")
		proxy := httputil.NewSingleHostReverseProxy(u)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			// 执行原处理函数
			originalDirector(req)
			// 修改header
			req.Header.Set("x-token", "value")
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		c.Writer.Write([]byte(err.Error()))
	}
}

// ProxyHttpChangeResponse gin
// router.Any("/proxy/*name", proxyHandler)
func ProxyHttpChangeResponse(c *gin.Context, addr string) {

	u, err := url.Parse(addr)
	if err == nil {
		c.Request.URL.Path = c.Param("name")
		proxy := httputil.NewSingleHostReverseProxy(u)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			// 执行原处理函数
			originalDirector(req)
			// 修改header
			req.Header.Set("x-token", "value")
		}
		// ModifyResponse默认nil 直接提供即可
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Del("Access-Control-Allow-Origin")
			resp.Header.Set("x_token","value")
			return nil
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		c.Writer.Write([]byte(err.Error()))
	}
}
