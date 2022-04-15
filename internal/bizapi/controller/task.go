package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Task = _TaskAPI{}

type _TaskAPI struct{}

func (api *_TaskAPI) Log(ctx context.Context, req *v1.TaskLogReq) (*v1.TaskLogRes, error) {
	entries, err := log.GetTaskLogByID(ctx, req.TaskId)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to get task log")
	}
	res := &v1.TaskLogRes{
		LogEntry: make([]string, 0, len(entries)),
	}
	for i := range entries {
		res.LogEntry = append(res.LogEntry, entries[i].Readable())
	}
	return res, nil
}

func (api *_TaskAPI) Status(ctx context.Context, req *v1.TaskStatusReq) (*v1.TaskStatusRes, error) {
	entry, err := log.GetTaskStatusByID(ctx, req.TaskId)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to get task log")
	}
	if entry == nil {
		return nil, util.FinalError(gcode.CodeNotFound, err, "Task not found")
	}
	return &v1.TaskStatusRes{Status: entry.Readable()}, nil
}
