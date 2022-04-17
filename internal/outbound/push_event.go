package outbound

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
	log "github.com/L-LYR/pns/internal/service/push_log"
)

func PushTaskEventConsumer(e event_queue.Event) error {
	pe, ok := e.(*model.PushTaskEvent)
	if !ok {
		return errors.New("not PushEvent")
	}

	log.PutTaskLog(
		pe.GetCtx(), pe.GetTask().GetLogMeta(),
		"task handle", "success",
	)

	var err error
	taskHint := "success"

	switch pe.GetTask().GetPusher() {
	case model.MQTTPusher:
		err = MQTTPusherManager.Handle(pe.GetCtx(), pe.GetTask())
	default:
		panic("unreachable")
	}

	if err != nil {
		taskHint = "fail"
	}

	log.PutTaskLog(
		pe.GetCtx(), pe.GetTask().GetLogMeta(),
		"task done", taskHint,
	)

	return err
}

func PutPushTaskEvent(ctx context.Context, task model.PushTask) error {
	return event_queue.EventQueueManager.Put(
		config.PushEventTopic(),
		&model.PushTaskEvent{
			Ctx:  ctx,
			Task: task,
		},
	)
}
