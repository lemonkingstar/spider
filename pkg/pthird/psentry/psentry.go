package psentry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

const (
	valuesKey = "sentry"
	sentryContext = "sentryContext"
)

// CaptureException Concurrency is not safe
func CaptureException(err error, extras map[string]interface{}) {
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetExtras(extras)
			sentry.CaptureException(err)
		})
	}
}

// CaptureMessage Concurrency is not safe
func CaptureMessage(message string, extras map[string]interface{}) {
	if message != "" {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetExtras(extras)
			sentry.CaptureMessage(message)
		})
	}
}

// BindMiddleware Concurrency is safe
// Usage:
// f := BindMiddleware(r, true)
// defer f(false)
func BindMiddleware(r *http.Request, track bool) (*http.Request, func(bool)) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func (scope *sentry.Scope) {
		scope.SetRequest(r)
		scope.SetTags(map[string]string{})
		scope.SetContexts(map[string]sentry.Context{})
		scope.SetExtras(map[string]interface{}{})
		scope.SetUser(sentry.User{})
	})
	// hub.Recover(err)
	// hub.CaptureException(err)
	ctx := r.Context()
	ctx = SetHub(ctx, hub)
	if track {
		span := sentry.StartSpan(ctx, "http.sentry",
			sentry.TransactionName(fmt.Sprintf("%s %s", r.Method, r.URL.Path)),
		)
		r = r.WithContext(span.Context())
		return r, func(b bool) {
			if err := recover(); err != nil {
				HandlePanic(ctx, r, b, err)
			}
			span.Finish()
		}
	}
	r = r.WithContext(ctx)
	return r, func(b bool) {
		if err := recover(); err != nil {
			HandlePanic(ctx, r, b, err)
		}
	}
}

func SetHub(ctx context.Context, hub *sentry.Hub) context.Context {
	return sentry.SetHubOnContext(ctx, hub)
}

func GetHub(ctx context.Context) *sentry.Hub {
	return sentry.GetHubFromContext(ctx)
}

func HandlePanic(ctx context.Context, r *http.Request, rePanic bool, err interface{}) {
	hub := GetHub(ctx)
	if hub != nil {
		_ = hub.RecoverWithContext(
			context.WithValue(ctx, sentry.RequestContextKey, r),
			err,
		)
	}
	if rePanic {
		panic(err)
	}
}

// GetHubWithGin 可以根据需要保存 local hub到指定的对象中
// GetHubFromContext retrieves attached *sentry.Hub instance from gin.Context.
func GetHubWithGin(ctx *gin.Context) *sentry.Hub {
	if hub, ok := ctx.Get(valuesKey); ok {
		if hub, ok := hub.(*sentry.Hub); ok {
			return hub
		}
	}
	return nil
}

func SetHubWithGin(ctx *gin.Context, hub *sentry.Hub) {
	ctx.Set(valuesKey, hub)
}

func GetContextWithGin(ctx *gin.Context) context.Context {
	if hub, ok := ctx.Get(sentryContext); ok {
		if ctx, ok := hub.(context.Context); ok {
			return ctx
		}
	}
	return nil
}

func SetContextWithGin(c *gin.Context, ctx context.Context) {
	c.Set(sentryContext, ctx)
	context.WithValue(c, sentry.RequestContextKey, ctx)
}

func GetRootSpanContext(ctx context.Context) context.Context {
	if span := sentry.TransactionFromContext(ctx); span != nil {
		return span.Context()
	}
	return nil
}
