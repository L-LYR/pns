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
	ID     string `json:"deviceId" bson:"_id" copier:"DeviceId"`
	Os     string `json:"os" bson:"os"`
	Brand  string `json:"brand" bson:"brand"`
	Model  string `json:"model" bson:"model"`
	TzName string `json:"tzName" bson:"tzName"`

	Tokens map[string]string `json:"tokens" bson:"tokens"`

	InAppPushStatus    PushStatus            `json:"inAppPushStatus" bson:"inAppPushStatus"`
	SystemPushStatus   PushStatus            `json:"systemPushStatus" bson:"systemPushStatus"`
	PrivacyPushStatus  PushStatus            `json:"privacyPushStatus" bson:"privacyPushStatus"`
	BusinessPushStatus map[string]PushStatus `json:"businessPushStatus" bson:"businessPushStatus"`

	CreateTime      time.Time `json:"createTime" bson:"createTime"`
	LastActiveTime  time.Time `json:"lastActiveTime" bson:"lastActiveTime"`
	TokenUpdateTime time.Time `json:"tokenUpdateTime" bson:"tokenUpdateTime"`
	InfoUpdateTime  time.Time `json:"infoUpdateTime" bson:"infoUpdateTime"`
}

type App struct {
	ID             int    `json:"appId" bson:"appId" copier:"AppId"`
	AppVersion     string `json:"appVersion" bson:"appVersion"`
	PushSdkVersion string `json:"pushSdkVersion" bson:"pushSdkVersion"`
	Language       string `json:"language" bson:"language"`
}
