package controller

import (
	"context"
	"strconv"
	"time"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/outbound"
	"github.com/L-LYR/pns/internal/service/cache"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/service/target"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Push = _PushAPI{}

type _PushAPI struct{}

func (api *_PushAPI) Push(ctx context.Context, req *v1.PushReq) (*v1.PushRes, error) {
	appName, ok := cache.Config.GetAppNameByAppId(req.AppId)
	if !ok {
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Unknown app id")
	}

	target, err := target.Query(ctx, appName, req.DeviceId)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to query target")
	}
	if target == nil {
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Target not found")
	}

	task := &model.PushTask{
		ID:     util.GeneratePushTaskId(),
		Type:   model.PersonalPush,
		Target: target,
		Message: &model.Message{
			Title:   req.Title,
			Content: req.Content,
		},
	}

	meta := task.LogMeta()
	if err := log.PutTaskRequestLog(ctx, meta); err != nil {
		util.GLog.Warning(ctx, "Fail to add task list entry")
	}

	if err := log.PutLogEvent(
		ctx, meta,
		time.Now().UnixMilli(), "task creation", "success",
	); err != nil {
		util.GLog.Warning(ctx, "Fail to put task creation log")
	}

	if err := outbound.PutMQTTPushEvent(ctx, task); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to send push task")
	}

	return &v1.PushRes{
		PushTaskId: strconv.FormatInt(int64(task.ID), 10),
	}, nil
}
