package pgin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/lemonkingstar/spider/pkg/server"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// D2 通用返回数据结构
type D2 struct {
	Ret       int64       `json:"ret"`
	Msg       string      `json:"msg"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id"`
}

func DzCheckError(c *gin.Context, logger *zap.Logger, code int64, msg string, data ...interface{}) {
	requestID := c.Request.Header.Get(iserver.PXHTTPCCRequestID)
	var result = D2{Ret: code, Msg: msg, Timestamp: time.Now().Unix(), RequestID: requestID}
	if data != nil {
		result.Data = data[0]
	}
	c.JSON(http.StatusOK, result)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	r, _ := json.Marshal(result)
	if code != 0 {
		logger.Debug("request", zap.String("result", string(r)), zap.String("request-id", requestID))
	}
}

func GenerateRID() string {
	unused := "0000"
	guid := xid.New()
	return fmt.Sprintf("cc%s%s", unused, guid.String())
}

// RequestIDMiddleware
// g.Use(RequestIDMiddleware())
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Request.Header.Get(iserver.PXHTTPCCRequestID)
		if rid == "" {
			rid = GenerateRID()
			c.Request.Header.Set(iserver.PXHTTPCCRequestID, rid)
		}
		c.Writer.Header().Set(iserver.PXHTTPCCRequestID, rid)
		c.Set(iserver.ContextRequestID, rid)
	}
}
