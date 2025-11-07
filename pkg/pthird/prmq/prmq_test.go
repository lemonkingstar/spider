package prmq

import (
	"fmt"
	"testing"
	"time"
)

func TestRmq(t *testing.T) {
	exType, exName, rtKey, quName := "direct", "xxx-exchange", "xxx-key", "xxx-queue"
	cli, _ := NewDefault("amqp://admin:123456@localhost:5672/", exType, exName, rtKey, quName)
	defer cli.Close()
	// 发送消息
	for i := 0; i < 10; i++ {
		cli.Publish(fmt.Sprintf("Hello world! %d", i))
	}
	time.Sleep(5 * time.Second)
	// 开始消费
	cli.StartConsume(func(b []byte) error {
		t.Log(string(b))
		return nil
	})
	select {}
}

func TestPeriodicConsume(t *testing.T) {
	exType, exName, rtKey, quName := "direct", "xxx-exchange", "xxx", "xxx-queue"
	cli, _ := NewDefault("amqp://admin:123456@localhost:5672/", exType, exName, rtKey, quName)
	defer cli.Close()
	// 发送消息
	go func() {
		start := 0
		for ; start < 100; start++ {
			cli.Publish(fmt.Sprintf("Hello world! %d", start))
			time.Sleep(time.Second)
		}
	}()

	master := false
	go func() {
		for {
			time.Sleep(5 * time.Second)
			master = !master
		}
	}()

	// 判断
	isMaster := func() bool {
		t.Log("=== 间隔检测是否为主节点", master)
		return master
	}

	cli.PeriodicConsume(func(b []byte) error {
		t.Log(string(b))
		return nil
	}, isMaster, 9)
	select {}
}
