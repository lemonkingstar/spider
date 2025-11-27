package psafe

import (
	"testing"
	"time"
)

func TestGolimit(t *testing.T) {
	gl := NewGl(3)
	defer func() {
		// Call Wait() if you want to finish.
		gl.Wait()
		t.Log("all tasks finished.")
	}()

	for i := 0; i < 10; i++ {
		taskID := i
		gl.Run(func() {
			t.Logf("Goroutine %d: start.", taskID)
			time.Sleep(2 * time.Second)
			t.Logf("Goroutine %d: done.", taskID)
		})
	}
	for i := 0; i < 5; i++ {
		taskID := i
		gl.Run(func() {
			t.Logf("Other Goroutine %d: start.", taskID)
			time.Sleep(2 * time.Second)
			t.Logf("Other Goroutine %d: done.", taskID)
		})
	}
}
