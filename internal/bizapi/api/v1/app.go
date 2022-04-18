package v1

import "github.com/gogf/gf/v2/frame/g"

type CreateAppReq struct {
	g.Meta  `path:"/app" method:"post"`
	AppId   int    `json:"appId" dc:"app id, must be unique" v:"required|app-not-exist#app id is required"`
	AppName string `json:"appName" dc:"app name, must be unique" v:"required#app name is required"`
}

type CreateAppRes struct{}

type CreateAppConfigReq struct {
	g.Meta     `path:"/app/config" method:"post"`
	AppId      int               `json:"appId" dc:"app id" v:"required|app-exist#app id is required"`
	ConfigName string            `json:"configName" dc:"pusher config name, must be fcm,apns" v:"required#pusher name is required"`
	Config     map[string]string `json:"config" dc:"pusher config, only for fcm and apns"`
}

type CreateAppConfigRes struct{}

type OpenMQTTReq struct {
	g.Meta `path:"/app/config/mqtt" method:"post"`
	AppId  int `json:"appId" dc:"app id" v:"required|app-exist#app id is required"`
}

type MQTTConfig struct {
	PusherKey      string `json:"pusherKey" dc:"app key, used to push"`
	PusherSecret   string `json:"pusherSecret" dc:"app secret, used to push"`
	ReceiverKey    string `json:"receiverKey" dc:"client key, used to receive"`
	ReceiverSecret string `json:"receiverSecret" dc:"client secret, used to receive"`
}

type OpenMQTTRes struct {
	*MQTTConfig `json:"mqttConfig" dc:"mqtt config"`
}

type QueryAppConfigReq struct {
	g.Meta     `path:"/app/config" method:"get"`
	AppId      int    `json:"appId" dc:"app id" v:"required|app-exist#app id is required"`
	ConfigName string `json:"configName" dc:"pusher config name, must be mqtt,gcm,apns" v:"required#pusher name is required"`
}

type QueryAppConfigRes struct {
	Config interface{} `json:"config"`
}
