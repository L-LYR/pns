package storage

import (
	"time"

	"github.com/L-LYR/pns/mobile/push_sdk/util"
	jsoniter "github.com/json-iterator/go"
)

// TODO: implement persistance

type Config struct {
	ClientId string
	DeviceId string
	Token    map[string]string
	Settings
}

func (c *Config) GetAddress() string {
	return "tcp://" + c.Settings.Sdk.MQTT.Broker + ":" + c.Settings.Sdk.MQTT.Port
}

func (c *Config) GetRetryInterval() time.Duration {
	return time.Duration(c.Settings.Sdk.MQTT.RetryInterval) * time.Millisecond
}

func (c *Config) GetConnectTimeout() time.Duration {
	return time.Duration(c.Settings.Sdk.MQTT.ConnectTimeout) * time.Second
}

func (c *Config) GetDeviceId() string { return c.DeviceId }

func (c *Config) GetAppId() int { return c.App.ID }

func MustNewConfigFromString(s string) *Config {
	c := &Config{}
	if err := jsoniter.UnmarshalFromString(s, c); err != nil {
		panic(err)
	}
	c.DeviceId = util.GenerateDeviceId()
	c.ClientId = util.GenerateClientId("pns-target", c.DeviceId, c.App.ID)
	c.Token = make(map[string]string)
	return c
}

func DefaultConfig() *Config {
	deviceId := util.GenerateDeviceId()
	appId := 12345
	appName := "test_app_name"
	return &Config{
		ClientId: util.GenerateClientId("pns-target", deviceId, appId),
		DeviceId: deviceId,
		Token:    make(map[string]string),
		Settings: Settings{
			App: APPSettings{
				ID:      appId,
				Key:     appName,
				Secret:  appName,
				Version: "0.0.1",
			},
			Sdk: SDKSettings{
				Version: "0.0.1",
				MQTT: MQTTSettings{
					Broker:         "192.168.137.1",
					Port:           "18830",
					RetryInterval:  1000,
					ConnectTimeout: 60,
				},
				Inbound: InboundSettings{
					Base: "192.168.137.1:10086",
				},
			},
		},
	}
}

var (
	GlobalConfig = &Config{}
)

func SetGlobalConfig(c *Config) {
	GlobalConfig = c
}
