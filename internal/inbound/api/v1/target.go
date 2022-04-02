package v1

import (
	"github.com/L-LYR/pns/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type UpdateTargetReq struct {
	g.Meta   `path:"/target" method:"post"`
	DeviceId string `json:"deviceId" v:"required#device id is required" copier:"DeviceId"`
	Os       string `json:"os" v:"required#device os is required"`
	Brand    string `json:"brand" v:"required#device brand is required"`
	Model    string `json:"model" v:"required#device model is required"`
	TzName   string `json:"tzName" v:"required#time zone is required"`

	InAppPushStatus    int            `json:"inAppPushStatus"`
	SystemPushStatus   int            `json:"systemPushStatus"`
	PrivacyPushStatus  int            `json:"privacyPushStatus"`
	BusinessPushStatus map[string]int `json:"businessPushStatus"`

	AppId          int    `json:"appId" v:"required#app id is required" copier:"AppId"`
	AppVersion     string `json:"appVersion" v:"required# app version is required"`
	PushSdkVersion string `json:"pushSdkVersion" v:"required#push sdk version is required"`
	Language       string `json:"language" v:"required#language is required"`
}

type UpdateTargetRes struct{}

type QueryTargetReq struct {
	g.Meta   `path:"/target" method:"get"`
	DeviceId string `json:"deviceId" v:"required#device id is required"`
	AppId    int    `json:"appId" v:"required#app id is required"`
}

type QueryTargetRes struct {
	Target *model.Target `json:"target"`
}

type GetTokenReq struct {
	g.Meta   `path:"/token" method:"get"`
	DeviceId string `json:"deviceId" v:"required#device id is required"`
	AppId    int    `json:"appId" v:"required#app id is required"`
}

type GetTokenRes struct {
	Token string `json:"token"`
}
