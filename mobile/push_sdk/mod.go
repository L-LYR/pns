package push_sdk

import (
	"context"

	"github.com/L-LYR/pns/mobile/push_sdk/net"
	"github.com/L-LYR/pns/mobile/push_sdk/storage"
	"github.com/L-LYR/pns/mobile/push_sdk/util"
)

// TODO: refactor mod.go

var (
	_ctx          = context.Background()
	_GlobalConfig = &storage.Config{
		ClientID:       "pns_target-12345678-12345",
		AppId:          12345,
		Key:            "test_app_name",
		Secret:         "test_app_name",
		DeviceId:       "12345678",
		Broker:         "192.168.137.1",
		Port:           "18830",
		RetryInterval:  1000,
		ConnectTimeout: 60,
		Token:          make(map[string]string),
	}

	_BizHTTPClient    = net.MustNewHTTPClient("http://192.168.137.1:10087")
	_GlobalHTTPClient = net.MustNewHTTPClient("http://192.168.137.1:10086")
	_GlobalMQTTClient = net.MustNewMQTTClient(_ctx, _GlobalConfig)
)

func PushMyself(title string, content string) {
	payload, err := _BizHTTPClient.POST("/push", net.Payload{
		"deviceId": _GlobalConfig.DeviceId,
		"appId":    _GlobalConfig.AppId,
		"title":    title,
		"content":  content,
	})
	if err != nil {
		util.Log("error: %s", err.Error())
	} else {
		util.Log("payload: %+v", payload)
	}
}

func GetToken() {
	payload, err := _GlobalHTTPClient.GET("/token", net.Payload{
		"deviceId": _GlobalConfig.DeviceId,
		"appId":    _GlobalConfig.AppId,
	})
	if err != nil {
		util.Log("error: %s", err.Error())
		return
	}
	token, ok := payload["token"]
	if !ok {
		util.Log("no token")
		return
	}
	_GlobalConfig.Token["self"], ok = token.(string)
	if !ok {
		util.Log("wrong token")
		return
	}
	util.Log("token: %s", token)
}

func RegisterPersonalPushHandler(fn net.MessageHandler) {
	_GlobalMQTTClient.SubscribePersonalPush(fn)
}

func RegisterBroadcastPushHandler(fn net.MessageHandler) {
	_GlobalMQTTClient.SubscribeBroadcastPushHandler(fn)
}

func RegisterGlobalLogHandler(fn util.LogHandler) {
	util.SetGlobalLogHandler(fn)
}
