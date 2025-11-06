package pgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lemonkingstar/spider/pkg/pgin/response"
)

var (
	successStatus = 0
	failStatus    = 1

	// 通用成功码/错误码
	successAction = 200
	failAction    = 500
)

type appError struct {
	code    int
	message string
}

func (p *appError) Error() string {
	return p.message
}

func (p *appError) Code() int {
	return p.code
}

func (p *appError) IsAppError() bool {
	return true
}

func customResponse(c *gin.Context, httpCode, status, code int, msg string, data interface{}) {
	c.JSON(httpCode, response.D{
		Status:  status,
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// Err 错误信息返回
// Err(c, "error message", 10000)
func Err(c *gin.Context, msg string, v ...interface{}) {
	errCode := failAction
	if len(v) > 0 {
		if code, ok := v[0].(int); ok {
			errCode = code
		}
	}
	customResponse(c, http.StatusOK, failStatus, errCode, msg, nil)
}

func Bomb(c *gin.Context, msg string, v ...interface{}) {
	errCode := failAction
	if len(v) > 0 {
		if code, ok := v[0].(int); ok {
			errCode = code
		}
	}
	customResponse(c, http.StatusInternalServerError, failStatus, errCode, msg, nil)
}

// OK 成功数据返回
func OK(c *gin.Context, data interface{}) {
	customResponse(c, http.StatusOK, successStatus, successAction, "", data)
}

// Success 成功返回
func Success(c *gin.Context) {
	customResponse(c, http.StatusOK, successStatus, successAction, "success", nil)
}

// Assert 条件断言
// 当断言条件为 假 时触发panic,返回错误信息
// Assert(result, "expression is false")
// Assert(result, "expression is false", 10100111)
func Assert(condition bool, errMsg string, v ...interface{}) {
	if !condition {
		errCode := successAction
		if len(v) > 0 {
			if code, ok := v[0].(int); ok {
				errCode = code
			}
		}

		panic(appError{
			code:    errCode,
			message: errMsg,
		})
	}
}

// CheckErr 错误检查
// CheckErr(err)
// CheckErr(err, "invalid param")
// CheckErr(err, "invalid param", 10100111)
func CheckErr(err error, v ...interface{}) {
	// msg = v[0]
	// code = v[1]
	if err != nil {
		errMsg := err.Error()
		if len(v) > 0 {
			if msg, ok := v[0].(string); ok {
				errMsg = msg
			}
		}
		errCode := successAction
		if len(v) > 1 {
			if code, ok := v[1].(int); ok {
				errCode = code
			}
		}

		panic(&appError{
			code:    errCode,
			message: errMsg,
		})
	}
}
