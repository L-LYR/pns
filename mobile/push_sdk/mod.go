package push_sdk

import (
	"context"

	"github.com/L-LYR/pns/mobile/push_sdk/net"
	"github.com/L-LYR/pns/mobile/push_sdk/storage"
	"github.com/L-LYR/pns/mobile/push_sdk/util"
)

// TODO: refactor mod.go

const (
	PUSH_SDK_VERSION = "0.0.1"
)

var (
	_ctx          = context.Background()
	_GlobalConfig = _GenerateGlobalConfig()

	_BizHTTPClient    = net.MustNewHTTPClient("http://192.168.137.1:10087")
	_GlobalHTTPClient = net.MustNewHTTPClient("http://192.168.137.1:10086")
	_GlobalMQTTClient = net.MustNewMQTTClient(_ctx, _GlobalConfig)
)

func _GenerateGlobalConfig() *storage.Config {
	deviceId := util.GenerateDeviceId()
	appId := 12345
	appName := "test_app_name"
	return &storage.Config{
		ClientID:       util.GenerateClientId("pns-target", deviceId, appId),
		AppId:          appId,
		Key:            appName,
		Secret:         appName,
		DeviceId:       deviceId,
		Broker:         "192.168.137.1",
		Port:           "18830",
		RetryInterval:  1000,
		ConnectTimeout: 60,
		Token:          make(map[string]string),
	}
}

func GetConfig() *storage.Config {
	return _GlobalConfig
}

func UpdateTargetInfo() {
	_, err := _GlobalHTTPClient.POST("/target", net.Payload{
		"deviceId":           _GlobalConfig.DeviceId,
		"os":                 "windows",
		"brand":              "chrome",
		"model":              "chrome",
		"tzName":             "Asia/Shanghai",
		"appId":              _GlobalConfig.AppId,
		"appVersion":         "3.3.3",
		"pushSDKVersion":     PUSH_SDK_VERSION,
		"language":           "cn",
		"inAppPushStatus":    1,
		"systemPushStatus":   1,
		"privacyPushStatus":  1,
		"businessPushStatus": make(map[string]int),
	})
	if err != nil {
		util.Log("Error: %s", err.Error())
	} else {
		util.Log("Info: Succeed to update info")
	}
}

func PushMyself(title string, content string) {
	payload, err := _BizHTTPClient.POST("/push", net.Payload{
		"deviceId": _GlobalConfig.DeviceId,
		"appId":    _GlobalConfig.AppId,
		"title":    title,
		"content":  content,
	})
	if err != nil {
		util.Log("Error: %s", err.Error())
	} else {
		if s, ok := payload["pushTaskId"]; ok {
			util.Log("Info: Task ID: %+v", s)
		} else {
			util.Log("Error: no push task id")
		}
	}
}

func GetToken() {
	payload, err := _GlobalHTTPClient.GET("/token", net.Payload{
		"deviceId": _GlobalConfig.DeviceId,
		"appId":    _GlobalConfig.AppId,
	})
	if err != nil {
		util.Log("Error: %s", err.Error())
		return
	}
	token, ok := payload["token"]
	if !ok {
		util.Log("Error: no token")
		return
	}
	_GlobalConfig.Token["self"], ok = token.(string)
	if !ok {
		util.Log("Error: wrong token")
		return
	}
	util.Log("Info: token: %s", token)
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
