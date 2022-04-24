package push_sdk

import (
	"github.com/L-LYR/pns/mobile/push_sdk/storage"
	"github.com/L-LYR/pns/mobile/push_sdk/util"
)

func GetConfig() *storage.Config { return storage.GlobalConfig }

func UpdateTargetInfo() { PushSDK.UpdateTargetInfo() }

func GetToken() { PushSDK.GetToken() }

func RegisterPersonalPushHandler(fn MessageHandler) {
	PushSDK.RegisterPersonalPushHandler(fn)
}

func RegisterBroadcastPushHandler(fn MessageHandler) {
	PushSDK.RegisterBroadcastPushHandler(fn)
}

func MustNewConfigFromString(s string, deviceId string) *storage.Config {
	return storage.MustNewConfigFromString(s, deviceId)
}

func GenerateUUIDAsDeviceId() string {
	return util.GenerateDeviceId()
}
