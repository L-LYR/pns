package target

import (
	"errors"

	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
)

func TargetEventConsumer(e event_queue.Event) error {
	te, ok := e.(*model.TargetEvent)
	if !ok {
		return errors.New("not TargetEvent")
	}
	switch te.EventType() {
	case model.CreateTarget:
		return Create(e.GetCtx(), te.GetTarget())
	case model.UpdateTarget:
		return Update(e.GetCtx(), te.GetTarget())
	default:
		return errors.New("unknown event type")
	}
}
