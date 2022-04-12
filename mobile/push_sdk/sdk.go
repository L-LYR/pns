package push_sdk

import (
	"context"

	"github.com/L-LYR/pns/mobile/push_sdk/net/http"
	"github.com/L-LYR/pns/mobile/push_sdk/net/mqtt"
	"github.com/L-LYR/pns/mobile/push_sdk/storage"
	"github.com/L-LYR/pns/mobile/push_sdk/util"
)

const (
	PUSH_SDK_VERSION = "0.0.1"
)

var (
	PushSDK *_PushSDK
)

type LogHandler func(fmt string, v ...interface{})

type _PushSDK struct {
	ctx context.Context
	cfg *storage.Config

	gLogHandler LogHandler
	gHTTPClient *http.Client
	gMQTTClient *mqtt.Client
}

func MustInitialize(cfg *storage.Config, fn LogHandler) {
	ctx := context.Background()
	gHTTPClient := http.MustNewHTTPClient("http://192.168.137.1:10086")

	PushSDK = &_PushSDK{
		ctx:         ctx,
		cfg:         cfg,
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

func DefaultConfig() *storage.Config {
	deviceId := util.GenerateDeviceId()
	appId := 12345
	appName := "test_app_name"
	return &storage.Config{
		ClientId:       util.GenerateClientId("pns-target", deviceId, appId),
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

func (sdk *_PushSDK) GetConfig() *storage.Config {
	return sdk.cfg
}

func (sdk *_PushSDK) ReportLog(payload http.Payload) {
	_, err := sdk.gHTTPClient.POST("/log", payload)
	if err != nil {
		sdk.gLogHandler("Error: %s", err.Error())
	}
}

func (sdk *_PushSDK) UpdateTargetInfo() {
	_, err := sdk.gHTTPClient.POST("/target", http.Payload{
		"deviceId":           PushSDK.cfg.DeviceId,
		"os":                 "windows",
		"brand":              "chrome",
		"model":              "chrome",
		"tzName":             "Asia/Shanghai",
		"appId":              PushSDK.cfg.AppId,
		"appVersion":         "3.3.3",
		"pushSDKVersion":     PUSH_SDK_VERSION,
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
		"deviceId": sdk.cfg.DeviceId,
		"appId":    sdk.cfg.AppId,
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
	PushSDK.cfg.Token["self"], ok = token.(string)
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
