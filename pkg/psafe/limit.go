package psafe

import (
	"sync"
)

// Golimit limits the maximum concurrent number of coroutines.
type Golimit struct {
	num    int
	ch     chan struct{}
	wg     sync.WaitGroup
	closed bool
}

// NewGl returns a *Golimit.
// usage: psafe.NewGl(10).Run(func() {})
func NewGl(num int) *Golimit {
	if num <= 0 {
		panic("invalid concurrent number")
	}
	return &Golimit{
		num: num,
		ch:  make(chan struct{}, num),
	}
}

func (g *Golimit) Run(f func()) {
	if g.closed {
		return
	}
	g.wg.Add(1)
	g.ch <- struct{}{}
	go func() {
		defer g.wg.Done()
		defer func() { <-g.ch }()
		Call(f)
	}()
}

func (g *Golimit) Wait() {
	g.wg.Wait()
	g.Close()
}

func (g *Golimit) Close() {
	if !g.closed {
		g.closed = true
		close(g.ch)
	}
}
