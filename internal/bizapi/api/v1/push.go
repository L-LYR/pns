package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PushReq struct {
	// TODO: authorization
	g.Meta   `path:"/push" method:"post"`
	AppId    int    `json:"appId" dc:"registered app id" v:"required|app-exist#app id is required"`
	DeviceId string `json:"deviceId" dc:"available device id" v:"required#device id is required"`
	Title    string `json:"title" dc:"push message title" v:"required#title is required"`
	Content  string `json:"content" dc:"push message content" v:"required#content is required"`
}

type PushRes struct {
	PushTaskId string `json:"pushTaskId" dc:"use task id to tracing it"`
}
