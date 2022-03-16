package model

type AppConfig struct {
	ID   int    `json:"appId"`
	Name string `json:"name"`
	// TODO: other fields, such as: Key and Secret
}
