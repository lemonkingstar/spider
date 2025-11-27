package predis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client is the interface for redis client
type Client interface {
	Subscribe(ctx context.Context, channels ...string) PubSub
	PSubscribe(ctx context.Context, channels ...string) PubSub
	Client() *redis.Client

	Command
	DistributedLock
	ListQueue
}

type client struct {
	cli *redis.Client
}

// NewClient returns a client to the Redis Server specified by Options
func NewClient(opt *redis.Options) Client {
	return &client{
		cli: redis.NewClient(opt),
	}
}

// NewFailoverClient returns a Redis client that uses Redis Sentinel for automatic failover
func NewFailoverClient(failoverOpt *redis.FailoverOptions) Client {
	return &client{
		cli: redis.NewFailoverClient(failoverOpt),
	}
}

func (c *client) Client() *redis.Client {
	return c.cli
}

func (c *client) Subscribe(ctx context.Context, channels ...string) PubSub {
	return c.cli.Subscribe(ctx, channels...)
}

func (c *client) PSubscribe(ctx context.Context, channels ...string) PubSub {
	return c.cli.PSubscribe(ctx, channels...)
}

func (c *client) Pipeline() Pipeliner {
	return c.cli.Pipeline()
}

func (c *client) BRPop(ctx context.Context, timeout time.Duration, keys ...string) StringSliceResult {
	return c.cli.BRPop(ctx, timeout, keys...)
}

func (c *client) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) StringResult {
	return c.cli.BRPopLPush(ctx, source, destination, timeout)
}

func (c *client) Close() error {
	return c.cli.Close()
}

func (c *client) Del(ctx context.Context, keys ...string) IntResult {
	return c.cli.Del(ctx, keys...)
}

func (c *client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) Result {
	return c.cli.Eval(ctx, script, keys, args...)
}

func (c *client) Exists(ctx context.Context, keys ...string) IntResult {
	return c.cli.Exists(ctx, keys...)
}

func (c *client) Expire(ctx context.Context, key string, expiration time.Duration) BoolResult {
	return c.cli.Expire(ctx, key, expiration)
}

func (c *client) FlushDB(ctx context.Context) StatusResult {
	return c.cli.FlushDB(ctx)
}

func (c *client) Get(ctx context.Context, key string) StringResult {
	return c.cli.Get(ctx, key)
}

func (c *client) HDel(ctx context.Context, key string, fields ...string) IntResult {
	return c.cli.HDel(ctx, key, fields...)
}

func (c *client) HGet(ctx context.Context, key, field string) StringResult {
	return c.cli.HGet(ctx, key, field)
}

func (c *client) HGetAll(ctx context.Context, key string) StringStringMapResult {
	return c.cli.HGetAll(ctx, key)
}

func (c *client) HIncrBy(ctx context.Context, key, field string, incr int64) IntResult {
	return c.cli.HIncrBy(ctx, key, field, incr)
}

func (c *client) HKeys(ctx context.Context, key string) StringSliceResult {
	return c.cli.HKeys(ctx, key)
}

func (c *client) HMGet(ctx context.Context, key string, fields ...string) SliceResult {
	return c.cli.HMGet(ctx, key, fields...)
}

func (c *client) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ScanResult {
	return c.cli.HScan(ctx, key, cursor, match, count)
}

func (c *client) Scan(ctx context.Context, cursor uint64, match string, count int64) ScanResult {
	return c.cli.Scan(ctx, cursor, match, count)
}

func (c *client) HSet(ctx context.Context, key string, values ...interface{}) IntResult {
	return c.cli.HSet(ctx, key, values...)
}

func (c *client) Incr(ctx context.Context, key string) IntResult {
	return c.cli.Incr(ctx, key)
}

func (c *client) Keys(ctx context.Context, pattern string) StringSliceResult {
	return c.cli.Keys(ctx, pattern)
}

func (c *client) LLen(ctx context.Context, key string) IntResult {
	return c.cli.LLen(ctx, key)
}

func (c *client) LPush(ctx context.Context, key string, values ...interface{}) IntResult {
	return c.cli.LPush(ctx, key, values...)
}

