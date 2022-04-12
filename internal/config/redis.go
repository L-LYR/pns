package config

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisDaoConfig struct {
	User string
	Pass string
	Host string
	Port string
}

func (c *RedisDaoConfig) Options() *redis.Options {
	return &redis.Options{
		Addr:     c.Host + ":" + c.Port,
		Username: c.User,
		Password: c.Pass,
	}
}

func MustLoadRedisDaoConfig(ctx context.Context, name string) *RedisDaoConfig {
	cfg := &RedisDaoConfig{}
	MustLoadConfig(ctx, "database."+name, cfg)
	return cfg
}
