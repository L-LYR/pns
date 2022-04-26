package controller

import (
	"context"
	"errors"
	"strconv"
	"time"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/bizcore"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/service/target"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Push = _PushAPI{}

type _PushAPI struct{}

func (api *_PushAPI) DirectPush(ctx context.Context, req *v1.DirectPushReq) (*v1.DirectPushRes, error) {
	task, err := _BuildDirectPushTask(ctx, req)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to build push task")
	}

	log.PutTaskLogEvent(ctx, task.GetLogMeta(), "direct push task creation", "success")

	if err := bizcore.PutTaskValidationEvent(ctx, task); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to send push task")
	}

	return &v1.DirectPushRes{
		PushTaskId: strconv.FormatInt(int64(task.ID), 10),
	}, nil
}

func _BuildDirectPushTask(ctx context.Context, req *v1.DirectPushReq) (*model.DirectPushTask, error) {
	target, err := target.Query(ctx, req.AppId, req.DeviceId)
	if err != nil {
		return nil, errors.New("fail to query target")
	}
	if target == nil {
		return nil, errors.New("target not found")
	}
	return &model.DirectPushTask{
		ID:     util.GeneratePushTaskId(),
		Pusher: model.MQTTPusher,
		Qos:    config.CommonTaskQos(),
		Target: target,
		Message: &model.Message{
			Title:   req.Title,
			Content: req.Content,
		},
		PushTaskMeta: &model.PushTaskMeta{
			RetryCounter: &model.RetryCounter{
				Counter: model.RetryTimes(req.Retry),
			},
			CreationTime: time.Now(),
		},
	}, nil
}

func (api *_PushAPI) BroadcastPush(ctx context.Context, req *v1.BroadcastPushReq) (*v1.BroadcastPushRes, error) {
	task, err := _BuildBroadcastPushTask(ctx, req)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to build push task")
	}

	log.PutTaskLogEvent(ctx, task.GetLogMeta(), "broadcast push task creation", "success")

	if err := bizcore.PutTaskValidationEvent(ctx, task); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to send push task")
	}

	return &v1.BroadcastPushRes{
		PushTaskId: strconv.FormatInt(int64(task.ID), 10),
	}, nil
}

func _BuildBroadcastPushTask(ctx context.Context, req *v1.BroadcastPushReq) (*model.BroadcastTask, error) {
	return &model.BroadcastTask{
		ID:     util.GeneratePushTaskId(),
		AppId:  req.AppId,
		Pusher: model.MQTTPusher,
		Qos:    config.CommonTaskQos(),
		Message: &model.Message{
			Title:   req.Title,
			Content: req.Content,
		},
		PushTaskMeta: &model.PushTaskMeta{
			RetryCounter: &model.RetryCounter{
				Counter: model.NeverRetry,
			},
			CreationTime: time.Now(),
		},
	}, nil
}
