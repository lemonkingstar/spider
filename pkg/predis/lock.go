package predis

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotObtained = errors.New("redis lock not obtained")
)

type DistributedLock interface {
	Lock(ctx context.Context, key string) (locked bool, err error)
	Unlock(ctx context.Context, key string) (err error)
	LockUntil(ctx context.Context, key string, expiration time.Duration) (locked bool, err error)
}

// 关于lock续租的问题 参考
// https://learnku.com/articles/46788
