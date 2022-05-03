package model

import "errors"

type AppConfig struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PusherType int8

const (
	MQTTPusher PusherType = 1
	FCMPusher  PusherType = 2
	APNsPusher PusherType = 3
)

func (t PusherType) Name() string {
	switch t {
	case MQTTPusher:
		return "mqtt"
	case FCMPusher:
		return "fcm"
	case APNsPusher:
		return "apns"
	default:
		panic("unreachable")
	}
}

func ParsePusherType(name string) (PusherType, error) {
	switch name {
	case "mqtt":
		return MQTTPusher, nil
	case "fcm":
		return FCMPusher, nil
	case "apns":
		return APNsPusher, nil
	}
	return 0, errors.New("unknown pusher type")
}

type PusherConfig interface {
	AppId() int
	PusherType() PusherType
}

func NewEmptyPusherConfig(appId int, t PusherType) PusherConfig {
	switch t {
	case MQTTPusher:
		return &MQTTConfig{ID: appId}
	case FCMPusher:
		return &FCMConfig{ID: appId}
	case APNsPusher:
		return &APNsConfig{ID: appId}
	default:
		panic("unreachable")
	}
}

type MQTTConfig struct {
	ID             int    `json:"-"`
	ReceiverKey    string `json:"receiverKey" copier:"must,nopanic"`
	ReceiverSecret string `json:"receiverSecret" copier:"must,nopanic"`
}

func (c *MQTTConfig) AppId() int             { return c.ID }
func (c *MQTTConfig) PusherType() PusherType { return MQTTPusher }

type FCMConfig struct {
	ID  int    `json:"-"`
	Key string `json:"key" copier:"must,nopanic"`
}

func (c *FCMConfig) AppId() int             { return c.ID }
func (c *FCMConfig) PusherType() PusherType { return FCMPusher }

type APNsConfig struct {
	ID      int    `json:"-"`
	AuthKey string `json:"authKey" copier:"must,nopanic"`
	KeyID   string `json:"keyID" copier:"must,nopanic"`
	TeamID  string `json:"teamID" copier:"must,nopanic"`
}

func (c *APNsConfig) AppId() int             { return c.ID }
func (c *APNsConfig) PusherType() PusherType { return APNsPusher }
