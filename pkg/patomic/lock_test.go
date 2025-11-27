package patomic

import (
	"sync"
	"testing"
	"time"
)

func TestSimpleLock(t *testing.T) {
	l := NewLock()
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			if l.TryLock() {
				defer l.Unlock()
				t.Logf("Goroutine %d: try lock success.", taskID)
				time.Sleep(500 * time.Millisecond)
				t.Logf("Goroutine %d: business finished.", taskID)
				return
			}

			t.Logf("Goroutine %d: try lock failed, waiting for lock.", taskID)
			if l.LockWithTimeout(1 * time.Second) {
				defer l.Unlock()
				t.Logf("Goroutine %d: waiting for lock success.", taskID)
				time.Sleep(500 * time.Millisecond)
				t.Logf("Goroutine %d: business finished.", taskID)
			} else {
				t.Logf("Goroutine %d: waiting for lock failed, do nothing and exit.", taskID)
			}
		}(i)
	}
	wg.Wait()
}
