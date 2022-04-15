package model

type AppConfig struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// members of config are all string type

type MQTTConfig struct {
	PusherKey      string `json:"pusherKey" copier:"must,nopanic"`
	PusherSecret   string `json:"pusherSecret" copier:"must,nopanic"`
	ReceiverKey    string `json:"receiverKey" copier:"must,nopanic"`
	ReceiverSecret string `json:"receiverSecret" copier:"must,nopanic"`
}

type FCMConfig struct {
	Key string `json:"key" copier:"must,nopanic"`
}

type APNsConfig struct {
	AuthKey string `json:"authKey" copier:"must,nopanic"`
	KeyID   string `json:"keyID" copier:"must,nopanic"`
	TeamID  string `json:"teamID" copier:"must,nopanic"`
}
