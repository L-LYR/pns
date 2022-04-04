package config

import (
	"context"
	"fmt"
)

type MongoDaoConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string // database name
}

func (c *MongoDaoConfig) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", c.User, c.Pass, c.Host, c.Port)
}

func MustLoadMongoDaoConfig(ctx context.Context, name string) *MongoDaoConfig {
	cfg := &MongoDaoConfig{}
	MustLoadConfig(ctx, "database."+name, cfg)
	return cfg
}
