package psafe

import (
	"testing"
	"time"
)

func TestGolimit(t *testing.T) {
	gl := NewGl(3)
	defer gl.Close()

	for i := 0; i < 10; i++ {
		taskID := i
		gl.Run(func() {
			t.Logf("Goroutine %d: start.", taskID)
			time.Sleep(2 * time.Second)
			t.Logf("Goroutine %d: done.", taskID)
		})
	}

	// Call Wait() if you want to finish.
	gl.Wait()
	t.Log("all tasks finished.")
}
