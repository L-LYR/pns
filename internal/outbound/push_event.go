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
		pe.GetCtx(), pe.GetTask().LogMeta(),
		"task handle", "success",
	)
	return MQTTPusherManager.Handle(pe.GetCtx(), pe.GetTask(), pe.PusherType())
}

func PutPushTaskEvent(ctx context.Context, task *model.PushTask, pusherType model.PusherType) error {
	return event_queue.EventQueueManager.Put(
		config.PushEventTopic(),
		&model.PushTaskEvent{
			Ctx:    ctx,
			Pusher: pusherType,
			Task:   task,
		},
	)
}
