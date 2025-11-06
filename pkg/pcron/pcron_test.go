package pcron

import "testing"

func TestCronJob(t *testing.T) {
	s := New()
	s.AddJob("* * * * *", func() {
		t.Log("hello world!")
	})
	s.AddJobDetail("job-test", "* * * * *", true, true, func() {
		t.Log("hello job test!")
	})
	s.Start()
	select {}
}
