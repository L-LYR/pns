package pusher_config

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"github.com/L-LYR/pns/internal/util"
)

func Authorization(ctx context.Context, key string, secret string, clientId string) (bool, string) {
	appId, isPusher, err := util.ParseClientID(clientId)
	if err != nil {
		util.GLog.Errorf(ctx, "%+v", err)
		return false, "unknown client"
	}
	config := &model.MQTTConfig{}
	if err := dao.FindConfigByKey(ctx, model.MQTTPusher, appId, config); err != nil {
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
