package model

type AppConfig struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MQTTConfig struct {
	ID             int    `json:"id"`
	PusherKey      string `json:"pusherKey"`
	PusherSecret   string `json:"pusherSecret"`
	ReceiverKey    string `json:"receiverKey"`
	ReceiverSecret string `json:"receiverSecret"`
}
