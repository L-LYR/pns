package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/local_storage"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/target"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Push = _PushAPI{}

type _PushAPI struct{}

func (api *_PushAPI) Push(ctx context.Context, req *v1.PushReq) (*v1.PushRes, error) {
	pushTaskId := util.GeneratePushTaskId()

	response := &v1.PushRes{
		PushTaskId: pushTaskId,
	}

	appName, ok := local_storage.GetAppNameByAppId(req.AppId)
	if !ok {
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Unknown app id")
	}

	target, err := target.Query(ctx, req.DeviceId, appName)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to query target")
	}
	if target == nil {
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Target not found")
	}

	task := &model.PushTask{
		ID:     pushTaskId,
		Type:   model.PersonalPush,
		Target: target,
		Message: &model.Message{
			Title:   req.Title,
			Content: req.Content,
		},
	}

	if err := event_queue.EventQueueManager.Put(
		config.PushEventTopic(),
		&model.PushEvent{
			Ctx:    ctx,
			Pusher: model.MQTTPusher,
			Task:   task,
		},
	); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to send push task")
	}

	return response, nil
}
