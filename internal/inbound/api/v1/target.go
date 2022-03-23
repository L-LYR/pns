package v1

import (
	"github.com/L-LYR/pns/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type TargetCreateReq struct {
	g.Meta         `path:"/target" method:"post"`
	DeviceId       string `json:"deviceId" v:"required#device id is required"`
	Os             string `json:"os" v:"required#device os is required"`
	Brand          string `json:"brand" v:"required#device brand is required"`
	Model          string `json:"model" v:"required#device model is required"`
	TzName         string `json:"tzName" v:"required#time zone is required"`
	AppId          int    `json:"appId" v:"required#app id is required"`
	AppVersion     string `json:"appVersion" v:"required# app version is required"`
	PushSdkVersion string `json:"pushSdkVersion" v:"required#push sdk version is required"`
	Language       string `json:"language" v:"required#language is required"`
}

type TargetCreateRes struct{}

type TargetUpdateReq struct {
	g.Meta   `path:"/target" method:"patch"`
	DeviceId string            `json:"deviceId" v:"required#device id is required"`
	Os       string            `json:"os"`
	Brand    string            `json:"brand"`
	Model    string            `json:"model"`
	TzName   string            `json:"tzName"`
	Tokens   map[string]string `json:"tokens"`

	InAppPushStatus    int            `json:"inAppPushStatus"`
	SystemPushStatus   int            `json:"systemPushStatus"`
	PrivacyPushStatus  int            `json:"privacyPushStatus"`
	BusinessPushStatus map[string]int `json:"businessPushStatus"`

	AppId          int    `json:"appId" v:"required#app id is required"`
	AppVersion     string `json:"appVersion"`
	PushSdkVersion string `json:"pushSdkVersion"`
	Language       string `json:"language"`
}

type TargetUpdateRes struct{}

type TargetQueryReq struct {
	g.Meta   `path:"/target" method:"get"`
	DeviceId string `json:"deviceId" v:"required#device id is required"`
	AppId    int    `json:"appId" v:"required#app id is required"`
}

type TargetQueryRes struct {
	Target *model.Target `json:"target"`
}
