package push_sdk

import (
	"github.com/L-LYR/pns/mobile/push_sdk/storage"
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

func MustNewConfigFromString(s string) *storage.Config {
	return storage.MustNewConfigFromString(s)
}
