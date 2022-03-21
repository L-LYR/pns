package event_queue

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/target"
)

const (
	_TargetEventTopic = "target_event"
)

type TargetEventType = int8

const (
	CreateTarget TargetEventType = 1
	UpdateTarget TargetEventType = 2
)

type _TargetEvent struct {
	Type   TargetEventType
	Ctx    context.Context
	Target *model.Target
}

func (e *_TargetEvent) GetCtx() context.Context    { return e.Ctx }
func (e *_TargetEvent) GetTarget() *model.Target   { return e.Target }
func (e *_TargetEvent) EventType() TargetEventType { return e.Type }

func SendTargetEvent(
	ctx context.Context,
	target *model.Target,
	t TargetEventType,
) error {
	return _Manager.Queue.Put(
		_TargetEventTopic,
		&_TargetEvent{Type: t, Ctx: ctx, Target: target},
	)
}

var (
	_TargetEventWorker = _NewWorker(_TargetEventTopic, 1, _TargetEventConsumer)
)

func _TargetEventConsumer(e Event) error {
	te, ok := e.(*_TargetEvent)
	if !ok {
		return errors.New("not TargetEvent")
	}
	switch te.EventType() {
	case CreateTarget:
		return target.Create(e.GetCtx(), te.GetTarget())
	case UpdateTarget:
		return target.Update(e.GetCtx(), te.GetTarget())
	default:
		return errors.New("unknown event type")
	}
}
