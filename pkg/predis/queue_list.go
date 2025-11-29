package predis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lemonkingstar/spider/pkg/plog"
	"github.com/lemonkingstar/spider/pkg/psafe"
)

var (
	logger = plog.WithField("[PACKET]", "predis")
)

type Message struct {
	Topic string
	Body  []byte
}

// Handler 返回值代表消息是否消费成功
type Handler func(msg *Message) error

type ListQueue interface {
	LPublish(ctx context.Context, msg *Message) error
	LConsume(ctx context.Context, topic string, h Handler)
	LConsumeEx(ctx context.Context, topic string, h Handler)
}

func (c *client) LPublish(ctx context.Context, msg *Message) error {
	return c.cli.LPush(ctx, msg.Topic, msg.Body).Err()
}

func (c *client) LConsume(ctx context.Context, topic string, h Handler) {
	psafe.Go(func() {
		loop := func(topic string, h Handler) error {
			for {
				result, err := c.cli.BRPop(ctx, 0, topic).Result()
				if err != nil {
					return err
				}
				err = h(&Message{
					Topic: result[0],
					Body:  []byte(result[1]),
				})
				if err != nil {
					return err
				}
			}
		}
		for {
			err := loop(topic, h)
			if err != nil {
				logger.Errorf("LConsume error: %s, Wait for reconsume...", err.Error())
			}
			time.Sleep(9 * time.Second)
		}
	})
}

// LConsumeEx 实现ACK机制
// 使用lindex从队列取出消息，如果消费成功再使用rpop删除消息
func (c *client) LConsumeEx(ctx context.Context, topic string, h Handler) {
	psafe.Go(func() {
		loop := func(topic string, h Handler) error {
			for {
				body, err := c.cli.LIndex(ctx, topic, -1).Bytes()
				if err != nil && !errors.Is(err, redis.Nil) {
					return err
				}
				if errors.Is(err, redis.Nil) {
					time.Sleep(time.Second)
					continue
				}
				err = h(&Message{
					Topic: topic,
					Body:  body,
				})
				if err != nil {
					continue
				}
				if err = c.cli.RPop(ctx, topic).Err(); err != nil {
					return err
				}
			}
		}
		for {
			err := loop(topic, h)
			if err != nil {
				logger.Errorf("LConsumeEx error: %s, Wait for reconsume...", err.Error())
			}
			time.Sleep(9 * time.Second)
		}
	})
}
