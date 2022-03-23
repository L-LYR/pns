package event_queue

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/outbound"
)

const (
	_PushEventTopic = "push_event"
)

type PushEventType = int8

const (
	Push PushEventType = 1
)

type _PushEvent struct {
	Type   PushEventType
	Ctx    context.Context
	Pusher model.PusherType
	Task   *model.PushTask
}

func (e *_PushEvent) GetCtx() context.Context      { return e.Ctx }
func (e *_PushEvent) GetTask() *model.PushTask     { return e.Task }
func (e *_PushEvent) EventType() PushEventType     { return e.Type }
func (e *_PushEvent) PusherType() model.PusherType { return e.Pusher }

func SendPushEvent(
	ctx context.Context,
	task *model.PushTask,
	t PushEventType,
	pusher model.PusherType,
) error {
	return _Manager.Queue.Put(
		_PushEventTopic,
		&_PushEvent{Type: t, Ctx: ctx, Task: task, Pusher: pusher},
	)
}

var (
	_PushEventWorker = _MustNewWorker(_PushEventTopic, 1, _PushEventConsumer)
)

func _PushEventConsumer(e _Event) error {
	pe, ok := e.(*_PushEvent)
	if !ok {
		return errors.New("not PushEvent")
	}
	switch pe.EventType() {
	case Push:
		return outbound.Handle(pe.GetCtx(), pe.GetTask(), pe.PusherType())
	default:
		return errors.New("unknown event type")
	}
}
