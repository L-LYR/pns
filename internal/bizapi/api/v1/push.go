package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type DirectPushReq struct {
	g.Meta   `path:"/push/direct" method:"post"`
	DeviceId string `json:"deviceId" dc:"available device id, only" v:"required#device id is required"`
	AppId    int    `json:"appId" dc:"registered app id" v:"required|app-exist#app id is required"`
	Title    string `json:"title" dc:"push message title" v:"required#title is required"`
	Content  string `json:"content" dc:"push message content" v:"required#content is required"`
	Retry    int    `json:"retry" dc:"retry times, -1 for always retry" v:"min:-1#retry time is greater than -1"`
}

type DirectPushRes struct {
	PushTaskId string `json:"pushTaskId" dc:"use task id to tracing it"`
}

type FilterParams struct {
	// TODO: add more
	TimeLimit     *int    `json:"timeLimit" dc:"time limit of broadcast"`
	MinAppVersion *string `json:"minAppVersion" dc:"min app version"`
	MaxAppVersion *string `json:"maxAppVersion" dc:"max app version"`
}

type BroadcastPushReq struct {
	g.Meta  `path:"/push/broadcast" method:"post"`
	AppId   int    `json:"appId" dc:"registered app id" v:"required|app-exist#app id is required"`
	Title   string `json:"title" dc:"push message title" v:"required#title is required"`
	Content string `json:"content" dc:"push message content" v:"required#content is required"`
	// TODO: add
	// FilterParams *FilterParams `json:"filterParams" dc:"filter parameters, optional"`
}

type BroadcastPushRes struct {
	PushTaskId string `json:"pushTaskId" dc:"use task id to tracing it"`
}
