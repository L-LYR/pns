package event_manager

import (
	"errors"

	"github.com/L-LYR/pns/internal/service/target"
)

const (
	_TargetEventTopic = "target_event"
)

var (
	_TargetEventManager = newManager(_TargetEventTopic)
)

func MustStart() {
	_TargetEventManager.Start()
	_TargetEventManager.Subscribe(_TargetEventTopic, targetEventHandler)
}

func EmitTargetEvent(event *Event) error {
	return _TargetEventManager.Put(event)
}

func MustShutdown() {
	_TargetEventManager.Shutdown()
}

func targetEventHandler(e *Event) error {
	p, ok := e.Payload.(TargetEventPayload)
	if !ok {
		return errors.New("unknown event")
	}
	switch p.TargetEventType {
	case CreateTarget:
		return target.Create(e.Ctx, p.Target)
	case UpdateTarget:
		return target.Update(e.Ctx, p.Target)
	default:
		return errors.New("unknown event type")
	}
}
