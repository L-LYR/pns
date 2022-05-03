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
	"golang.org/x/time/rate"
)

var Push = &_PushAPI{
	limiter: rate.NewLimiter(1100, 2000),
}

type _PushAPI struct {
	limiter *rate.Limiter
}

func (api *_PushAPI) DirectPush(
	ctx context.Context, req *v1.DirectPushReq,
) (*v1.DirectPushRes, error) {
	res, err := api.createTask(
		ctx,
		task.NewTaskBuilder(ctx, model.DirectPush).
			SetTaskMeta(req.Retry, req.IgnoreFreqCtrl, req.IgnoreOnlineCheck).
			SetMessage(req.Message).
			SetDirectPushBase(req.DirectPushBase),
	)
	return &v1.DirectPushRes{PushResBase: res}, err
}

func (api *_PushAPI) TemplateDirectPush(
	ctx context.Context, req *v1.TemplateDirectPushReq,
) (*v1.TemplateDirectPushRes, error) {
	res, err := api.createTask(
		ctx,
		task.NewTaskBuilder(ctx, model.DirectPush).
			SetTaskMeta(req.Retry, req.IgnoreFreqCtrl, req.IgnoreOnlineCheck).
			SetTemplateMessage(req.Message).
			SetDirectPushBase(req.DirectPushBase),
	)
	return &v1.TemplateDirectPushRes{PushResBase: res}, err
}

func (api *_PushAPI) BroadcastPush(
	ctx context.Context, req *v1.BroadcastPushReq,
) (*v1.BroadcastPushRes, error) {
	res, err := api.createTask(
		ctx,
		task.NewTaskBuilder(ctx, model.BroadcastPush).
			SetTaskMeta(-1, req.IgnoreFreqCtrl, false).
			SetMessage(req.Message).
			SetBroadcastPushBase(req.BroadcastPushBase),
	)

	return &v1.BroadcastPushRes{PushResBase: res}, err
}

func (api *_PushAPI) TemplateBroadcastPush(
	ctx context.Context, req *v1.TemplateBroadcastPushReq,
) (*v1.TemplateBroadcastPushRes, error) {
	res, err := api.createTask(
		ctx,
		task.NewTaskBuilder(ctx, model.BroadcastPush).
			SetTaskMeta(-1, req.IgnoreFreqCtrl, false).
			SetTemplateMessage(req.Message).
			SetBroadcastPushBase(req.BroadcastPushBase),
	)
	return &v1.TemplateBroadcastPushRes{PushResBase: res}, err
}

func (api *_PushAPI) RangePush(
	ctx context.Context, req *v1.RangePushReq,
) (*v1.RangePushRes, error) {
	res, err := api.createTask(
		ctx,
		task.NewTaskBuilder(ctx, model.RangePush).
			SetTaskMeta(-1, req.IgnoreFreqCtrl, false).
			SetMessage(req.Message).
			SetBroadcastPushBase(req.BroadcastPushBase).
			SetFilterParams(req.FilterParams),
	)

	return &v1.RangePushRes{PushResBase: res}, err
}

func (api *_PushAPI) TemplateRangePush(
	ctx context.Context, req *v1.TemplateRangePushReq,
) (*v1.TemplateRangePushRes, error) {
	res, err := api.createTask(
		ctx,
		task.NewTaskBuilder(ctx, model.RangePush).
			SetTaskMeta(-1, req.IgnoreFreqCtrl, false).
			SetTemplateMessage(req.Message).
			SetBroadcastPushBase(req.BroadcastPushBase).
			SetFilterParams(req.FilterParams),
	)
	return &v1.TemplateRangePushRes{PushResBase: res}, err
}

func (api *_PushAPI) createTask(ctx context.Context, taskBuilder task.TaskBuilder) (*v1.PushResBase, error) {
	t, err := taskBuilder.Build()
	if err != nil {
		log.PutTaskLogEvent(ctx, t.GetLogMeta(), model.TaskCreation, "failure")
		return nil, err
	}

	log.PutTaskLogEvent(ctx, t.GetLogMeta(), model.TaskCreation, "success")
	if t.GetType() == model.RangePush {
		// range task runner
		go task.NewRangePushTaskRunner(ctx, model.AsRangePushTask(t), api.limiter).Run()
	} else {
		api.limiter.Wait(ctx)
		if err := bizcore.PutTaskValidationEvent(ctx, t); err != nil {
			return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to send push task")
		}
	}

	return &v1.PushResBase{
		PushTaskId: t.GetID(),
	}, nil
}
