package controller

import (
	"context"
	"strconv"

	v1 "github.com/L-LYR/pns/internal/inbound/api/v1"
	"github.com/L-LYR/pns/internal/model"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Log = &_LogAPI{}

type _LogAPI struct{}

func (api *_LogAPI) Log(ctx context.Context, req *v1.LogReq) (*v1.LogRes, error) {
	taskId, err := strconv.ParseInt(req.TaskId, 10, 64)
	if err != nil {
		return nil, util.FinalError(gcode.CodeValidationFailed, err, "Fail to parse task id")
	}

	ts, err := strconv.ParseInt(req.Timestamp, 10, 64)
	if err != nil {
		return nil, util.FinalError(gcode.CodeValidationFailed, err, "Fail to parse timestamp")
	}

	meta := &model.LogMeta{
		TaskId:   taskId,
		AppId:    req.AppId,
		DeviceId: req.DeviceId,
	}

	log.PutPushLogEvent(ctx, meta, model.PushTaskStage(req.Where), ts, req.Hint)

	return &v1.LogRes{}, nil
}
