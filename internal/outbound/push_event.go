package outbound

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
)

func PushTaskEventConsumer(e event_queue.Event) error {
	pe, ok := e.(*model.PushTaskEvent)
	if !ok {
		return errors.New("not PushEvent")
	}
	switch pe.GetTask().Type {
	case model.PersonalPush:
		return MQTTPusherManager.Handle(pe.GetCtx(), pe.GetTask(), pe.PusherType())
	case model.BroadcastPush:
		panic("not implemented yet")
	default:
		panic("unreachable")
	}
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
