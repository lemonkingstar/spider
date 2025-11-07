package predis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestRedisCommand(t *testing.T) {
	MyClient()
}

func MyClient() {
	client := NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       10,
	})

	DBOps2(client)
}

func DBOps(cli Client) {
	ctx := context.Background()
	key := "mykey"
	listName := "mylist"
	listName2 := "mylist2"
	hashKey := "myHashKey"
	setKey := "mySetKey"

	pipe := cli.Pipeline()
	pipe.Set(ctx, "aaa", 99, 0)
	pipe.Get(ctx, "aaa")
	vals, err := pipe.Exec(ctx)
	checkErr(err)
	fmt.Println("Pipeline", vals)

	err = cli.Set(ctx, key, "Hello,man!", 0).Err()
	checkErr(err)

	intVal, err := cli.Exists(ctx, key).Result()
	checkErr(err)
	fmt.Println("Exists", intVal)

	strVal, err := cli.Get(ctx, key).Result()
	checkErr(err)
	fmt.Println("Get", key, strVal)

	interfVal, err := cli.Eval(ctx, "return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}", []string{"key1", "key2"}, "arg1", "arg2").Result()
	checkErr(err)
	fmt.Println("Eval:", interfVal)

	statusVal, err := cli.Ping(ctx).Result()
	checkErr(err)
	fmt.Println("Ping:", statusVal)

	cli.Set(ctx, "key1", "value111", 0)
	cli.Set(ctx, "key2", "value222", 0)
	interfSliVal, err := cli.MGet(ctx, "key1", "key2").Result()
	checkErr(err)
	fmt.Println("MGet:", interfSliVal)

	intVal, err = cli.Del(ctx, "key1", "key2").Result()
	checkErr(err)
	fmt.Println("Del:", intVal)

	sub := cli.Subscribe(ctx, "channels")
	checkErr(err)

	go func() {
		time.Sleep(time.Second)
		cli.Publish(ctx, "channels", "hello,a subscribe test")
	}()
	msg, err := sub.ReceiveMessage(ctx)
	checkErr(err)
	fmt.Println("ReceiveMessage:", msg)

	err = sub.Unsubscribe(ctx, "channels")
	checkErr(err)

	err = sub.Close()
	checkErr(err)

	go func() {
		// 订阅消费 Channel
		sub2 := cli.Subscribe(ctx, "channels")
		defer func() {
			sub2.Close()
		}()
		for msg := range sub.Channel() {
			fmt.Println(msg)
		}
	}()

	intVal, err = cli.Incr(ctx, "key").Result()
	checkErr(err)
	fmt.Println("Incr:", "key", intVal)

	intVal, err = cli.LPush(ctx, listName, "111").Result()
	checkErr(err)
	fmt.Println("LPush", listName, intVal)

	strSliVal, err := cli.BRPop(ctx, time.Second*30, listName).Result()
	checkErr(err)
	fmt.Println("BRPop", listName, strSliVal)

	cli.LPush(ctx, listName, "333")
	strVal, err = cli.BRPopLPush(ctx, listName, listName2, time.Second).Result()
	checkErr(err)
	fmt.Println("BRPopLPush", strVal)

	intVal, err = cli.LLen(ctx, listName2).Result()
	checkErr(err)
	fmt.Println("LLen", strVal)

	statusVal, err = cli.LTrim(ctx, listName2, 0, 100).Result()
	checkErr(err)
	fmt.Println("LTrim", strVal)

	intVal, err = cli.HIncrBy(ctx, hashKey, "field", 5).Result()
	checkErr(err)
	fmt.Println("HIncrBy", intVal)

	intVal, err = cli.SAdd(ctx, setKey, "m1", "m2", "m3").Result()
	checkErr(err)
	fmt.Println("SAdd", intVal)

	intVal, err = cli.SRem(ctx, setKey, "m2").Result()
	checkErr(err)
	fmt.Println("SRem", intVal)

	strSliVal, err = cli.SMembers(ctx, setKey).Result()
	checkErr(err)
	fmt.Println("SAdd", strSliVal)
}

func checkErr(err error) {
	if err != nil {
		if err == redis.Nil {
			return
		}
		panic(err)
	}
}

func DBOps2(cli Client) {
	ctx := context.Background()
	v, err := cli.HGet(ctx, "test", "aaa").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
	if err := cli.HSet(ctx, "test", "aaa", "111").Err(); err != nil {
		fmt.Println(err)
	}
	v, err = cli.HGet(ctx, "test", "bbb").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
}

func TestListQueue(t *testing.T) {
	cli, err := NewDefault(Config{
		Address: "10.10.6.65:6379", Database: 2,
	})
	if err != nil {
		t.Fatalf("redis init error: %s", err.Error())
	}

	topic := "phoenix:queue"
	h := func(msg *Message) error {
		t.Log(string(msg.Body))
		return nil
	}
	cli.LConsume(context.Background(), topic, h)
	go func() {
		for i := 0; i < 100; i++ {
			msg := &Message{
				Topic: topic, Body: []byte(fmt.Sprintf("test message %d", i)),
			}
			cli.LPublish(context.Background(), msg)
		}
	}()
	select {}
}
