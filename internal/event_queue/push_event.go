package event_queue

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

const (
	_PushEventTopic = "push_event"
)

type PushEventType = int8

const (
	Push PushEventType = 1
)

type PushEvent struct {
	Type PushEventType
	Ctx  context.Context
	Task *model.PushTask
}

func (e *PushEvent) GetCtx() context.Context  { return e.Ctx }
func (e *PushEvent) GetTask() *model.PushTask { return e.Task }
func (e *PushEvent) EventType() PushEventType { return e.Type }

func SendPushEvent(
	ctx context.Context,
	task *model.PushTask,
	t PushEventType,
) error {
	return _Manager.Queue.Put(
		_PushEventTopic,
		&PushEvent{Type: t, Ctx: ctx, Task: task},
	)
}

var (
	PushWorker = _NewWorker(_PushEventTopic, 1, nil)
)
