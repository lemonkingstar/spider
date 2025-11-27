package psafe

import (
	"fmt"
	"sync"
)

// Group can concurrently execute multiple functions and obtain all the execution results.
type Group struct {
	wg   sync.WaitGroup
	mu   sync.Mutex
	errs []error
}

func NewGroup() *Group {
	return &Group{}
}

// Run calls the given function in a new goroutine.
func (g *Group) Run(f func() error) {
	if f == nil {
		return
	}
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				g.mu.Lock()
				defer g.mu.Unlock()
				g.errs = append(g.errs, fmt.Errorf("panic recovered: %v", err))
			}
		}()
		if err := f(); err != nil {
			g.mu.Lock()
			defer g.mu.Unlock()
			g.errs = append(g.errs, err)
		}
	}()
}

func (g *Group) WaitErrors() []error {
	g.wg.Wait()
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.errs) == 0 {
		return nil
	}

	errs := make([]error, len(g.errs))
	copy(errs, g.errs)
	return errs
}

// WaitError only returns the first error.
func (g *Group) WaitError() error {
	errs := g.WaitErrors()
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (g *Group) Clear() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.errs = nil
}
