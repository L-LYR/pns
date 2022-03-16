package v1

type TargetCreateRequest struct {
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

type TargetUpdateRequest struct {
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

type TargetQueryRequest struct {
	DeviceId string `json:"deviceId" v:"required#device id is required"`
	AppId    int    `json:"appId" v:"required#app is required"`
}