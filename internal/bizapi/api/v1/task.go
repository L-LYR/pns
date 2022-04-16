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
