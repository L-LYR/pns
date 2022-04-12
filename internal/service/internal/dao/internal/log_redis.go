package internal

import (
	"context"

	"github.com/L-LYR/pns/internal/config"
	"github.com/go-redis/redis/v8"
)

const (
	_RedisDaoConfigName = "redis"
)

type LogRedisDao struct {
	client *redis.Client
}

func NewLogRedisDao(ctx context.Context) (*LogRedisDao, error) {
	cfg := config.MustLoadRedisDaoConfig(ctx, _RedisDaoConfigName)
	return &LogRedisDao{
		client: redis.NewClient(cfg.Options()),
	}, nil
}

// provide a continues single connection
func (dao *LogRedisDao) Conn(ctx context.Context) *redis.Conn {
	return dao.client.Conn(ctx)
}

// run commands
func (dao *LogRedisDao) Client(ctx context.Context) *redis.Client {
	return dao.client.WithContext(ctx)
}

func (dao *LogRedisDao) Shutdown(ctx context.Context) error {
	return dao.client.Close()
}
