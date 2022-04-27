package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/bizcore"
	"github.com/L-LYR/pns/internal/model"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/task"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Push = &_PushAPI{}

type _PushAPI struct{}

func (api *_PushAPI) DirectPush(
	ctx context.Context, req *v1.DirectPushReq,
) (*v1.DirectPushRes, error) {
	res, err := _CreateTask(
		ctx,
		task.NewTaskBuilder(ctx, model.DirectPush).
			SetTaskMeta(req.Retry).
			SetMessage(req.Message).
			SetDirectPushBase(req.DirectPushBase),
	)
	return &v1.DirectPushRes{PushResBase: res}, err
}

func (api *_PushAPI) TemplateDirectPush(
	ctx context.Context, req *v1.TemplateDirectPushReq,
) (*v1.TemplateDirectPushRes, error) {
	res, err := _CreateTask(
		ctx,
		task.NewTaskBuilder(ctx, model.DirectPush).
			SetTaskMeta(req.Retry).
			SetTemplateMessage(req.Message).
			SetDirectPushBase(req.DirectPushBase),
	)
	return &v1.TemplateDirectPushRes{PushResBase: res}, err
}

func (api *_PushAPI) BroadcastPush(
	ctx context.Context, req *v1.BroadcastPushReq,
) (*v1.BroadcastPushRes, error) {
	res, err := _CreateTask(
		ctx,
		task.NewTaskBuilder(ctx, model.BroadcastPush).
			SetTaskMeta(-1).
			SetMessage(req.Message).
			SetBroadcastPushBase(req.BroadcastPushBase),
	)

	return &v1.BroadcastPushRes{PushResBase: res}, err
}

func (api *_PushAPI) TemplateBroadcastPush(
	ctx context.Context, req *v1.TemplateBroadcastPushReq,
) (*v1.TemplateBroadcastPushRes, error) {
	res, err := _CreateTask(
		ctx,
		task.NewTaskBuilder(ctx, model.BroadcastPush).
			SetTaskMeta(-1).
			SetTemplateMessage(req.Message).
			SetBroadcastPushBase(req.BroadcastPushBase),
	)
	return &v1.TemplateBroadcastPushRes{PushResBase: res}, err
}

func _CreateTask(ctx context.Context, taskBuilder task.TaskBuilder) (*v1.PushResBase, error) {
	task, err := taskBuilder.Build()
	if err != nil {
		log.PutTaskLogEvent(ctx, task.GetLogMeta(), "task creation", "failure")
		return nil, err
	}

	log.PutTaskLogEvent(ctx, task.GetLogMeta(), "task creation", "success")
	
	if err := bizcore.PutTaskValidationEvent(ctx, task); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to send push task")
	}


	return &v1.PushResBase{
		PushTaskId: task.GetID(),
	}, nil
}
