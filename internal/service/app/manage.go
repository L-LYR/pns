package app

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/cache"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"github.com/L-LYR/pns/internal/util"
)

func ACLCheck(ctx context.Context, key string, clientId string) (bool, string) {
	if clientId == util.GetRootClientID() &&
		key == util.GetRootClientUser() {
		return true, ""
	}

	appId, err := util.ParseClientID(clientId)
	if err != nil {
		util.GLog.Errorf(ctx, "client id: %s, error: %+v", clientId, err)
		return false, "unknown client"
	}
	config, ok := cache.Config.GetMQTTPusherConfigByAppId(appId)
	if !ok {
		return false, "cache miss"
	}
	if config.ReceiverKey == key {
		return true, ""
	}
	return false, "unauthorized"
}

func Authorization(ctx context.Context, key string, secret string, clientId string) (bool, string) {
	if clientId == util.GetRootClientID() &&
		key == util.GetRootClientUser() &&
		secret == util.GetRootClientPass() {
		return true, ""
	}

	appId, err := util.ParseClientID(clientId)
	if err != nil {
		util.GLog.Errorf(ctx, "client id: %s, error: %+v", clientId, err)
		return false, "unknown client"
	}
	config, ok := cache.Config.GetMQTTPusherConfigByAppId(appId)
	if !ok {
		return false, "cache miss"
	}
	if config.ReceiverKey == key && config.ReceiverSecret == secret {
		return true, ""
	}
	return false, "unauthorized"
}

func Create(ctx context.Context, appName string, appId int) error {
	cfg := &model.AppConfig{
		ID:   appId,
		Name: appName,
	}
	if err := dao.CreateApp(ctx, cfg); err != nil {
		return err
	}
	cache.Config.AddAppConfig(cfg)
	return nil
}

func CreateConfig(ctx context.Context, config model.PusherConfig) error {
	if err := dao.CreateConfig(ctx, config); err != nil {
		return err
	}
	cache.Config.AddPusherConfig(config)
	return nil
}

func QueryConfig(ctx context.Context, appId int, pusher model.PusherType) (map[string]string, error) {
	v, err := dao.FindConfigByKey(ctx, appId, pusher)
	if err != nil {
		return nil, err
	}
	return v.MapStrStr(), nil
}
