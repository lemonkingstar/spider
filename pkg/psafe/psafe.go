// Package psafe
// @Title  psafe.go
// @Description  golang security operations
// @Author  wuzi
// @Update  wuzi
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

func Go(g func()) {
	go func() {
		defer Recover()
		if g != nil {
			g()
		}
	}()
}

func GoContext(g func(context.Context), ctx context.Context) {
	go func() {
		defer Recover()
		if g != nil {
			g(ctx)
		}
	}()
}

func Call(c func()) {
	defer Recover()
	if c != nil {
		c()
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

func GoLoop(g func(), interval ...int) {
	loop := func(g func()) {
		defer Recover()
		if g != nil {
			g()
		}
	}

	t := 2
	if len(interval) > 0 {
		t = interval[0]
	}
	go func() {
		for {
			loop(g)
			time.Sleep(time.Duration(t) * time.Second)
		}
	}()
}
