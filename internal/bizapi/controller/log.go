package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/model"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Log = &_LogAPI{}

type _LogAPI struct{}

func (api *_LogAPI) TaskLog(ctx context.Context, req *v1.TaskLogReq) (*v1.TaskLogRes, error) {
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

func (api *_LogAPI) TaskStatus(ctx context.Context, req *v1.TaskStatusReq) (*v1.TaskStatusRes, error) {
	entry, err := log.GetTaskStatusByID(ctx, req.TaskId)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to get task log")
	}
	if entry == nil {
		return nil, nil
	}
	if entry.Where == model.TaskDone {
		if status, err := log.GetTaskStatisticsByID(ctx, req.TaskId); err == nil {
			return &v1.TaskStatusRes{Status: status}, nil
		}
		// ignore this error
	}
	return &v1.TaskStatusRes{Status: entry.Status()}, nil
}

func (api *_LogAPI) PushLog(ctx context.Context, req *v1.PushLogReq) (*v1.PushLogRes, error) {
	entries, err := log.GetPushLogByMeta(ctx, &model.LogMeta{
		TaskId:   req.TaskId,
		AppId:    req.AppId,
		DeviceId: req.DeviceId,
	})
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to get push log")
	}
	res := &v1.PushLogRes{
		LogEntry: make([]string, 0, len(entries)),
	}
	for i := range entries {
		res.LogEntry = append(res.LogEntry, entries[i].Readable())
	}
	return res, nil
}

func (api *_LogAPI) DeviceLogs(ctx context.Context, req *v1.DeviceLogReq) (*v1.DeviceLogRes, error) {
	taskIds, err := log.GetTaskEntryListByMeta(ctx, &model.LogMeta{
		AppId:    req.AppId,
		DeviceId: req.DeviceId,
	})
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to get device log")
	}
	return &v1.DeviceLogRes{TaskIds: taskIds}, nil
}

func (api *_LogAPI) AppLogs(ctx context.Context, req *v1.AppLogReq) (*v1.AppLogRes, error) {
	taskIds, err := log.GetTaskEntryListByMeta(ctx, &model.LogMeta{
		AppId: req.AppId,
	})
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to get device log")
	}
	return &v1.AppLogRes{TaskIds: taskIds}, nil
}
