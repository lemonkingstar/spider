// Package psafe
// @Title psafe.go
// @Description golang security operations
// @Author spider
// @Update spider
package psafe

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/lemonkingstar/spider/pkg/plog"
)

var (
	logger = plog.WithField("[PACKET]", "psafe")
)

func Go(f func()) {
	go func() {
		defer Recover()
		if f != nil {
			f()
		}
	}()
}

func GoCtx(f func(context.Context), ctx context.Context) {
	go func() {
		defer Recover()
		if f != nil {
			f(ctx)
		}
	}()
}

func Call(f func()) {
	defer Recover()
	if f != nil {
		f()
	}
}

func Recover(handlers ...func(interface{})) {
	if err := recover(); err != nil {
		logger.Errorf("Panic: %v", err)
		logger.Errorf("Stack:\n %s", debug.Stack())
		for _, handler := range handlers {
			handler(err)
		}
	}
}

func GoForever(f func(), interval ...int) {
	loop := func(f func()) {
		defer Recover()
		if f != nil {
			f()
		}
	}

	t := 2
	if len(interval) > 0 {
		t = interval[0]
	}
	go func() {
		for {
			loop(f)
			time.Sleep(time.Duration(t) * time.Second)
		}
	}()
}
