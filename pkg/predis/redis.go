package predis

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
)

// Config define redis config
type Config struct {
	Address          string
	Password         string
	Database         int
	MasterName       string
	SentinelPassword string
}

// New returns new redis client from config
func New(cfg Config) (Client, error) {
	var client Client
	if cfg.MasterName == "" {
		option := &redis.Options{
			Addr:     cfg.Address,
			Password: cfg.Password,
			DB:       cfg.Database,
			//PoolSize: cfg.MaxOpenConns,
		}
		client = NewClient(option)
	} else {
		hosts := strings.Split(cfg.Address, ",")
		option := &redis.FailoverOptions{
			MasterName:       cfg.MasterName,
			SentinelAddrs:    hosts,
			Password:         cfg.Password,
			DB:               cfg.Database,
			SentinelPassword: cfg.SentinelPassword,
		}
		client = NewFailoverClient(option)
	}

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, err
}

// IsNilErr returns whether err is nil error
func IsNilErr(err error) bool {
	return redis.Nil == err
}
