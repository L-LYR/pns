package model

import "time"

type Target struct {
	*Device
	*App
}

type PushStatus = int8

const (
	Pushable   PushStatus = 1
	UnPushable PushStatus = 2
)

type Device struct {
	ID     string `json:"deviceId" copier:"DeviceId"`
	Os     string `json:"os"`
	Brand  string `json:"brand"`
	Model  string `json:"model"`
	TzName string `json:"tzName"`

	Tokens map[string]string `json:"tokens"`

	InAppPushStatus    PushStatus            `json:"inAppPushStatus"`
	SystemPushStatus   PushStatus            `json:"systemPushStatus"`
	PrivacyPushStatus  PushStatus            `json:"privacyPushStatus"`
	BusinessPushStatus map[string]PushStatus `json:"businessPushStatus"`

	CreateTime      time.Time `json:"createTime"`
	LastActiveTime  time.Time `json:"lastActiveTime"`
	TokenUpdateTime time.Time `json:"tokenUpdateTime"`
	InfoUpdateTime  time.Time `json:"infoUpdateTime"`
}

type App struct {
	ID             int    `json:"appId" copier:"AppId"`
	AppVersion     string `json:"appVersion"`
	PushSdkVersion string `json:"pushSdkVersion"`
	Language       string `json:"language"`
}
