package predis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type PubSub interface {
	Channel(opts ...redis.ChannelOption) <-chan *redis.Message
	ChannelSize(size int) <-chan *redis.Message
	ChannelWithSubscriptions(_ context.Context, size int) <-chan interface{}
	Close() error
	PSubscribe(ctx context.Context, patterns ...string) error
	PUnsubscribe(ctx context.Context, patterns ...string) error
	Ping(ctx context.Context, payload ...string) error
	Receive(ctx context.Context) (interface{}, error)
	ReceiveMessage(ctx context.Context) (*redis.Message, error)
	ReceiveTimeout(ctx context.Context, timeout time.Duration) (interface{}, error)
	String() string
	Subscribe(ctx context.Context, channels ...string) error
	Unsubscribe(ctx context.Context, channels ...string) error
}
