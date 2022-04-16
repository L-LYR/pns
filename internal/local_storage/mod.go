package local_storage

import (
	"github.com/L-LYR/pns/internal/model"
)

// NOTICE: mock

func GetAppNameByAppId(id int) (string, bool) {
	if id == 12345 {
		return "test_app_name", true
	}
	return "", false
}

func GetAppConfigByAppId(id int) *model.AppConfig {
	if id == 12345 {
		return &model.AppConfig{
			ID:   12345,
			Name: "test_app_name",
		}
	}
	return nil
}

func GetPusherAuthByAppId(appId int, pusherId model.PusherType) interface{} {
	if appId == 12345 && pusherId == model.MQTTPusher {
		return &model.MQTTConfig{
			PusherKey:    "test_app_name",
			PusherSecret: "test_app_name",
		}
	}
	return nil
}
