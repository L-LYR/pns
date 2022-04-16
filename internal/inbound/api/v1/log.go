package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type LogReq struct {
	g.Meta    `path:"/log" method:"post"`
	Where     string `json:"where" dc:"log location" v:"required#location is required"`
	Timestamp string `json:"timestamp" dc:"log timestamp" v:"required#timestamp is required"`
	Hint      string `json:"hint" dc:"log hint" v:"required#hint is required"`
	TaskId    string `json:"taskId" dc:"push task id" v:"required#task id is required"`
	AppId     int    `json:"appId" dc:"app id" v:"required|app-exist#app id is required"`
	DeviceId  string `json:"deviceId" dc:"device id" v:"required#device id is required"`
}

type LogRes struct{}
