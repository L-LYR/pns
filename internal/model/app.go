package model

type AppConfig struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MQTTConfig struct {
	ID             int    `json:"id"`
	PusherKey      string `json:"pusher_key"`
	PusherSecret   string `json:"pusher_secret"`
	ReceiverKey    string `json:"receiver_key"`
	ReceiverSecret string `json:"receiver_secret"`
}
