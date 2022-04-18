package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type TaskLogReq struct {
	g.Meta `path:"/task/log" method:"get"`
	TaskId int `json:"taskId"  dc:"push task id, returned by successful /push response" v:"required#task id is required"`
}

type TaskLogRes struct {
	LogEntry []string `json:"logEntry" dc:"readable log entries"`
}

type TaskStatusReq struct {
	g.Meta `path:"/task/status" method:"get"`
	TaskId int `json:"taskId"  dc:"push task id, returned by successful /push response" v:"required#task id is required"`
}

type TaskStatusRes struct {
	Status string
}

type PushLogReq struct {
	g.Meta   `path:"/push/log" method:"get"`
	TaskId   int    `json:"taskId" dc:"push task id" v:"required#push task id is required"`
	AppId    int    `json:"appId" dc:"app id" v:"required#app id is required"`
	DeviceId string `json:"deviceId" dc:"device id" v:"required#device id is required"`
}

type PushLogRes struct {
	LogEntry []string `json:"logEntry" dc:"readable log entries"`
}

// The apis bellow just returns the list of task ids.
// TODO: join the logs

type DeviceLogReq struct {
	g.Meta   `path:"/device/log" method:"get"`
	AppId    int    `json:"appId" dc:"app id" v:"required#app id is required"`
	DeviceId string `json:"deviceId" dc:"device id" v:"required#device id is required"`
}

type DeviceLogRes struct {
	TaskIds []string `json:"taskIds" dc:"broadcast task ids of app"`
}

type AppLogReq struct {
	g.Meta `path:"/app/log" method:"get"`
	AppId  int `json:"appId" dc:"app id" v:"required#app id is required"`
}

type AppLogRes struct {
	TaskIds []string `json:"taskIds" dc:"broadcast task ids of app"`
}
