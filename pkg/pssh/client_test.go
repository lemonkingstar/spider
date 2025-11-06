package pssh

import (
	"fmt"
	"testing"
	"time"
)

var options = &ClientOptions{
	Host:     "10.10.10.10",
	Password: "dangerous",
}

func TestRun(t *testing.T) {
	client := NewClient(options)
	ret, _ := client.Run("hostname")
	fmt.Println(ret)
}

func TestRun2f(t *testing.T) {
	client := NewClient(options)
	go func() {
		// 5s后退出
		time.Sleep(5 * time.Second)
		client.Close()
	}()
	err := client.Run2f("hostname && sleep 5 && pwd && sleep 2 && echo hello world", func(line string) {
		fmt.Println(line)
	})
	fmt.Println(err)
}

func TestRun2fx(t *testing.T) {
	options.KeyFile = "./private_key"
	client := NewClient(options)
	err := client.Run2f("hostname; sleep 5; pwd; sleep 2; echo hello world", func(line string) {
		fmt.Println(line)
	})
	fmt.Println(err)
}
