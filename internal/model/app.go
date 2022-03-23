package model

type AppConfig struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MQTTConfig struct {
	ID     int    `json:"id"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}
