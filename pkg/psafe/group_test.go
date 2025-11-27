package psafe

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	g := NewGroup()

	f := func(num int) error {
		time.Sleep(time.Duration(num) * time.Millisecond)
		if num%2 == 0 {
			return nil
		}
		if num > 3 {
			return nil
		}
		return fmt.Errorf("error: %v", num)
	}
	for i := 0; i < 10; i++ {
		g.Run(func() error { return f(rand.Intn(10)) })
	}
	err := g.WaitError()
	if err != nil {
		t.Logf("init failed: %v.", err)
	} else {
		t.Logf("init success.")
	}
}
