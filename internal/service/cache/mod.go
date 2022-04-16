package cache

import (
	"context"
	"errors"
	"sync"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
)

/*
	cache features:

	1. very low-frequent insert operation
	2. no update or delete operations
	3. heavy concurrent read load

*/

type _ConfigCache struct {
	appConfig        sync.Map
	mqttPusherConfig sync.Map
	fcmPusherConfig  sync.Map
	apnsPusherConfig sync.Map
}

func (c *_ConfigCache) GetAppNameByAppId(id int) (string, bool) {
	cfg, ok := c.appConfig.Load(id)
	if !ok {
		return "", false
	}
	return cfg.(*model.AppConfig).Name, true
}

func (c *_ConfigCache) GetAppConfigByAppId(id int) (*model.AppConfig, bool) {
	cfg, ok := c.appConfig.Load(id)
	if !ok {
		return nil, false
	}
	return cfg.(*model.AppConfig), true
}

func (c *_ConfigCache) GetMQTTPusherConfigByAppId(id int) (*model.MQTTConfig, bool) {
	cfg, ok := c.mqttPusherConfig.Load(id)
	if !ok {
		return nil, false
	}
	return cfg.(*model.MQTTConfig), true
}

func (c *_ConfigCache) AddAppConfig(config *model.AppConfig) {
	c.appConfig.Store(config.ID, config)
}

func (c *_ConfigCache) AddPusherConfig(config model.PusherConfig) {
	switch config.PusherType() {
	case model.MQTTPusher:
		c.mqttPusherConfig.Store(config.AppId(), config)
	case model.FCMPusher:
		c.fcmPusherConfig.Store(config.AppId(), config)
	case model.APNsPusher:
		c.apnsPusherConfig.Store(config.AppId(), config)
	default:
		panic("unreachable")
	}
}

// used in initialize pusher
func (c *_ConfigCache) RangePusherConfig(t model.PusherType, fn func(appId int, config model.PusherConfig) error) error {
	var err error = nil
	switch t {
	case model.MQTTPusher:
		c.mqttPusherConfig.Range(func(key, value interface{}) bool {
			err = fn(key.(int), value.(model.PusherConfig))
			return err != nil
		})
	case model.FCMPusher:
		return errors.New("not implemented yet")
	case model.APNsPusher:
		return errors.New("not implemented yet")
	default:
		panic("unreachable")
	}
	return err
}

func (c *_ConfigCache) MustLoad(ctx context.Context) {
	appConfigs, err := dao.LoadAllAppConfig(ctx)
	if err != nil {
		panic(err)
	}
	for _, config := range appConfigs {
		c.AddAppConfig(config)
	}
	pusherConfigs, err := dao.LoadAllPusherConfig(ctx)
	if err != nil {
		panic(err)
	}
	for _, config := range pusherConfigs {
		c.AddPusherConfig(config)
	}
}

var (
	Config = &_ConfigCache{}
)
