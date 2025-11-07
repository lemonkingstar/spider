package pmaster

import (
	"testing"
	"time"

	"github.com/lemonkingstar/spider/pkg/predis"
)

func TestMater(t *testing.T) {
	_, err := predis.NewDefault(predis.Config{
		Address: "localhost:6379", Database: 0,
	})
	if err != nil {
		t.Fatalf("redis init error: %s", err.Error())
	}
	Run("")
	// check master
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			if Get() {
				t.Log("当前节点类型：master")
			} else {
				t.Log("当前节点类型：normal")
			}
		}
	}
}
