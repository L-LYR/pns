package model

import (
	"time"
)

type Target struct {
	*Device `bson:"inline"`
	*App    `bson:"inline"`
	Tokens  *TokenSet `json:"tokens" bson:"tokens"`
}

func (l *Target) Equal(r *Target) bool {
	return l.App.Equal(r.App) && l.Device.Equal(r.Device)
}

type PushStatus = int8

const (
	Pushable   PushStatus = 1
	UnPushable PushStatus = 2
)

type Device struct {
	ID     string `json:"deviceId" bson:"_id" copier:"DeviceId"`
	Os     string `json:"os,omitempty" bson:"os"`
	Brand  string `json:"brand,omitempty" bson:"brand"`
	Model  string `json:"model,omitempty" bson:"model"`
	TzName string `json:"tzName,omitempty" bson:"tzName"`

	InAppPushStatus    PushStatus            `json:"inAppPushStatus,omitempty" bson:"inAppPushStatus"`
	SystemPushStatus   PushStatus            `json:"systemPushStatus,omitempty" bson:"systemPushStatus"`
	PrivacyPushStatus  PushStatus            `json:"privacyPushStatus,omitempty" bson:"privacyPushStatus"`
	BusinessPushStatus map[string]PushStatus `json:"businessPushStatus,omitempty" bson:"businessPushStatus"`

	CreateTime      time.Time `json:"createTime,omitempty" bson:"createTime"`
	LastActiveTime  time.Time `json:"lastActiveTime,omitempty" bson:"lastActiveTime"`
	TokenUpdateTime time.Time `json:"tokenUpdateTime,omitempty" bson:"tokenUpdateTime"`
	InfoUpdateTime  time.Time `json:"infoUpdateTime,omitempty" bson:"infoUpdateTime"`
}

func (l *Device) Equal(r *Device) bool {
	return l.ID == r.ID &&
		l.Os == r.Os &&
		l.Brand == r.Os &&
		l.Model == r.Model &&
		l.TzName == r.TzName
}

type TokenSet struct {
	Self string `json:"self" bson:"self"`
	// TODO: other tokens
}

func (l *TokenSet) Equal(r *TokenSet) bool {
	return l.Self == r.Self
}

type App struct {
	ID             int    `json:"appId" bson:"appId" copier:"AppId"`
	AppVersion     string `json:"appVersion,omitempty" bson:"appVersion"`
	PushSdkVersion string `json:"pushSdkVersion,omitempty" bson:"pushSdkVersion"`
	Language       string `json:"language,omitempty" bson:"language"`
}

func (l *App) Equal(r *App) bool {
	return l.ID == r.ID &&
		l.AppVersion == r.AppVersion &&
		l.PushSdkVersion == r.PushSdkVersion &&
		l.Language == r.Language
}
