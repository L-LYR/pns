package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/target"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Push = _PushAPI{}

type _PushAPI struct{}

func (api *_PushAPI) Push(ctx context.Context, request *v1.PushReq) (*v1.PushRes, error) {
	pushTaskId := util.GeneratePushTaskId()

	response := &v1.PushRes{
		PushTaskId: pushTaskId,
	}

	target, err := target.Query(ctx, request.DeviceId, request.AppId)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to query target")
	}
	if target == nil {
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Target not found")
	}

	task := &model.PushTask{
		ID:     pushTaskId,
		Target: target,
		Message: &model.Message{
			Title:   request.Title,
			Content: request.Content,
		},
	}

	if err := event_queue.EventQueueManager.Put(
		config.PushEventTopic(),
		&model.PushEvent{
			Ctx:    ctx,
			Type:   model.Push,
			Pusher: model.MQTTPusher,
			Task:   task,
		},
	); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to send push task")
	}

	return response, nil
}