func (c *client) LRange(ctx context.Context, key string, start, stop int64) StringSliceResult {
	return c.cli.LRange(ctx, key, start, stop)
}

func (c *client) LRem(ctx context.Context, key string, count int64, value interface{}) IntResult {
	return c.cli.LRem(ctx, key, count, value)
}

func (c *client) LTrim(ctx context.Context, key string, start, stop int64) StatusResult {
	return c.cli.LTrim(ctx, key, start, stop)
}

func (c *client) LIndex(ctx context.Context, key string, index int64) StringResult {
	return c.cli.LIndex(ctx, key, index)
}

func (c *client) MGet(ctx context.Context, keys ...string) SliceResult {
	return c.cli.MGet(ctx, keys...)
}

func (c *client) MSet(ctx context.Context, values ...interface{}) StatusResult {
	return c.cli.MSet(ctx, values...)
}

func (c *client) Ping(ctx context.Context) StatusResult {
	return c.cli.Ping(ctx)
}

func (c *client) Publish(ctx context.Context, channel string, message interface{}) IntResult {
	return c.cli.Publish(ctx, channel, message)
}

func (c *client) Rename(ctx context.Context, key, newkey string) StatusResult {
	return c.cli.Rename(ctx, key, newkey)
}

func (c *client) RenameNX(ctx context.Context, key, newkey string) BoolResult {
	return c.cli.RenameNX(ctx, key, newkey)
}

func (c *client) RPop(ctx context.Context, key string) StringResult {
	return c.cli.RPop(ctx, key)
}

func (c *client) RPopLPush(ctx context.Context, source, destination string) StringResult {
	return c.cli.RPopLPush(ctx, source, destination)
}

func (c *client) RPush(ctx context.Context, key string, values ...interface{}) IntResult {
	return c.cli.RPush(ctx, key, values...)
}

func (c *client) SAdd(ctx context.Context, key string, members ...interface{}) IntResult {
	return c.cli.SAdd(ctx, key, members...)
}

func (c *client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) StatusResult {
	return c.cli.Set(ctx, key, value, expiration)
}

func (c *client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) BoolResult {
	return c.cli.SetNX(ctx, key, value, expiration)
}

func (c *client) TxPipeline(ctx context.Context) Pipeliner {
	return c.cli.TxPipeline()
}

func (c *client) Discard(ctx context.Context, pipe Pipeliner) error {
	return pipe.Discard()
}
func (c *client) MSetNX(ctx context.Context, values ...interface{}) BoolResult {
	return c.cli.MSetNX(ctx, values...)
}

func (c *client) SMembers(ctx context.Context, key string) StringSliceResult {
	return c.cli.SMembers(ctx, key)
}

func (c *client) SIsMember(ctx context.Context, key string, member interface{}) BoolResult {
	return c.cli.SIsMember(ctx, key, member)
}

func (c *client) SRem(ctx context.Context, key string, members ...interface{}) IntResult {
	return c.cli.SRem(ctx, key, members...)
}

func (c *client) TTL(ctx context.Context, key string) DurationResult {
	return c.cli.TTL(ctx, key)
}

func (c *client) Lock(ctx context.Context, key string) (locked bool, err error) {
	locked, err = c.SetNX(ctx, key, time.Now(), time.Minute).Result()
	if err != nil {
		return false, err
	} else if !locked {
		return false, ErrNotObtained
	}
	return locked, nil
}

func (c *client) Unlock(ctx context.Context, key string) (err error) {
	_, err = c.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// LockUntil redis设置超时锁
// usage:
// b, _ := client.Lock(ctx, lockKey, time.Minute);
// defer client.Unlock(ctx, lockKey) // 如果需要在指定时间内都不想要执行业务逻辑，也可以不主动释放，待自动释放即可
func (c *client) LockUntil(ctx context.Context, key string, expiration time.Duration) (locked bool, err error) {
	locked, err = c.SetNX(ctx, key, time.Now(), expiration).Result()
	if err != nil {
		return false, err
	} else if !locked {
		return false, ErrNotObtained
	}
	return locked, nil
}
