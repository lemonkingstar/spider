package predis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Pipeliner interface {
	redis.StatefulCmdable
	Close() error
	Discard() error
	Exec(ctx context.Context) ([]redis.Cmder, error)
}
