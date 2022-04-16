package app

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"github.com/L-LYR/pns/internal/util"
)

func Authorization(ctx context.Context, key string, secret string, clientId string) (bool, string) {
	appId, isPusher, err := util.ParseClientID(clientId)
	if err != nil {
		util.GLog.Errorf(ctx, "client id: %s, error: %+v", clientId, err)
		return false, "unknown client"
	}
	config := &model.MQTTConfig{}
	if v, err := dao.FindConfigByKey(ctx, appId, model.MQTTPusher); err != nil {
		util.GLog.Errorf(ctx, "%+v", err)
		return false, "internal error"
	} else if err := v.Struct(config); err != nil {
		util.GLog.Errorf(ctx, "%+v", err)
		return false, "internal error"
	}
	if isPusher && config.PusherKey == key && config.PusherSecret == secret {
		return true, ""
	} else if !isPusher && config.ReceiverKey == key && config.ReceiverSecret == secret {
		return true, ""
	}
	return false, "unauthorized"
}

func Create(ctx context.Context, appName string, appId int) error {
	return dao.CreateApp(ctx, appName, appId)
}

func CreateConfig(ctx context.Context, appId int, pusher model.PusherType, configPointer interface{}) error {
	return dao.CreateConfig(ctx, appId, pusher, configPointer)
}

func QueryConfig(ctx context.Context, appId int, pusher model.PusherType) (map[string]string, error) {
	v, err := dao.FindConfigByKey(ctx, appId, pusher)
	if err != nil {
		return nil, err
	}
	return v.MapStrStr(), nil
}
