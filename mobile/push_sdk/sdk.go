package push_sdk

import (
	"context"
	"fmt"

	"github.com/L-LYR/pns/mobile/push_sdk/net/http"
	"github.com/L-LYR/pns/mobile/push_sdk/net/mqtt"
	"github.com/L-LYR/pns/mobile/push_sdk/storage"
)

var (
	PushSDK *_PushSDK
)

type LogHandler func(fmt string, v ...interface{})

type _PushSDK struct {
	ctx context.Context

	gLogHandler LogHandler
	gHTTPClient *http.Client
	gMQTTClient *mqtt.Client
}

func MustInitialize(cfg *storage.Config, fn LogHandler) {
	ctx := context.Background()
	gHTTPClient := http.MustNewHTTPClient(cfg.Sdk.Inbound.Base)

	if fn == nil {
		fn = DefaultLogHandler()
	}

	storage.SetGlobalConfig(cfg)

	PushSDK = &_PushSDK{
		ctx:         ctx,
		gLogHandler: fn,
		gHTTPClient: gHTTPClient,
	}

	options := mqtt.NewOptions()
	options.SetWithCfg(cfg)
	options.SetTopicSet(mqtt.NewTopicSet(cfg))
	options.SetLogHandler(mqtt.LogHandler(fn))
	options.SetRecvHandler(func(m map[string]interface{}) {
		PushSDK.ReportLog(m)
	})
	options.SetShowHandler(func(m map[string]interface{}) {
		PushSDK.ReportLog(m)
	})

	PushSDK.gMQTTClient = mqtt.MustNewMQTTClient(options)
}

func DefaultLogHandler() LogHandler {
	return func(s string, v ...interface{}) {
		fmt.Printf(s, v...)
	}
}

func (sdk *_PushSDK) ReportLog(payload http.Payload) {
	_, err := sdk.gHTTPClient.POST("/log", payload)
	if err != nil {
		sdk.gLogHandler("Error: %s", err.Error())
	}
}

func (sdk *_PushSDK) UpdateTargetInfo() {
	_, err := sdk.gHTTPClient.POST("/target", http.Payload{
		"deviceId":           storage.GlobalConfig.DeviceId,
		"os":                 "windows",
		"brand":              "chrome",
		"model":              "chrome",
		"tzName":             "Asia/Shanghai",
		"appId":              storage.GlobalConfig.App.ID,
		"appVersion":         storage.GlobalConfig.App.Version,
		"pushSDKVersion":     storage.GlobalConfig.Sdk.Version,
		"language":           "cn",
		"inAppPushStatus":    1,
		"systemPushStatus":   1,
		"privacyPushStatus":  1,
		"businessPushStatus": make(map[string]int),
	})
	if err != nil {
		sdk.gLogHandler("Error: %s", err.Error())
	} else {
		sdk.gLogHandler("Info: Succeed to update info")
	}
}

func (sdk *_PushSDK) GetToken() {
	payload, err := sdk.gHTTPClient.GET("/token", http.Payload{
		"deviceId": storage.GlobalConfig.DeviceId,
		"appId":    storage.GlobalConfig.App.ID,
	})
	if err != nil {
		sdk.gLogHandler("Error: %s", err.Error())
		return
	}
	token, ok := payload["token"]
	if !ok {
		sdk.gLogHandler("Error: no token")
		return
	}
	storage.GlobalConfig.Token["self"], ok = token.(string)
	if !ok {
		sdk.gLogHandler("Error: wrong token")
		return
	}
	sdk.gLogHandler("Info: token: %s", token)
}

type MessageHandler mqtt.MessageHandler

func (sdk *_PushSDK) RegisterPersonalPushHandler(fn MessageHandler) {
	sdk.gMQTTClient.SubscribePersonalPush(mqtt.MessageHandler(fn))
}

func (sdk *_PushSDK) RegisterBroadcastPushHandler(fn MessageHandler) {
	sdk.gMQTTClient.SubscribeBroadcastPush(mqtt.MessageHandler(fn))
}
