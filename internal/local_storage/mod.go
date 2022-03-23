package local_storage

import (
	"github.com/L-LYR/pns/internal/model"
)

// NOTICE: mock

func GetAppNameByAppId(id int) string {
	return "test_app_name"
}

func GetAppConfigByAppId(id int) *model.AppConfig {
	return &model.AppConfig{
		ID:   12345,
		Name: "test_app_name",
	}
}

func GetPusherAuthByAppId(appId int, pusherId model.PusherType) interface{} {
	return &model.MQTTConfig{
		ID:     12345,
		Key:    "test_app_name",
		Secret: "test_app_name",
	}
}
