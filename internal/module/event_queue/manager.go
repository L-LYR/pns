package event_queue

import (
	"errors"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/module/event_queue/internal"
	"github.com/L-LYR/pns/internal/service/target"
)

const (
	TargetEventTopic = "target_event"
)

var (
	_EventQueueManager = internal.NewEventQueue([]string{TargetEventTopic})
)

func MustInit() {
	_EventQueueManager.Start()
	_EventQueueManager.Subscribe(TargetEventTopic, targetEventHandler)
}

func EmitTargetEvent(e *TargetEvent) error {
	return _EventQueueManager.Put(TargetEventTopic, e)
}

func MustShutdown() {
	_EventQueueManager.Shutdown()
}

func targetEventHandler(e internal.Event) error {
	p, ok := e.GetPayload().(*model.Target)
	if !ok {
		return errors.New("unknown event")
	}
	switch e.(*TargetEvent).EventType() {
	case CreateTarget:
		return target.Create(e.GetCtx(), p)
	case UpdateTarget:
		return target.Update(e.GetCtx(), p)
	default:
		return errors.New("unknown event type")
	}
}
