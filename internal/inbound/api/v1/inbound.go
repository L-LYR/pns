package v1

type TargetUpsertReq struct {
	DeviceId       string `json:"deviceId" v:"required#device id is required"`
	Os             string `json:"os"`
	Brand          string `json:"brand"`
	Model          string `json:"model"`
	TzName         string `json:"tzName"`
	AppId          int    `json:"appId" v:"required#app is is required"`
	AppVersion     string `json:"appVersion"`
	PushSdkVersion string `json:"pushSdkVersion"`
	Language       string `json:"language"`
}

type TargetUpsertRes struct {
	*CommonRes
}
